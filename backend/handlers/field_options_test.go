package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFieldOptionsInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/search-fields/options", GetFieldOptions(nil))

	req := httptest.NewRequest("POST", "/api/search-fields/options", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// ShouldBindJSON 失败 → 返回 400，不访问 router（传 nil 安全）
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
