package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/test/blog/utils"
)

const (
	RequestIDKey    = "request_id"
	RequestIDHeader = "X-Request-ID"
)

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取请求ID，如果没有则生成新的
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 设置请求ID到上下文
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDHeader, requestID)

		// 记录请求开始
		utils.LogInfo("Request started",
			utils.WithRequestID(requestID),
			utils.WithMethod(c.Request.Method),
			utils.WithPath(c.Request.URL.Path),
		)

		// 处理请求
		c.Next()

		// 记录请求结束
		utils.LogInfo("Request completed",
			utils.WithRequestID(requestID),
			utils.WithMethod(c.Request.Method),
			utils.WithPath(c.Request.URL.Path),
			utils.WithStatusCode(c.Writer.Status()),
		)
	}
}

// GetRequestID 从上下文获取请求ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
