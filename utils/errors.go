package utils

import (
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// 自定义错误类型
var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrBadRequest        = errors.New("bad request")
	ErrInternalServer    = errors.New("internal server error")
	ErrDatabaseError     = errors.New("database error")
	ErrValidationError   = errors.New("validation error")
	ErrDuplicateResource = errors.New("resource already exists")
)

// AppError 应用错误结构
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// NewAppError 创建新的应用错误
func NewAppError(code int, message string, err error) *AppError {
	appErr := &AppError{
		Code:    code,
		Message: message,
	}

	if err != nil {
		appErr.Error = err.Error()
	}

	return appErr
}

// GetHTTPStatus 根据错误类型获取HTTP状态码
func GetHTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrValidationError):
		return http.StatusBadRequest
	case errors.Is(err, ErrDuplicateResource):
		return http.StatusConflict
	case errors.Is(err, ErrDatabaseError):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// LogError 记录错误日志
func LogError(err error, context string) {
	logger := GetLogger()
	logger.Error("application error",
		zap.Error(err),
		zap.String("context", context),
	)
}

// LogInfo 记录信息日志
func LogInfo(message string, fields ...zap.Field) {
	logger := GetLogger()
	logger.Info(message, fields...)
}

// LogWarn 记录警告日志
func LogWarn(message string, fields ...zap.Field) {
	logger := GetLogger()
	logger.Warn(message, fields...)
}

// LogDebug 记录调试日志
func LogDebug(message string, fields ...zap.Field) {
	logger := GetLogger()
	logger.Debug(message, fields...)
}
