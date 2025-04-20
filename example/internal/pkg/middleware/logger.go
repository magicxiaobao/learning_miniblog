package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"example/internal/pkg/known"
	"example/internal/pkg/log"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写 ResponseWriter 的 Write 方法
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Logger 返回 gin 中间件，用于记录 API 请求和响应的详细信息
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装 ResponseWriter，以便我们可以捕获写入的内容
		w := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)
		// 获取 RequestID
		requestID := c.GetString(known.XRequestIDKey)
		// 获取状态码
		statusCode := c.Writer.Status()
		// 获取客户端 IP
		clientIP := c.ClientIP()
		// 获取用户代理
		userAgent := c.Request.UserAgent()
		// 响应大小
		responseSize := c.Writer.Size()

		// 构建 URL
		url := path
		if query != "" {
			url = path + "?" + query
		}

		// 日志字段
		fields := []interface{}{
			"status", statusCode,
			"latency", latency,
			"ip", clientIP,
			"method", method,
			"uri", url,
			"size", responseSize,
			"user-agent", userAgent,
		}

		// 如果有请求ID，则添加到日志字段中
		if requestID != "" {
			fields = append(fields, "request-id", requestID)
		}

		// 根据状态码确定日志级别
		switch {
		case statusCode >= 500:
			// 添加请求和响应体到日志中，因为这是一个错误
			fields = append(fields, "req", string(requestBody), "resp", w.body.String())
			log.Errorw("HTTP API Request", fields...)
		case statusCode >= 400:
			// 添加请求体到日志中，因为这是一个客户端错误
			fields = append(fields, "req", string(requestBody))
			log.Warnw("HTTP API Request", fields...)
		case statusCode >= 200 && statusCode < 300:
			// 正常响应
			log.Infow("HTTP API Request", fields...)
		default:
			// 其他情况
			log.Infow("HTTP API Request", fields...)
		}
	}
}
