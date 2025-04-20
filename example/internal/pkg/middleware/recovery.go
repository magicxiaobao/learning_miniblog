package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"example/internal/pkg/errno"
	"example/internal/pkg/known"
	"example/internal/pkg/log"
)

// Recovery 返回一个中间件，该中间件可以从任何 panic 恢复，并写入 500 响应。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查连接是否已断开
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				// 获取请求ID
				requestID := c.GetString(known.XRequestIDKey)

				// 如果连接已断开，不继续尝试写入响应
				if brokenPipe {
					log.Errorw("Recovery from brokenPipe",
						"request_id", requestID,
						"error", err,
						"request", string(httpRequest),
					)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				// 记录详细错误日志
				log.Errorw("Recovery from panic",
					"request_id", requestID,
					"error", err,
					"request", string(httpRequest),
					"stack", string(debug.Stack()),
				)

				// 返回 500 错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    errno.InternalServerError.Code,
					"message": "服务器内部错误，请稍后再试",
				})
			}
		}()

		c.Next()
	}
}

// TimeoutMiddleware 超时中间件，如果处理时间超过指定时间，则返回超时错误
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用channel完成同步
		finish := make(chan struct{}, 1)
		// 使用 panicChan 捕获处理过程中的 panic
		panicChan := make(chan interface{}, 1)

		// 在 goroutine 中处理请求
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			// 处理请求
			c.Next()
			finish <- struct{}{}
		}()

		// 等待请求处理完成或超时
		select {
		case p := <-panicChan:
			// 处理过程中出现了 panic，重新抛出以便让 Recovery 中间件处理
			panic(p)

		case <-finish:
			// 请求正常完成，不做任何操作

		case <-time.After(timeout):
			// 请求处理超时，返回超时错误
			requestID := c.GetString(known.XRequestIDKey)
			log.Warnw("Request processing timed out",
				"request_id", requestID,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"timeout", fmt.Sprintf("%v", timeout),
			)

			// 返回超时错误响应
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"code":    504,
				"message": "请求处理超时，请稍后再试",
			})
		}
	}
}
