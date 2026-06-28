package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jkd-website/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSearchFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/search-fields", GetSearchFields)

	req := httptest.NewRequest("GET", "/api/search-fields", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 200, resp.Code)

	fields, ok := resp.Data.([]interface{})
	require.True(t, ok)
	assert.Len(t, fields, len(models.AllFields))

	// 每个字段都有 id/label/searchType
	for _, f := range fields {
		m := f.(map[string]interface{})
		assert.NotEmpty(t, m["id"])
		assert.NotEmpty(t, m["label"])
		assert.Contains(t, []interface{}{"select", "input"}, m["searchType"])
	}
}
