package db

import (
	"strings"
	"testing"

	"jkd-website/backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── buildWhere 单元测试 ───

func TestBuildWhereEmpty(t *testing.T) {
	where, params := buildWhere(nil)
	assert.Equal(t, "1=1", where)
	assert.Empty(t, params)

	where, params = buildWhere([]models.SearchGroup{})
	assert.Equal(t, "1=1", where)
	assert.Empty(t, params)
}

func TestBuildWhereSingleCondition(t *testing.T) {
	groups := []models.SearchGroup{
		{Conditions: []models.SearchCondition{
			{Field: "teacherName", MatchType: "contains", Value: "高珍"},
		}},
	}
	where, params := buildWhere(groups)
	assert.Equal(t, "(t.teacherName LIKE ?)", where)
	require.Len(t, params, 1)
	assert.Equal(t, "%高珍%", params[0])
}

func TestBuildWhereNotContains(t *testing.T) {
	groups := []models.SearchGroup{
		{Conditions: []models.SearchCondition{
			{Field: "campus", MatchType: "not_contains", Value: "嘉定校区"},
		}},
	}
	where, params := buildWhere(groups)
	assert.Equal(t, "(cp.campusI18n NOT LIKE ?)", where)
	require.Len(t, params, 1)
	assert.Equal(t, "%嘉定校区%", params[0])
}

func TestBuildWhereMultipleConditionsInGroup(t *testing.T) {
	// 组内 OR
	groups := []models.SearchGroup{
		{Conditions: []models.SearchCondition{
			{Field: "teacherName", MatchType: "contains", Value: "高珍"},
			{Field: "teacherName", MatchType: "contains", Value: "卫志华"},
		}},
	}
	where, params := buildWhere(groups)
	assert.Equal(t, "(t.teacherName LIKE ? OR t.teacherName LIKE ?)", where)
	assert.Equal(t, []interface{}{"%高珍%", "%卫志华%"}, params)
}

func TestBuildWhereMultipleGroups(t *testing.T) {
	// 组间 AND
	groups := []models.SearchGroup{
		{Conditions: []models.SearchCondition{
			{Field: "teacherName", MatchType: "contains", Value: "高珍"},
		}},
		{Conditions: []models.SearchCondition{
			{Field: "campus", MatchType: "contains", Value: "四平路校区"},
		}},
	}
	where, params := buildWhere(groups)
	assert.Equal(t,
		"(t.teacherName LIKE ?) AND (cp.campusI18n LIKE ?)",
		where,
	)
	assert.Equal(t, []interface{}{"%高珍%", "%四平路校区%"}, params)
}

func TestBuildWhereMixedGroup(t *testing.T) {
	// 组内 OR + 组间 AND
	groups := []models.SearchGroup{
		{Conditions: []models.SearchCondition{
			{Field: "courseName", MatchType: "contains", Value: "数学"},
			{Field: "courseCode", MatchType: "contains", Value: "002"},
		}},
		{Conditions: []models.SearchCondition{
			{Field: "campus", MatchType: "contains", Value: "四平路校区"},
		}},
	}
	where, params := buildWhere(groups)
	assert.True(t, strings.HasPrefix(where, "(cd.courseName LIKE ? OR cd.code LIKE ?) AND (cp.campusI18n LIKE ?)"))
	assert.Equal(t, []interface{}{"%数学%", "%002%", "%四平路校区%"}, params)
}

func TestBuildWhereSkipEmptyValue(t *testing.T) {
	groups := []models.SearchGroup{
		{Conditions: []models.SearchCondition{
			{Field: "teacherName", MatchType: "contains", Value: "高珍"},
			{Field: "courseName", MatchType: "contains", Value: ""}, // 跳过
		}},
	}
	where, params := buildWhere(groups)
	assert.Equal(t, "(t.teacherName LIKE ?)", where)
	assert.Len(t, params, 1)
}

func TestBuildWhereSkipUnknownField(t *testing.T) {
	groups := []models.SearchGroup{
		{Conditions: []models.SearchCondition{
			{Field: "nonexistent", MatchType: "contains", Value: "test"},
			{Field: "courseName", MatchType: "contains", Value: "数学"},
		}},
	}
	where, params := buildWhere(groups)
	assert.Equal(t, "(cd.courseName LIKE ?)", where)
	assert.Len(t, params, 1)
}

func TestBuildWhereSelectFields(t *testing.T) {
	// 验证 select 和 input 字段都正确映射
	tests := []struct {
		field      string
		expectCol  string
		searchType string
	}{
		{"courseCode", "cd.code", "input"},
		{"newCode", "cd.newCode", "input"},
		{"courseName", "cd.courseName", "input"},
		{"teacherName", "t.teacherName", "input"},
		{"teacherCode", "t.teacherCode", "input"},
		{"schedule", "t.arrangeInfoText", "input"},
		{"campus", "cp.campusI18n", "select"},
		{"faculty", "f.facultyI18n", "select"},
		{"courseType", "cn.courseLabelName", "select"},
		{"major", "m.name", "select"},
	}

	for _, tc := range tests {
		t.Run(tc.field, func(t *testing.T) {
			def, ok := models.FieldByID[tc.field]
			require.True(t, ok, "字段 %s 未定义", tc.field)
			assert.Equal(t, tc.expectCol, def.DBColumn)
			assert.Equal(t, tc.searchType, def.SearchType)
		})
	}
}

func TestBuildWhereAllFieldsReachable(t *testing.T) {
	// 每个搜索字段都能生成正确的 WHERE
	for _, def := range models.AllFields {
		groups := []models.SearchGroup{
			{Conditions: []models.SearchCondition{
				{Field: def.ID, MatchType: "contains", Value: "test"},
			}},
		}
		where, params := buildWhere(groups)
		assert.Contains(t, where, def.DBColumn)
		assert.Equal(t, "%test%", params[0])
	}
}

// ─── dedupLines 单元测试 ───

func TestDedupLinesEmpty(t *testing.T) {
	assert.Equal(t, "", dedupLines(""))
}

func TestDedupLinesSingleLine(t *testing.T) {
	line := "高珍(05119) 星期二3-4节 [4-10双] 复楼F102"
	assert.Equal(t, line, dedupLines(line))
}

func TestDedupLinesDuplicateLines(t *testing.T) {
	input := "line1\nline1\nline2\nline2\nline1"
	expected := "line1\nline2"
	assert.Equal(t, expected, dedupLines(input))
}

func TestDedupLinesBlankLines(t *testing.T) {
	input := "line1\n\n\nline2\n"
	expected := "line1\nline2"
	assert.Equal(t, expected, dedupLines(input))
}

func TestDedupLinesWhitespaceTrim(t *testing.T) {
	input := "  line1  \n  line1\n  line2  "
	expected := "line1\nline2"
	assert.Equal(t, expected, dedupLines(input))
}
