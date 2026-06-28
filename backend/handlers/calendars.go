package handlers

import (
	"net/http"

	"jkd-website/backend/db"
	"jkd-website/backend/models"

	"github.com/gin-gonic/gin"
)

func GetCalendars(router *db.Router) gin.HandlerFunc {
	return func(c *gin.Context) {
		cals, err := router.Calendars()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Err(500, "查询学期失败: "+err.Error()))
			return
		}
		result := make([]models.CalendarInfo, len(cals))
		for i, cal := range cals {
			result[i] = models.CalendarInfo{
				CalendarID:   cal.CalendarID,
				CalendarName: cal.CalendarName,
			}
		}
		c.JSON(http.StatusOK, models.OK(result))
	}
}
