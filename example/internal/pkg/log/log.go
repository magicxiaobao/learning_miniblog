package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Options 包含配置日志的选项
type Options struct {
	Level             string   // 日志级别
	Format            string   // 日志格式 (json 或 console)
	OutputPaths       []string // 日志输出路径
	DisableCaller     bool     // 是否禁用调用者信息
	DisableStacktrace bool     // 是否禁用堆栈跟踪
}

// NewOptions 创建一个带有默认值的 Options 对象
func NewOptions() *Options {
	return &Options{
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             "info",
		Format:            "console",
		OutputPaths:       []string{"stdout"},
	}
}

var (
	mu     sync.Mutex
	logger *Logger
)

// Logger 是一个简单的日志记录器，模拟 zap 日志库的行为
type Logger struct {
	level             string
	format            string
	outputs           []string
	disableCaller     bool
	disableStacktrace bool
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

	format := "console"
	if opts != nil && opts.Format != "" {
		format = opts.Format
	}

	disableCaller := false
	if opts != nil {
		disableCaller = opts.DisableCaller
	}

	disableStacktrace := false
	if opts != nil {
		disableStacktrace = opts.DisableStacktrace
	}

	logger = &Logger{
		level:             level,
		format:            format,
		outputs:           outputs,
		disableCaller:     disableCaller,
		disableStacktrace: disableStacktrace,
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

// Infow 是 Logger 实例的方法，用于记录 info 级别的结构化日志
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	logMessage(os.Stdout, "INFO", msg, keysAndValues...)
}

// Errorw 是 Logger 实例的方法，用于记录 error 级别的结构化日志
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	logMessage(os.Stderr, "ERROR", msg, keysAndValues...)
}

// Debugw 是 Logger 实例的方法，用于记录 debug 级别的结构化日志
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	// 只有在 debug 级别时才输出
	if l.level == "debug" {
		logMessage(os.Stdout, "DEBUG", msg, keysAndValues...)
	}
}

// Warnw 是 Logger 实例的方法，用于记录 warn 级别的结构化日志
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	logMessage(os.Stdout, "WARN", msg, keysAndValues...)
}

// Fatalw 是 Logger 实例的方法，用于记录 fatal 级别的结构化日志
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	logMessage(os.Stderr, "FATAL", msg, keysAndValues...)
	os.Exit(1)
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

// Debugw 记录 debug 级别的结构化日志
func Debugw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		// 默认初始化
		Init(nil)
	}

	// 只有在 debug 级别时才输出
	if logger.level == "debug" {
		logMessage(os.Stdout, "DEBUG", msg, keysAndValues...)
	}
}

// Warnw 记录 warn 级别的结构化日志
func Warnw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		// 默认初始化
		Init(nil)
	}

	logMessage(os.Stdout, "WARN", msg, keysAndValues...)
}

// logMessage 格式化并输出日志消息
func logMessage(output *os.File, level, msg string, keysAndValues ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 获取调用者信息
	var caller string
	if !logger.disableCaller {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			// 获取文件名的短路径
			parts := strings.Split(file, "/")
			if len(parts) > 2 {
				file = strings.Join(parts[len(parts)-2:], "/")
			}
			caller = fmt.Sprintf(" %s:%d", file, line)
		}
	}

	// 基础日志格式
	fmt.Fprintf(output, "[%s] [%s]%s %s", timestamp, level, caller, msg)

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

	// 如果是错误或致命错误且启用了堆栈跟踪，则打印堆栈
	if (level == "ERROR" || level == "FATAL") && !logger.disableStacktrace {
		// 简单的堆栈跟踪实现
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		fmt.Fprintf(output, "Stack trace:\n%s\n", buf[:n])
	}
}
