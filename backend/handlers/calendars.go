package handlers

import (
	"net/http"

	"jkd-website/backend/db"
	"jkd-website/backend/models"

	"github.com/gin-gonic/gin"
)

func GetCalendars(router *db.Router) gin.HandlerFunc {
	return func(c *gin.Context) {
		calendars := router.Calendars()
		result := make([]models.CalendarInfo, len(calendars))
		for i, cal := range calendars {
			result[i] = models.CalendarInfo{
				CalendarID:   cal.CalendarID,
				CalendarName: cal.CalendarName,
			}
		}
		c.JSON(http.StatusOK, models.OK(result))
	}
}
