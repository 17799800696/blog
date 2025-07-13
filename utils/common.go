package utils

import (
	"fmt"
	"os"
	"strconv"
)

// GetEnv 获取环境变量，如果不存在则返回错误
func GetEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s is required but not set", key)
}

// GetEnvWithDefault 获取环境变量，如果不存在则返回默认值
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt 获取整数类型的环境变量，如果不存在或无效则返回错误
func GetEnvInt(key string) (int, error) {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue, nil
		}
		return 0, fmt.Errorf("environment variable %s must be a valid integer", key)
	}
	return 0, fmt.Errorf("environment variable %s is required but not set", key)
}

// GetEnvIntWithDefault 获取整数类型的环境变量，如果不存在则返回默认值
func GetEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvBool 获取布尔类型的环境变量，如果不存在或无效则返回错误
func GetEnvBool(key string) (bool, error) {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue, nil
		}
		return false, fmt.Errorf("environment variable %s must be a valid boolean", key)
	}
	return false, fmt.Errorf("environment variable %s is required but not set", key)
}

// GetEnvBoolWithDefault 获取布尔类型的环境变量，如果不存在则返回默认值
func GetEnvBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// MustGetEnv 获取环境变量，如果不存在则panic（用于向后兼容）
func MustGetEnv(key string) string {
	value, err := GetEnv(key)
	if err != nil {
		panic(err.Error())
	}
	return value
}

// MustGetEnvInt 获取整数类型的环境变量，如果不存在或无效则panic（用于向后兼容）
func MustGetEnvInt(key string) int {
	value, err := GetEnvInt(key)
	if err != nil {
		panic(err.Error())
	}
	return value
}

// MustGetEnvBool 获取布尔类型的环境变量，如果不存在或无效则panic（用于向后兼容）
func MustGetEnvBool(key string) bool {
	value, err := GetEnvBool(key)
	if err != nil {
		panic(err.Error())
	}
	return value
}
