package errno

// 错误码规则:
// - 错误码为 5 位数字
// - 第 1 位为错误产生来源: 1 为服务端, 2 为客户端
// - 第 2-3 位为模块代码
// - 最后 2 位为具体错误代码

// OK 代表请求成功，无错误发生
var OK = &Errno{Code: 0, Message: "OK"}

// 通用错误, 前 3 位为 100
var (
	// InternalServerError 表示所有未知的服务端错误
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}

	// ErrBind 表示参数绑定错误
	ErrBind = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct"}

	// ErrValidation 表示参数验证失败
	ErrValidation = &Errno{Code: 10003, Message: "Validation failed"}

	// ErrTokenInvalid 表示 JWT Token 格式错误
	ErrTokenInvalid = &Errno{Code: 10004, Message: "Token is invalid"}

	// ErrPageNotFound 表示路由不匹配错误
	ErrPageNotFound = &Errno{Code: 10005, Message: "Page not found"}
)

// 用户模块错误, 前 3 位为 101
var (
	// ErrUserAlreadyExist 表示用户已经存在
	ErrUserAlreadyExist = &Errno{Code: 10101, Message: "User already exists"}

	// ErrUserNotFound 表示用户不存在
	ErrUserNotFound = &Errno{Code: 10102, Message: "User was not found"}

	// ErrPasswordIncorrect 表示密码不正确
	ErrPasswordIncorrect = &Errno{Code: 10103, Message: "Password is incorrect"}
)

// 博客模块错误, 前 3 位为 102
var (
	// ErrPostNotFound 表示博客不存在
	ErrPostNotFound = &Errno{Code: 10201, Message: "Post was not found"}

	// ErrPostCreateFailed 表示创建博客失败
	ErrPostCreateFailed = &Errno{Code: 10202, Message: "Post create failed"}
)
