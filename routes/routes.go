package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/test/blog/handlers"
	"github.com/test/blog/middleware"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 添加请求ID中间件
	r.Use(middleware.RequestID())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Blog API is running",
			"status":  "ok",
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			authorized.GET("/profile", handlers.GetProfile)
			authorized.POST("/posts", handlers.CreatePost)
			authorized.PUT("/posts/:id", handlers.UpdatePost)
			authorized.DELETE("/posts/:id", handlers.DeletePost)
			authorized.POST("/posts/:id/comments", handlers.CreateComment)
		}

		// 公开路由
		api.GET("/posts", handlers.GetPosts)
		api.GET("/posts/:id", handlers.GetPost)
		api.GET("/posts/:id/comments", handlers.GetComments)
	}
}
