package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllFieldsDistinct(t *testing.T) {
	seen := make(map[string]bool)
	for _, f := range AllFields {
		assert.False(t, seen[f.ID], "字段 ID 重复: %s", f.ID)
		seen[f.ID] = true
	}
}

func TestFieldByIDComplete(t *testing.T) {
	assert.Equal(t, len(AllFields), len(FieldByID),
		"FieldByID 应与 AllFields 长度一致")
	for _, f := range AllFields {
		def, ok := FieldByID[f.ID]
		assert.True(t, ok, "FieldByID 缺少 %s", f.ID)
		assert.Equal(t, f, def)
	}
}

func TestOKResponse(t *testing.T) {
	resp := OK("test-data")
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "查询成功", resp.Msg)
	assert.Equal(t, "test-data", resp.Data)
}

func TestErrResponse(t *testing.T) {
	resp := Err(400, "参数错误")
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "参数错误", resp.Msg)
	assert.Nil(t, resp.Data)
}

func TestOKPaginatedResponse(t *testing.T) {
	data := []string{"a", "b"}
	resp := OKPaginated(data, 1, 2, 100)
	assert.Equal(t, 200, resp.Code)

	pd, ok := resp.Data.(PaginatedData)
	assert.True(t, ok)
	assert.Equal(t, 1, pd.Page)
	assert.Equal(t, 2, pd.PageSize)
	assert.Equal(t, 100, pd.TotalCount)
	assert.Len(t, pd.Items, 2)
}

func TestSearchResultStruct(t *testing.T) {
	// 确保 NewCode 是指针类型（NULL 安全）
	r := SearchResult{}
	assert.Nil(t, r.NewCode)
	assert.Nil(t, r.TotalHours)
	assert.Nil(t, r.Capacity)
	assert.Nil(t, r.Enrolled)
}
