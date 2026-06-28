package models

// ─── API 统一响应（与主项目一致） ───

type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type PaginatedData struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalCount int         `json:"totalCount"`
}

func OK(data interface{}) APIResponse {
	return APIResponse{Code: 200, Msg: "查询成功", Data: data}
}

func OKPaginated(data interface{}, page, pageSize, total int) APIResponse {
	return APIResponse{
		Code: 200,
		Msg:  "查询成功",
		Data: PaginatedData{
			Items:      data,
			Page:       page,
			PageSize:   pageSize,
			TotalCount: total,
		},
	}
}

func Err(code int, msg string) APIResponse {
	return APIResponse{Code: code, Msg: msg, Data: nil}
}

// ─── 日历 ───

type CalendarInfo struct {
	CalendarID   int    `json:"calendarId"`
	CalendarName string `json:"calendarName"`
	DbName       string `json:"-"`
}

// ─── 字段定义（后端唯一权威源，前端按 SearchType 决定交互） ───

type FieldDef struct {
	ID         string `json:"id"`         // 英文标识，搜索请求用这个值
	Label      string `json:"label"`      // 中文显示名
	SearchType string `json:"searchType"` // "select" 下拉 | "input" 自由输入
	DBColumn   string `json:"-"`          // SQL 列名（白名单注入，不暴露给前端）
}

var AllFields = []FieldDef{
	{ID: "courseCode", Label: "课程序号", SearchType: "input", DBColumn: "cd.code"},
	{ID: "newCode", Label: "新课程序号", SearchType: "input", DBColumn: "cd.newCode"},
	{ID: "courseName", Label: "课程名称", SearchType: "input", DBColumn: "cd.courseName"},
	{ID: "teacherName", Label: "授课教师", SearchType: "input", DBColumn: "t.teacherName"},
	{ID: "teacherCode", Label: "教师工号", SearchType: "input", DBColumn: "t.teacherCode"},
	{ID: "schedule", Label: "排课信息", SearchType: "input", DBColumn: "t.arrangeInfoText"},
	{ID: "campus", Label: "校区", SearchType: "select", DBColumn: "cp.campusI18n"},
	{ID: "faculty", Label: "开课学院", SearchType: "select", DBColumn: "f.facultyI18n"},
	{ID: "courseType", Label: "课程性质", SearchType: "select", DBColumn: "cn.courseLabelName"},
	{ID: "major", Label: "听课专业", SearchType: "select", DBColumn: "m.name"},
}

var FieldByID = make(map[string]FieldDef)

func init() {
	for _, f := range AllFields {
		FieldByID[f.ID] = f
	}
}

// ─── 字段候选项（仅 select 类型字段） ───

type FieldOption struct {
	Value string `json:"value"` // 发送到后端的值
	Label string `json:"label"` // 前端显示
}

type FieldOptionsRequest struct {
	Field       string `json:"field"`
	CalendarIDs []int  `json:"calendar_ids"`
}

// ─── 搜索 ───

type SearchCondition struct {
	Field     string `json:"field"`     // 字段 id（如 "courseCode"）
	MatchType string `json:"matchType"` // "contains" | "not_contains"
	Value     string `json:"value"`
}

type SearchGroup struct {
	Conditions []SearchCondition `json:"conditions"`
}

type SearchRequest struct {
	Groups      []SearchGroup `json:"groups"`
	CalendarIDs []int         `json:"calendar_ids"`
	Page        int           `json:"page"`
	PageSize    int           `json:"page_size"`
}

type SearchResult struct {
	CalendarID     int    `json:"calendarId"`
	CalendarName   string `json:"term"`
	NewCode        *string `json:"newCode"`
	CourseCode     string `json:"courseCode"`
	CourseName     string `json:"courseName"`
	Campus         string `json:"campus"`
	Faculty        string `json:"faculty"`
	Majors         string `json:"majors"`
	TotalHours     *int   `json:"totalHours"`
	CourseType     string `json:"courseType"`
	AssessmentMode string `json:"assessmentMode"`
	Capacity       *int   `json:"capacity"`
	Enrolled       *int   `json:"enrolled"`
	Teachers       string `json:"teachers"`
	Schedule       string `json:"schedule"`
}
