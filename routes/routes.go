package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/test/blog/handlers"
	"github.com/test/blog/middleware"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// API路由组
	api := r.Group("/api")
	{
		// 认证相关路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", handlers.GetProfile)
			// 后续会添加文章和评论相关的路由
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog API is running",
		})
	})
}
