//go:build integration

package db

import (
	"testing"

	"jkd-website/backend/config"
	"jkd-website/backend/db/testutil"
	"jkd-website/backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) *Router {
	t.Helper()
	cfg := testutil.DefaultCfg()
	require.NoError(t, testutil.Setup(cfg), "init test DB")
	testutil.SetEnv(cfg)

	c, err := config.Load()
	require.NoError(t, err)
	r, err := NewRouter(c)
	require.NoError(t, err)
	t.Cleanup(func() { r.Close() })
	return r
}

func TestDB_Calendars(t *testing.T) {
	r := setupDB(t)
	cals, err := r.Calendars()
	require.NoError(t, err)
	require.NotEmpty(t, cals)
	assert.Equal(t, testutil.TestCalID, cals[0].CalendarID)
}

func TestDB_CalendarsRefreshAfterInsert(t *testing.T) {
	r := setupDB(t)

	// 启动后插入新学期
	_, err := r.Meta().Exec(
		"INSERT INTO calendar_registry (calendarId, calendarIdI18n) VALUES (998, '新测试学期')")
	require.NoError(t, err)
	defer r.Meta().Exec("DELETE FROM calendar_registry WHERE calendarId = 998")

	// 查询应能立即看到新加学期
	cals, err := r.Calendars()
	require.NoError(t, err)

	found := false
	for _, c := range cals {
		if c.CalendarID == 998 {
			found = true
			break
		}
	}
	assert.True(t, found, "启动后新增的学期应立即出现在列表中")
}

func TestDB_SearchAll(t *testing.T) {
	r := setupDB(t)
	results, err := r.Search(models.SearchRequest{
		CalendarIDs: []int{testutil.TestCalID},
		Page:        1, PageSize: 5,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, results)
	assert.NotEmpty(t, results[0].CourseCode)
}

func TestDB_SearchByTeacher(t *testing.T) {
	r := setupDB(t)
	results, err := r.Search(models.SearchRequest{
		CalendarIDs: []int{testutil.TestCalID},
		Groups: []models.SearchGroup{
			{Conditions: []models.SearchCondition{
				{Field: "teacherName", MatchType: "contains", Value: "关佶红"},
			}},
		},
		Page: 1, PageSize: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, results)
	assert.Contains(t, results[0].Teachers, "关佶红")
}

func TestDB_SearchByCourseCode(t *testing.T) {
	r := setupDB(t)
	results, err := r.Search(models.SearchRequest{
		CalendarIDs: []int{testutil.TestCalID},
		Groups: []models.SearchGroup{
			{Conditions: []models.SearchCondition{
				{Field: "courseCode", MatchType: "contains", Value: "140076"},
			}},
		},
		Page: 1, PageSize: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, results)
	assert.Equal(t, "公共营养学", results[0].CourseName)
}

func TestDB_FieldOptions(t *testing.T) {
	r := setupDB(t)
	opts, err := r.GetFieldOptions(models.FieldOptionsRequest{
		Field:       "campus",
		CalendarIDs: []int{testutil.TestCalID},
	})
	require.NoError(t, err)
	values := make([]string, len(opts))
	for i, o := range opts {
		values[i] = o.Value
	}
	assert.Contains(t, values, "四平路校区")
}

func TestDB_Pagination(t *testing.T) {
	r := setupDB(t)
	page1, _ := r.Search(models.SearchRequest{
		CalendarIDs: []int{testutil.TestCalID},
		Page:        1, PageSize: 1,
	})
	require.NotEmpty(t, page1)
}
