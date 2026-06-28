package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"jkd-website/backend/config"
	"jkd-website/backend/db"
	"jkd-website/backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	router, err := db.NewRouter(cfg)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer router.Close()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// RESTful 路由
	api := r.Group("/api")
	{
		// 可检索字段
		api.GET("/search-fields", handlers.GetSearchFields)
		api.POST("/search-fields/options", handlers.GetFieldOptions(router))

		// 学期列表
		api.GET("/calendars", handlers.GetCalendars(router))

		// 课程搜索
		api.POST("/courses/search", handlers.SearchCourses(router))

		// 最新更新时间
		api.GET("/last-update", handlers.GetLastUpdate(router))
	}

	// 健康检查
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "服务健康", "data": nil})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r.Handler(),
	}

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("正在关闭服务...")
		srv.Close()
	}()

	log.Printf("ics-website 后端启动于 :%s", cfg.ServerPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("服务异常退出: %v", err)
	}
}
