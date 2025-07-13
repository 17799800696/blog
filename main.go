package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/test/blog/config"
	"github.com/test/blog/routes"
	"github.com/test/blog/utils"
	"go.uber.org/zap"
)

func main() {
	// 验证配置
	config.ValidateConfig()

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化日志系统
	utils.InitLogger()

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

	// 创建HTTP服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Printf("Server starting on %s", serverAddr)
		utils.LogInfo("Server starting", zap.String("port", cfg.Server.Port))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	log.Println("Shutting down server...")
	utils.LogInfo("Server shutting down")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// 关闭数据库连接
	if db := config.GetDB(); db != nil {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}
		}
	}

	log.Println("Server exited")
	utils.LogInfo("Server exited successfully")
}
