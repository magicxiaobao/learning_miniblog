package errno

import (
	"fmt"
	"runtime"
	"strings"
)

// StackError 表示生成的错误附带有源文件和行号信息
type StackError struct {
	Err    error
	Stack  string
	Caller string
}

// 确保 StackError 实现了 error 接口
var _ error = (*StackError)(nil)

// Error 返回错误的详细文本表示
func (e *StackError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s\n%s\n%s", e.Err.Error(), e.Caller, e.Stack)
	}
	return fmt.Sprintf("Error occurred, but err is nil")
}

// New 根据传入的 err 创建一个带调用栈的 StackError 实例
func New(err error) *StackError {
	if err == nil {
		return nil
	}

	// 获取调用栈
	callers := make([]uintptr, 32)
	n := runtime.Callers(2, callers)
	if n == 0 {
		return &StackError{Err: err}
	}

	// 处理调用栈信息
	frames := runtime.CallersFrames(callers[:n])
	var stackInfo strings.Builder
	var caller string

	// 记录第一个调用位置
	if frame, more := frames.Next(); more {
		caller = fmt.Sprintf("Called from %s:%d", frame.File, frame.Line)
	}

	// 构建完整调用栈
	for i := 0; i < 10; i++ { // 限制堆栈深度
		frame, more := frames.Next()
		if !more {
			break
		}
		// 排除标准库和无关文件
		if strings.Contains(frame.File, "/runtime/") || strings.Contains(frame.File, "/reflect/") {
			continue
		}
		stackInfo.WriteString(fmt.Sprintf("%s:%d - %s\n", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	return &StackError{
		Err:    err,
		Stack:  stackInfo.String(),
		Caller: caller,
	}
}

// Wrap 将自定义错误和标准错误包装在一起，带有堆栈信息
func Wrap(err error, code int, message string) error {
	if err == nil {
		return &Errno{Code: code, Message: message}
	}

	// 创建 Errno 并关联原始错误
	errno := &Errno{
		Code:    code,
		Message: message,
		Err:     err,
	}

	// 包装错误和堆栈信息
	return New(errno)
}

// WrapC 将自定义错误和标准错误包装在一起，带有更多上下文信息和堆栈信息
func WrapC(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return &Errno{Code: code, Message: fmt.Sprintf(format, args...)}
	}

	// 创建 Errno 并关联原始错误
	errno := &Errno{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Err:     err,
	}

	// 包装错误和堆栈信息
	return New(errno)
}

// Cause 从错误链中提取最原始的错误
func Cause(err error) error {
	if err == nil {
		return nil
	}

	switch typed := err.(type) {
	case *StackError:
		return Cause(typed.Err)
	case *Errno:
		if typed.Err != nil {
			return Cause(typed.Err)
		}
		return typed
	default:
		return err
	}
}
