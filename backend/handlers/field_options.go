package handlers

import (
	"net/http"

	"jkd-website/backend/db"
	"jkd-website/backend/models"

	"github.com/gin-gonic/gin"
)

func GetFieldOptions(router *db.Router) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.FieldOptionsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Err(400, "请求格式错误"))
			return
		}

		options, err := router.GetFieldOptions(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Err(500, "查询字段选项失败: "+err.Error()))
			return
		}

		c.JSON(http.StatusOK, models.OK(options))
	}
}
