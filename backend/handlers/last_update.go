package handlers

import (
	"net/http"

	"jkd-website/backend/db"
	"jkd-website/backend/models"

	"github.com/gin-gonic/gin"
)

func GetLastUpdate(router *db.Router) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, msg, err := router.GetLatestUpdate()
		if err != nil || t == nil {
			c.JSON(http.StatusOK, models.OK(map[string]string{
				"fetchTime": "未知",
				"message":   "暂无更新记录",
			}))
			return
		}
		data := map[string]string{
			"fetchTime": t.Format("2006-01-02"),
		}
		if msg != "" {
			data["message"] = msg
		} else {
			data["message"] = "数据已更新"
		}
		c.JSON(http.StatusOK, models.OK(data))
	}
}
