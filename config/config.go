package config

import (
	"github.com/test/blog/utils"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int // 分钟
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string
	Format     string
	OutputPath string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: utils.GetEnvWithDefault("SERVER_PORT", "8080"),
			Mode: utils.GetEnvWithDefault("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:            utils.GetEnvWithDefault("DB_HOST", "localhost"),
			Port:            utils.GetEnvWithDefault("DB_PORT", "3306"),
			User:            utils.GetEnvWithDefault("DB_USER", "root"),
			Password:        utils.GetEnvWithDefault("DB_PASSWORD", ""),
			Database:        utils.GetEnvWithDefault("DB_NAME", "blog"),
			MaxIdleConns:    utils.GetEnvIntWithDefault("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    utils.GetEnvIntWithDefault("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: utils.GetEnvIntWithDefault("DB_CONN_MAX_LIFETIME", 60),
		},
		JWT: JWTConfig{
			Secret:          utils.GetEnvWithDefault("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpirationHours: utils.GetEnvIntWithDefault("JWT_EXPIRATION_HOURS", 24),
		},
		Log: LogConfig{
			Level:      utils.GetEnvWithDefault("LOG_LEVEL", "info"),
			Format:     utils.GetEnvWithDefault("LOG_FORMAT", "json"),
			OutputPath: utils.GetEnvWithDefault("LOG_OUTPUT_PATH", ""),
		},
	}
}
