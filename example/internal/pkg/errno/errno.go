package errno

import (
	"fmt"
	"net/http"
)

// Errno 定义了错误码和消息
type Errno struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	HTTP    int    `json:"-"` // HTTP状态码
	Err     error  `json:"-"` // 内部错误信息
	IsLog   bool   `json:"-"` // 是否需要记录日志
}

// 实现error接口
func (e *Errno) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Err - code: %d, message: %s, error: %s", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("Err - code: %d, message: %s", e.Code, e.Message)
}

// WithErr 添加内部错误
func (e *Errno) WithErr(err error) *Errno {
	newErr := *e
	newErr.Err = err
	return &newErr
}

// WithMsg 覆盖默认消息
func (e *Errno) WithMsg(msg string) *Errno {
	newErr := *e
	newErr.Message = msg
	return &newErr
}

// 预定义错误
var (
	// 系统级错误, Code 范围为 [10000, 19999]
	OK = &Errno{Code: 0, Message: "成功", HTTP: http.StatusOK}

	// 通用错误, Code 范围为 [10000, 10099]
	ErrParam            = &Errno{Code: 10001, Message: "参数错误", HTTP: http.StatusBadRequest}
	ErrValidation       = &Errno{Code: 10002, Message: "验证失败", HTTP: http.StatusBadRequest}
	ErrBind             = &Errno{Code: 10003, Message: "请求体绑定失败", HTTP: http.StatusBadRequest}
	ErrTokenInvalid     = &Errno{Code: 10004, Message: "无效的认证令牌", HTTP: http.StatusUnauthorized}
	ErrTokenExpired     = &Errno{Code: 10005, Message: "认证令牌已过期", HTTP: http.StatusUnauthorized}
	ErrUnauthorized     = &Errno{Code: 10006, Message: "未授权访问", HTTP: http.StatusUnauthorized}
	ErrForbidden        = &Errno{Code: 10007, Message: "禁止访问", HTTP: http.StatusForbidden}
	ErrNotFound         = &Errno{Code: 10008, Message: "资源不存在", HTTP: http.StatusNotFound}
	ErrMethodNotAllowed = &Errno{Code: 10009, Message: "方法不允许", HTTP: http.StatusMethodNotAllowed}
	ErrTooManyRequests  = &Errno{Code: 10010, Message: "请求过于频繁", HTTP: http.StatusTooManyRequests}
	ErrInternalServer   = &Errno{Code: 10011, Message: "服务器内部错误", HTTP: http.StatusInternalServerError}

	// 数据库错误, Code 范围为 [11000, 11099]
	ErrDatabase       = &Errno{Code: 11000, Message: "数据库错误", HTTP: http.StatusInternalServerError}
	ErrRecordExists   = &Errno{Code: 11001, Message: "记录已存在", HTTP: http.StatusConflict}
	ErrRecordNotFound = &Errno{Code: 11002, Message: "记录不存在", HTTP: http.StatusNotFound}

	// 用户错误, Code 范围为 [20100, 20199]
	ErrUserNotFound      = &Errno{Code: 20100, Message: "用户不存在", HTTP: http.StatusNotFound}
	ErrPasswordIncorrect = &Errno{Code: 20101, Message: "密码错误", HTTP: http.StatusBadRequest}
	ErrUserExists        = &Errno{Code: 20102, Message: "用户已存在", HTTP: http.StatusConflict}

	// InternalServerError 用于默认的内部服务器错误
	InternalServerError = &Errno{Code: 10011, Message: "服务器内部错误", HTTP: http.StatusInternalServerError}
)

// SetMessage 设置错误消息
func (e *Errno) SetMessage(format string, args ...interface{}) *Errno {
	err := *e
	err.Message = fmt.Sprintf(format, args...)
	return &err
}

// WithHTTPCode 为错误设置HTTP状态码
func (e *Errno) WithHTTPCode(httpCode int) *Errno {
	newErr := *e
	newErr.HTTP = httpCode
	return &newErr
}

// WithIsLog 设置是否需要记录日志
func (e *Errno) WithIsLog(isLog bool) *Errno {
	newErr := *e
	newErr.IsLog = isLog
	return &newErr
}

// ToHTTPStatusCode 根据 HTTP 状态码选择默认状态码
func (e *Errno) ToHTTPStatusCode() int {
	if e.HTTP != 0 {
		return e.HTTP
	}

	// 基于错误码范围选择 HTTP 状态码
	switch {
	case e.Code >= 10000 && e.Code <= 19999:
		return http.StatusInternalServerError // 5xx 服务器错误
	case e.Code >= 20000 && e.Code <= 29999:
		return http.StatusBadRequest // 4xx 客户端错误
	}

	// 默认返回 OK
	return http.StatusOK
}

// Error 定义通用错误结构
type Error struct {
	Err error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// Decode 分解错误，获取错误码和错误消息
func Decode(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Message
	case *Error:
		if innerErrno, ok := typed.Err.(*Errno); ok {
			return innerErrno.Code, innerErrno.Message
		}
		return InternalServerError.Code, typed.Error()
	case *StackError:
		if innerErrno, ok := typed.Err.(*Errno); ok {
			return innerErrno.Code, innerErrno.Message
		}
		return InternalServerError.Code, typed.Error()
	default:
	}

	// 未知错误统一返回服务器内部错误
	return InternalServerError.Code, InternalServerError.Message
}
