package handlers

import (
	"log"
	"net/http"

	"jkd-website/backend/db"
	"jkd-website/backend/models"

	"github.com/gin-gonic/gin"
)

func SearchCourses(router *db.Router) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SearchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Err(400, "请求格式错误"))
			return
		}

		if len(req.CalendarIDs) == 0 {
			c.JSON(http.StatusBadRequest, models.Err(400, "请选择至少一个学期"))
			return
		}

		hasConditions := false
		for _, g := range req.Groups {
			if len(g.Conditions) > 0 {
				hasConditions = true
				break
			}
		}
		if !hasConditions && len(req.CalendarIDs) > 2 {
			c.JSON(http.StatusBadRequest, models.Err(400, "至少选择1个检索条件！不允许在不选择条件的情况下查看超过两个学期的课程"))
			return
		}

		clientIP := c.ClientIP()
		log.Printf("[检索] IP=%s calendar_ids=%v groups=%d", clientIP, req.CalendarIDs, len(req.Groups))

		results, err := router.Search(req)
		if err != nil {
			log.Printf("[检索] 错误: %v", err)
			c.JSON(http.StatusInternalServerError, models.Err(500, "检索出错: "+err.Error()))
			return
		}

		// 默认分页
		page := req.Page
		pageSize := req.PageSize
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 50
		}
		if pageSize > 100 {
			pageSize = 100
		}

		total := len(results)
		start := (page - 1) * pageSize
		if start > total {
			start = total
		}
		end := start + pageSize
		if end > total {
			end = total
		}

		c.JSON(http.StatusOK, models.OKPaginated(results[start:end], page, pageSize, total))
	}
}
