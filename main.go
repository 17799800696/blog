package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/test/blog/config"
	"github.com/test/blog/routes"
)

func main() {
	// 验证配置
	config.ValidateConfig()

	// 加载配置
	cfg := config.LoadConfig()

	// 打印配置信息
	config.PrintConfig(cfg)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	config.InitDB(cfg)

	// 创建Gin引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
