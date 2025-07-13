package config

import (
	"github.com/test/blog/utils"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: utils.GetEnv("SERVER_PORT"),
			Mode: utils.GetEnv("GIN_MODE"),
		},
		Database: DatabaseConfig{
			Host:     utils.GetEnv("DB_HOST"),
			Port:     utils.GetEnv("DB_PORT"),
			User:     utils.GetEnv("DB_USER"),
			Password: utils.GetEnv("DB_PASSWORD"),
			Database: utils.GetEnv("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret:          utils.GetEnv("JWT_SECRET"),
			ExpirationHours: utils.GetEnvInt("JWT_EXPIRATION_HOURS"),
		},
	}
}
