package utils

import (
	"net/http"

	"go.uber.org/zap"
)

// CustomError 自定义错误类型
type CustomError struct {
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Success   bool   `json:"success"`
	ErrorType string `json:"error_type,omitempty"`
}

// Error 实现error接口
func (e CustomError) Error() string {
	return e.Message
}

// NewError 创建新的错误
func NewError(message string, code int) CustomError {
	return CustomError{
		Message: message,
		Code:    code,
		Success: false,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(message string) CustomError {
	return CustomError{
		Message:   message,
		Code:      http.StatusBadRequest,
		Success:   false,
		ErrorType: "validation_error",
	}
}

// NewAuthError 创建认证错误
func NewAuthError(message string) CustomError {
	return CustomError{
		Message:   message,
		Code:      http.StatusUnauthorized,
		Success:   false,
		ErrorType: "auth_error",
	}
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(message string) CustomError {
	return CustomError{
		Message:   message,
		Code:      http.StatusNotFound,
		Success:   false,
		ErrorType: "not_found_error",
	}
}

// NewInternalError 创建内部错误
func NewInternalError(message string) CustomError {
	return CustomError{
		Message:   message,
		Code:      http.StatusInternalServerError,
		Success:   false,
		ErrorType: "internal_error",
	}
}

// LogInfo 记录信息日志
func LogInfo(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(message, fields...)
	}
}

// LogError 记录错误日志
func LogError(message string, err error, fields ...zap.Field) {
	if Logger != nil {
		if err != nil {
			fields = append(fields, zap.Error(err))
		}
		Logger.Error(message, fields...)
	}
}

// LogWarn 记录警告日志
func LogWarn(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(message, fields...)
	}
}

// LogDebug 记录调试日志
func LogDebug(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(message, fields...)
	}
}
