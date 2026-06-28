package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"jkd-website/backend/models"
)

// ─── 公共方法 ───

// GetLatestUpdate 返回元数据库 fetchlog 最新同步时间。
func (r *Router) GetLatestUpdate() (*time.Time, string, error) {
	var t sql.NullTime
	var msg sql.NullString
	err := r.metaDB.QueryRow(
		"SELECT startTime, msg FROM fetchlog ORDER BY startTime DESC LIMIT 1",
	).Scan(&t, &msg)
	if err != nil {
		return nil, "", err
	}
	if t.Valid {
		return &t.Time, msg.String, nil
	}
	return nil, "", nil
}

// Search 跨库搜索：按 calendarId 列表遍历各学期库，合并结果。
func (r *Router) Search(req models.SearchRequest) ([]models.SearchResult, error) {
	var all []models.SearchResult
	for _, id := range req.CalendarIDs {
		cal, ok := r.CalendarByID(id)
		if !ok {
			continue
		}
		db, err := r.GetConnection(cal.DbName)
		if err != nil {
			return nil, err
		}
		results, err := searchInDB(db, req, cal)
		if err != nil {
			return nil, err
		}
		all = append(all, results...)
	}
	return all, nil
}

// GetFieldOptions 跨库获取 select 字段的候选项。
func (r *Router) GetFieldOptions(req models.FieldOptionsRequest) ([]models.FieldOption, error) {
	calendars := r.filterCalendarsByID(req.CalendarIDs)
	if len(calendars) == 0 {
		return nil, nil
	}
	def, ok := models.FieldByID[req.Field]
	if !ok {
		return nil, nil
	}
	query := selectFieldOptionQueries[def.ID]
	if query == "" {
		return nil, nil
	}

	seen := make(map[string]bool)
	var opts []models.FieldOption
	for _, cal := range calendars {
		db, err := r.GetConnection(cal.DbName)
		if err != nil {
			return nil, err
		}
		rows, err := db.Query(query)
		if err != nil {
			return nil, fmt.Errorf("查询字段选项失败: %w", err)
		}
		for rows.Next() {
			var v sql.NullString
			if err := rows.Scan(&v); err != nil {
				rows.Close()
				return nil, err
			}
			if v.Valid && v.String != "" && !seen[v.String] {
				seen[v.String] = true
				opts = append(opts, models.FieldOption{Value: v.String, Label: v.String})
			}
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	if len(opts) > 800 {
		opts = opts[:800]
	}
	return opts, nil
}

// ─── SQL 模板 ───

const searchBaseQuery = `
SELECT
	cd.code,
	cd.courseName,
	cp.campusI18n,
	f.facultyI18n,
	(SELECT GROUP_CONCAT(DISTINCT m2.name ORDER BY m2.name SEPARATOR ', ')
	 FROM majorandcourse mc2
	 LEFT JOIN major m2 ON m2.id = mc2.majorId
	 WHERE mc2.courseId = cd.id) AS majors,
	cd.period,
	cn.courseLabelName,
	a.assessmentModeI18n,
	cd.number,
	cd.elcNumber,
	GROUP_CONCAT(DISTINCT CONCAT(t.teacherName, ' (', t.teacherCode, ')')
		ORDER BY t.teacherName SEPARATOR ', ') AS teachers,
	GROUP_CONCAT(DISTINCT t.arrangeInfoText SEPARATOR '\n')  AS schedule
FROM coursedetail cd
LEFT JOIN coursenature cn ON cd.courseLabelId = cn.courseLabelId
LEFT JOIN assessment   a  ON cd.assessmentMode = a.assessmentMode
LEFT JOIN campus       cp ON cd.campus         = cp.campus
LEFT JOIN faculty      f  ON cd.faculty        = f.faculty
LEFT JOIN teacher      t  ON cd.id             = t.teachingClassId`

// selectFieldOptionQueries 按 field ID 索引。
var selectFieldOptionQueries = map[string]string{
	"campus":     `SELECT DISTINCT c.campusI18n FROM campus c INNER JOIN coursedetail cd ON c.campus = cd.campus WHERE c.campusI18n IS NOT NULL`,
	"faculty":    `SELECT DISTINCT f.facultyI18n FROM faculty f INNER JOIN coursedetail cd ON f.faculty = cd.faculty WHERE f.facultyI18n IS NOT NULL`,
	"courseType": `SELECT DISTINCT cn.courseLabelName FROM coursenature cn INNER JOIN coursedetail cd ON cn.courseLabelId = cd.courseLabelId WHERE cn.courseLabelName IS NOT NULL`,
	"major":      `SELECT DISTINCT m.name FROM major m INNER JOIN majorandcourse mc ON m.id = mc.majorId INNER JOIN coursedetail cd ON mc.courseId = cd.id ORDER BY m.grade DESC`,
}

// ─── 内部函数 ───

func searchInDB(db *sql.DB, req models.SearchRequest, cal models.CalendarInfo) ([]models.SearchResult, error) {
	where, params := buildWhere(req.Groups)
	if where == "" {
		return nil, nil
	}

	query := searchBaseQuery + "\nWHERE " + where + `
GROUP BY cd.id, cd.code, cd.courseName,
         cp.campusI18n, f.facultyI18n, cd.period,
         cn.courseLabelName, a.assessmentModeI18n,
         cd.number, cd.elcNumber
ORDER BY cd.code ASC`

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("搜索查询失败: %w", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		results = append(results, scanResult(rows, cal))
	}
	return results, rows.Err()
}

func scanResult(rows *sql.Rows, cal models.CalendarInfo) models.SearchResult {
	var r models.SearchResult
	var campus, faculty, courseType, assessment sql.NullString
	var majors, teachers, schedule sql.NullString
	var totalHours, capacity, enrolled sql.NullInt64

	_ = rows.Scan(
		&r.CourseCode, &r.CourseName,
		&campus, &faculty,
		&majors, &totalHours, &courseType, &assessment,
		&capacity, &enrolled,
		&teachers, &schedule,
	)

	r.CalendarID = cal.CalendarID
	r.CalendarName = cal.CalendarName
	r.Campus = campus.String
	r.Faculty = faculty.String
	r.Majors = majors.String
	r.CourseType = courseType.String
	r.AssessmentMode = assessment.String
	r.Teachers = teachers.String
	r.Schedule = dedupLines(schedule.String)

	if totalHours.Valid {
		v := int(totalHours.Int64)
		r.TotalHours = &v
	}
	if capacity.Valid {
		v := int(capacity.Int64)
		r.Capacity = &v
	}
	if enrolled.Valid {
		v := int(enrolled.Int64)
		r.Enrolled = &v
	}
	return r
}

// buildWhere 从搜索条件组构建 WHERE 子句（白名单映射，杜绝 SQL 注入）。
func buildWhere(groups []models.SearchGroup) (string, []interface{}) {
	var clauses []string
	var params []interface{}

	for _, g := range groups {
		var ors []string
		for _, c := range g.Conditions {
			def, ok := models.FieldByID[c.Field]
			if !ok || c.Value == "" {
				continue
			}
			switch c.MatchType {
			case "not_contains":
				ors = append(ors, def.DBColumn+" NOT LIKE ?")
			default:
				ors = append(ors, def.DBColumn+" LIKE ?")
			}
			params = append(params, "%"+c.Value+"%")
		}
		if len(ors) > 0 {
			clauses = append(clauses, "("+strings.Join(ors, " OR ")+")")
		}
	}

	if len(clauses) == 0 {
		return "", nil
	}
	return strings.Join(clauses, " AND "), params
}

func dedupLines(s string) string {
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	seen := make(map[string]bool)
	var uniq []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || seen[line] {
			continue
		}
		seen[line] = true
		uniq = append(uniq, line)
	}
	return strings.Join(uniq, "\n")
}

func (r *Router) filterCalendarsByID(ids []int) []models.CalendarInfo {
	if len(ids) == 0 {
		return r.Calendars()
	}
	var filtered []models.CalendarInfo
	for _, id := range ids {
		if c, ok := r.CalendarByID(id); ok {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
