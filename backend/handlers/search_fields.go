package handlers

import (
	"net/http"

	"jkd-website/backend/models"

	"github.com/gin-gonic/gin"
)

func GetSearchFields(c *gin.Context) {
	c.JSON(http.StatusOK, models.OK(models.AllFields))
}
