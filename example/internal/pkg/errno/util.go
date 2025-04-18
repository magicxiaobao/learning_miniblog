package errno

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 用于包装API响应
type Response struct {
	// 业务错误码
	Code int `json:"code"`

	// 业务错误消息
	Message string `json:"message"`

	// 数据负载
	Data interface{} `json:"data,omitempty"`
}

// WriteResponse 将标准的API响应写入 HTTP 响应体
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		// 获取错误码和消息
		code, message := Decode(err)

		// 检查是否需要记录错误日志
		logError(c, err)

		// 设置响应状态码
		httpCode := http.StatusOK
		if e, ok := err.(*Errno); ok {
			httpCode = e.ToHTTPStatusCode()
		}

		c.JSON(httpCode, &Response{
			Code:    code,
			Message: message,
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Code:    OK.Code,
		Message: OK.Message,
		Data:    data,
	})
}

// logError 根据错误类型决定是否记录错误日志
func logError(c *gin.Context, err error) {
	// 检查是否需要记录错误日志
	shouldLog := false

	switch typed := err.(type) {
	case *Errno:
		shouldLog = typed.IsLog
	case *Error:
		if e, ok := typed.Err.(*Errno); ok {
			shouldLog = e.IsLog
		} else {
			// 默认记录所有非 Errno 类型的错误
			shouldLog = true
		}
	default:
		// 默认记录所有未知类型的错误
		shouldLog = true
	}

	if shouldLog {
		// 模拟记录日志，实际项目中应该调用日志包
		fmt.Printf("[ERROR] %v\n", err)
	}
}

// JsonError 返回一个JSON格式的错误响应
func JsonError(c *gin.Context, err error) {
	code, message := Decode(err)

	// 检查是否需要记录错误日志
	logError(c, err)

	// 设置响应状态码
	httpCode := http.StatusOK
	if e, ok := err.(*Errno); ok {
		httpCode = e.ToHTTPStatusCode()
	}

	c.JSON(httpCode, &Response{
		Code:    code,
		Message: message,
	})
}

// NewError 快速创建一个自定义错误
func NewError(code int, msg string) *Errno {
	return &Errno{
		Code:    code,
		Message: msg,
	}
}

// NewErrorf 使用格式化字符串创建一个自定义错误
func NewErrorf(code int, format string, args ...interface{}) *Errno {
	return &Errno{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
