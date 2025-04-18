package log

import (
	"fmt"
	"os"
	"sync"
)

// Options 包含配置日志的选项
type Options struct {
	Level       string   // 日志级别
	Format      string   // 日志格式 (json 或 console)
	OutputPaths []string // 日志输出路径
}

var (
	mu     sync.Mutex
	logger *Logger
)

// Logger 是一个简单的日志记录器，模拟 zap 日志库的行为
type Logger struct {
	level   string
	outputs []string
}

// Init 初始化日志系统
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()

	// 如果没有提供日志级别，默认为 info
	level := "info"
	if opts != nil && opts.Level != "" {
		level = opts.Level
	}

	// 如果没有指定输出路径，默认输出到标准输出
	outputs := []string{"stdout"}
	if opts != nil && len(opts.OutputPaths) > 0 {
		outputs = opts.OutputPaths
	}

	logger = &Logger{
		level:   level,
		outputs: outputs,
	}

	fmt.Println("Log system initialized with level:", level)
}

// Sync 将缓存中的日志刷新到持久化存储中
func Sync() error {
	// 在真实实现中，这里会调用底层日志库的 Sync 方法
	// 在这个简化版本中，我们只是模拟这个行为
	fmt.Println("Log sync called")
	return nil
}

// C 返回带有请求上下文的日志记录器，用于追踪请求
func C(ctx interface{}) *Logger {
	if logger == nil {
		// 默认初始化
		Init(nil)
	}

	// 在实际项目中，这里会将请求ID等信息添加到日志记录器中
	return logger
}

// Infow 记录 info 级别的结构化日志
func Infow(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		// 默认初始化
		Init(nil)
	}

	logMessage(os.Stdout, "INFO", msg, keysAndValues...)
}

// Errorw 记录 error 级别的结构化日志
func Errorw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		// 默认初始化
		Init(nil)
	}

	logMessage(os.Stderr, "ERROR", msg, keysAndValues...)
}

// Fatalw 记录 fatal 级别的结构化日志，并终止程序
func Fatalw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		// 默认初始化
		Init(nil)
	}

	logMessage(os.Stderr, "FATAL", msg, keysAndValues...)
	os.Exit(1)
}

// logMessage 格式化并输出日志消息
func logMessage(output *os.File, level, msg string, keysAndValues ...interface{}) {
	fmt.Fprintf(output, "[%s] %s", level, msg)

	// 处理键值对参数
	if len(keysAndValues) > 0 {
		// 打印键值对
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				fmt.Fprintf(output, " %v=%v", keysAndValues[i], keysAndValues[i+1])
			} else {
				fmt.Fprintf(output, " %v=MISSING_VALUE", keysAndValues[i])
			}
		}
	}

	fmt.Fprintln(output)
}
