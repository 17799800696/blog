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

	requiredEnvVars := []string{
		"SERVER_PORT",
		"GIN_MODE",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"JWT_SECRET",
		"JWT_EXPIRATION_HOURS",
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
	expirationHours := utils.GetEnvInt("JWT_EXPIRATION_HOURS")
	if expirationHours <= 0 {
		log.Fatal("JWT_EXPIRATION_HOURS must be greater than 0")
	}

	// 验证服务器端口
	port := utils.GetEnv("SERVER_PORT")
	if port == "" {
		log.Fatal("SERVER_PORT is required")
	}

	// 验证Gin模式
	mode := utils.GetEnv("GIN_MODE")
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
	log.Printf("  JWT Expiration Hours: %d", cfg.JWT.ExpirationHours)
	log.Printf("  JWT Secret: %s", maskSecret(cfg.JWT.Secret))
}

// maskSecret 隐藏敏感信息
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "..." + secret[len(secret)-4:]
}
