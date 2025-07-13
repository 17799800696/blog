package config

import (
	"log"
	"os"
	"strings"

	"github.com/test/blog/utils"
)

// ValidateConfig 验证配置
func ValidateConfig() {
	log.Println("Validating configuration...")

	// 验证必需的环境变量
	requiredEnvVars := []string{
		"JWT_SECRET", // 只有JWT_SECRET是真正必需的
	}

	var missingVars []string

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		log.Fatalf("Missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	// 验证JWT过期时间
	expirationHours := utils.GetEnvIntWithDefault("JWT_EXPIRATION_HOURS", 24)
	if expirationHours <= 0 {
		log.Fatal("JWT_EXPIRATION_HOURS must be greater than 0")
	}

	// 验证服务器端口
	port := utils.GetEnvWithDefault("SERVER_PORT", "8080")
	if port == "" {
		log.Fatal("SERVER_PORT cannot be empty")
	}

	// 验证Gin模式
	mode := utils.GetEnvWithDefault("GIN_MODE", "debug")
	validModes := []string{"debug", "release", "test"}
	isValidMode := false
	for _, validMode := range validModes {
		if mode == validMode {
			isValidMode = true
			break
		}
	}
	if !isValidMode {
		log.Fatalf("GIN_MODE must be one of: %s", strings.Join(validModes, ", "))
	}

	// 验证数据库配置
	dbHost := utils.GetEnvWithDefault("DB_HOST", "localhost")
	dbPort := utils.GetEnvWithDefault("DB_PORT", "3306")
	dbUser := utils.GetEnvWithDefault("DB_USER", "root")
	dbName := utils.GetEnvWithDefault("DB_NAME", "blog")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatal("Database configuration cannot have empty values")
	}

	// 验证日志级别
	logLevel := utils.GetEnvWithDefault("LOG_LEVEL", "info")
	validLogLevels := []string{"debug", "info", "warn", "error"}
	isValidLogLevel := false
	for _, validLevel := range validLogLevels {
		if logLevel == validLevel {
			isValidLogLevel = true
			break
		}
	}
	if !isValidLogLevel {
		log.Fatalf("LOG_LEVEL must be one of: %s", strings.Join(validLogLevels, ", "))
	}

	// 验证数据库连接池配置
	maxIdleConns := utils.GetEnvIntWithDefault("DB_MAX_IDLE_CONNS", 10)
	maxOpenConns := utils.GetEnvIntWithDefault("DB_MAX_OPEN_CONNS", 100)
	connMaxLifetime := utils.GetEnvIntWithDefault("DB_CONN_MAX_LIFETIME", 60)

	if maxIdleConns <= 0 || maxOpenConns <= 0 || connMaxLifetime <= 0 {
		log.Fatal("Database connection pool settings must be greater than 0")
	}

	if maxIdleConns > maxOpenConns {
		log.Fatal("DB_MAX_IDLE_CONNS cannot be greater than DB_MAX_OPEN_CONNS")
	}

	log.Println("Configuration validation passed!")
}

// PrintConfig 打印当前配置（隐藏敏感信息）
func PrintConfig(cfg *Config) {
	log.Println("Current configuration:")
	log.Printf("  Server Port: %s", cfg.Server.Port)
	log.Printf("  Gin Mode: %s", cfg.Server.Mode)
	log.Printf("  Database Host: %s", cfg.Database.Host)
	log.Printf("  Database Port: %s", cfg.Database.Port)
	log.Printf("  Database User: %s", cfg.Database.User)
	log.Printf("  Database Name: %s", cfg.Database.Database)
	log.Printf("  Database Max Idle Conns: %d", cfg.Database.MaxIdleConns)
	log.Printf("  Database Max Open Conns: %d", cfg.Database.MaxOpenConns)
	log.Printf("  Database Conn Max Lifetime: %d minutes", cfg.Database.ConnMaxLifetime)
	log.Printf("  JWT Expiration Hours: %d", cfg.JWT.ExpirationHours)
	log.Printf("  JWT Secret: %s", maskSecret(cfg.JWT.Secret))
	log.Printf("  Log Level: %s", cfg.Log.Level)
	log.Printf("  Log Format: %s", cfg.Log.Format)
}

// maskSecret 隐藏敏感信息
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "..." + secret[len(secret)-4:]
}
