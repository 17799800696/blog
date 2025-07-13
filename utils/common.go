package utils

import (
	"fmt"
	"os"
	"strconv"
)

// GetEnv 获取环境变量，如果不存在则panic
func GetEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	panic(fmt.Sprintf("Environment variable %s is required but not set", key))
}

// GetEnvWithDefault 获取环境变量，如果不存在则返回默认值
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt 获取整数类型的环境变量，如果不存在则panic
func GetEnvInt(key string) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		panic(fmt.Sprintf("Environment variable %s must be a valid integer", key))
	}
	panic(fmt.Sprintf("Environment variable %s is required but not set", key))
}

// GetEnvIntWithDefault 获取整数类型的环境变量，如果不存在则返回默认值
func GetEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		panic(fmt.Sprintf("Environment variable %s must be a valid integer", key))
	}
	return defaultValue
}

// GetEnvBool 获取布尔类型的环境变量，如果不存在则panic
func GetEnvBool(key string) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		panic(fmt.Sprintf("Environment variable %s must be a valid boolean", key))
	}
	panic(fmt.Sprintf("Environment variable %s is required but not set", key))
}

// GetEnvBoolWithDefault 获取布尔类型的环境变量，如果不存在则返回默认值
func GetEnvBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		panic(fmt.Sprintf("Environment variable %s must be a valid boolean", key))
	}
	return defaultValue
}
