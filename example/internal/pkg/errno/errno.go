package errno

import "fmt"

// Errno 定义错误类型
type Errno struct {
	Code    int    // 错误码
	Message string // 错误消息
	Err     error  // 内部错误
}

// Error 实现 error 接口
func (e *Errno) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Error: %s, code: %d, message: %s, error: %s", e.Message, e.Code, e.Message, e.Err)
	}

	return fmt.Sprintf("Error: %s, code: %d, message: %s", e.Message, e.Code, e.Message)
}

// SetMessage 设置错误消息
func (e *Errno) SetMessage(format string, args ...interface{}) *Errno {
	err := *e
	err.Message = fmt.Sprintf(format, args...)
	return &err
}

// Decode 分解错误，获取错误码和错误消息
func Decode(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	// 未知错误统一返回服务器内部错误
	return InternalServerError.Code, InternalServerError.Message
}
