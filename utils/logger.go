package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger 初始化日志系统
func InitLogger() {
	// 获取日志级别
	logLevel := GetEnvWithDefault("LOG_LEVEL", "info")
	level := getLogLevel(logLevel)

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.MessageKey = "msg"
	encoderConfig.LevelKey = "level"
	encoderConfig.CallerKey = "caller"
	encoderConfig.StacktraceKey = "stacktrace"

	// 配置输出
	var core zapcore.Core
	logFormat := GetEnvWithDefault("LOG_FORMAT", "json")
	logOutputPath := GetEnvWithDefault("LOG_OUTPUT_PATH", "")

	if logFormat == "json" {
		encoder := zapcore.NewJSONEncoder(encoderConfig)
		if logOutputPath != "" {
			// 输出到文件
			file, err := os.OpenFile(logOutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic("Failed to open log file: " + err.Error())
			}
			core = zapcore.NewCore(encoder, zapcore.AddSync(file), level)
		} else {
			// 输出到控制台
			core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
		}
	} else {
		// 控制台格式
		encoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	}

	// 创建logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// getLogLevel 获取日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// WithRequestID 添加请求ID到日志字段
func WithRequestID(requestID string) zap.Field {
	return zap.String("request_id", requestID)
}

// WithUserID 添加用户ID到日志字段
func WithUserID(userID uint) zap.Field {
	return zap.Uint("user_id", userID)
}

// WithPostID 添加文章ID到日志字段
func WithPostID(postID uint) zap.Field {
	return zap.Uint("post_id", postID)
}

// WithCommentID 添加评论ID到日志字段
func WithCommentID(commentID uint) zap.Field {
	return zap.Uint("comment_id", commentID)
}

// WithMethod 添加HTTP方法到日志字段
func WithMethod(method string) zap.Field {
	return zap.String("method", method)
}

// WithPath 添加请求路径到日志字段
func WithPath(path string) zap.Field {
	return zap.String("path", path)
}

// WithStatusCode 添加HTTP状态码到日志字段
func WithStatusCode(statusCode int) zap.Field {
	return zap.Int("status_code", statusCode)
}

// WithDuration 添加请求持续时间到日志字段
func WithDuration(duration float64) zap.Field {
	return zap.Float64("duration_ms", duration)
}
