package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"jkd-website/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 假的 Router，只用于验证请求格式错误检测
type fakeRouter struct{}

func (f *fakeRouter) Search(req models.SearchRequest) ([]models.SearchResult, error) {
	return nil, nil
}

func TestSearchMissingCalendarIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/courses/search", SearchCourses(nil))

	body := `{"groups":[],"calendar_ids":[]}`
	req := httptest.NewRequest("POST", "/api/courses/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 400, resp.Code)
}

func TestSearchTooManyCalendarsNoCondition(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/courses/search", SearchCourses(nil))

	body := `{"groups":[],"calendar_ids":[117,118,119]}`
	req := httptest.NewRequest("POST", "/api/courses/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/courses/search", SearchCourses(nil))

	req := httptest.NewRequest("POST", "/api/courses/search", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
