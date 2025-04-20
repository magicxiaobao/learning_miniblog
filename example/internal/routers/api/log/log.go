package log

import (
	"github.com/gin-gonic/gin"

	"example/internal/pkg/errno"
	"example/internal/pkg/log"
)

// DemoDebug 演示 debug 级别日志
func DemoDebug(c *gin.Context) {
	log.C(c).Debugw("这是一条调试日志", "level", "debug", "user", "admin")

	errno.WriteResponse(c, nil, "Debug 日志已记录")
}

// DemoInfo 演示 info 级别日志
func DemoInfo(c *gin.Context) {
	log.C(c).Infow("这是一条信息日志", "level", "info", "user", "admin")

	errno.WriteResponse(c, nil, "Info 日志已记录")
}

// DemoWarn 演示 warn 级别日志
func DemoWarn(c *gin.Context) {
	log.C(c).Warnw("这是一条警告日志", "level", "warn", "user", "admin")

	errno.WriteResponse(c, nil, "Warn 日志已记录")
}

// DemoError 演示 error 级别日志
func DemoError(c *gin.Context) {
	log.C(c).Errorw("这是一条错误日志", "level", "error", "user", "admin")

	// 返回一个错误，触发错误日志记录
	errno.WriteResponse(c, errno.ErrDatabase, nil)
}

// DemoErrorWithStack 演示带有堆栈信息的错误日志
func DemoErrorWithStack(c *gin.Context) {
	// 创建一个带有堆栈信息的错误
	err := errno.Wrap(errno.ErrDatabase, errno.ErrDatabase.Code, "数据库连接失败")

	// 记录错误日志
	log.C(c).Errorw("这是一条带有堆栈信息的错误日志", "err", err)

	// 返回错误
	errno.WriteResponse(c, err, nil)
}
