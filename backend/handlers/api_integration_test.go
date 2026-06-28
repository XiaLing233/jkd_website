//go:build integration

package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"jkd-website/backend/config"
	"jkd-website/backend/db"
	"jkd-website/backend/db/testutil"
	"jkd-website/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAPI(t *testing.T) (*gin.Engine, *db.Router) {
	t.Helper()

	cfg := testutil.DefaultCfg()
	require.NoError(t, testutil.Setup(cfg), "init test DB")
	testutil.SetEnv(cfg)

	c, err := config.Load()
	require.NoError(t, err)

	router, err := db.NewRouter(c)
	require.NoError(t, err)
	t.Cleanup(func() { router.Close() })

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/search-fields", GetSearchFields)
	r.POST("/api/search-fields/options", GetFieldOptions(router))
	r.GET("/api/calendars", GetCalendars(router))
	r.POST("/api/courses/search", SearchCourses(router))
	r.GET("/api/last-update", GetLastUpdate(router))

	return r, router
}

func TestAPI_GetSearchFields(t *testing.T) {
	r, _ := setupAPI(t)
	req := httptest.NewRequest("GET", "/api/search-fields", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fields := resp.Data.([]interface{})
	assert.Len(t, fields, len(models.AllFields))
}

func TestAPI_GetCalendars(t *testing.T) {
	r, _ := setupAPI(t)
	req := httptest.NewRequest("GET", "/api/calendars", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	cals := resp.Data.([]interface{})
	assert.NotEmpty(t, cals)
}

func TestAPI_GetCalendarsExcludesSyncing(t *testing.T) {
	r, router := setupAPI(t)

	// 插入一个正在同步中的学期（calendarIdI18n = '数据同步中…'）
	_, err := router.Meta().Exec(
		"INSERT IGNORE INTO calendar_registry (calendarId, calendarIdI18n) VALUES (?, '数据同步中…')",
		888,
	)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/calendars", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	cals := resp.Data.([]interface{})

	hasNormal := false
	for _, cal := range cals {
		cm := cal.(map[string]interface{})
		name := cm["calendarName"].(string)
		assert.NotEqual(t, "数据同步中…", name,
			"同步中的学期不应出现在日历列表中")
		if name == "测试学期" {
			hasNormal = true
		}
	}
	assert.True(t, hasNormal, "正常学期仍应返回")
}

func TestAPI_SearchCoursesValid(t *testing.T) {
	r, _ := setupAPI(t)
	body := `{"groups":[],"calendar_ids":[999],"page":1,"page_size":5}`
	req := httptest.NewRequest("POST", "/api/courses/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 200, resp.Code)

	pd := resp.Data.(map[string]interface{})
	assert.NotEmpty(t, pd["items"])
}

func TestAPI_SearchByTeacher(t *testing.T) {
	r, _ := setupAPI(t)
	body := `{"groups":[{"conditions":[{"field":"teacherName","matchType":"contains","value":"关佶红"}]}],"calendar_ids":[999],"page":1,"page_size":5}`
	req := httptest.NewRequest("POST", "/api/courses/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	pd := resp.Data.(map[string]interface{})
	items := pd["items"].([]interface{})
	require.NotEmpty(t, items)
	first := items[0].(map[string]interface{})
	assert.Contains(t, first["teachers"], "关佶红")
}

func TestAPI_FieldOptions(t *testing.T) {
	r, _ := setupAPI(t)
	body := `{"field":"campus","calendar_ids":[999]}`
	req := httptest.NewRequest("POST", "/api/search-fields/options", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	opts := resp.Data.([]interface{})
	assert.NotEmpty(t, opts)
}

func TestAPI_GetLastUpdate(t *testing.T) {
	r, _ := setupAPI(t)
	req := httptest.NewRequest("GET", "/api/last-update", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAPI_SearchMissingCalendarIDs(t *testing.T) {
	r, _ := setupAPI(t)
	body := `{"groups":[],"calendar_ids":[],"page":1,"page_size":5}`
	req := httptest.NewRequest("POST", "/api/courses/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestAPI_SearchTooManyCalendarsNoCondition(t *testing.T) {
	r, _ := setupAPI(t)
	body := `{"groups":[],"calendar_ids":[111,222,333],"page":1,"page_size":5}`
	req := httptest.NewRequest("POST", "/api/courses/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
