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
	Level             string        // 日志级别
	Format            string        // 日志格式 (json 或 console)
	OutputPaths       []string      // 日志输出路径
	DisableCaller     bool          // 是否禁用调用者信息
	DisableStacktrace bool          // 是否禁用堆栈跟踪
	ErrorOutputPaths  []string      // 错误输出路径
	RotateConfig      *RotateConfig // 日志轮转配置
}

// NewOptions 创建一个带有默认值的 Options 对象
func NewOptions() *Options {
	return &Options{
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             "info",
		Format:            "console",
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		RotateConfig:      DefaultRotateConfig(),
	}
}

// StdLogger represents the global logger.
var StdLogger = &Logger{
	level:             "info",
	format:            "console",
	outputs:           []string{"stdout"},
	disableCaller:     false,
	disableStacktrace: false,
}

var mu sync.Mutex

// Logger 是一个简单的日志记录器
type Logger struct {
	level             string
	format            string
	outputs           []string
	disableCaller     bool
	disableStacktrace bool
}

// Init initializes the logger with the given options.
func Init(opts *Options) error {
	mu.Lock()
	defer mu.Unlock()

	// Set default options if not specified
	if opts == nil {
		opts = NewOptions()
	}

	// Initialize the standard logger
	StdLogger.level = strings.ToLower(opts.Level)
	StdLogger.format = strings.ToLower(opts.Format)
	StdLogger.disableCaller = opts.DisableCaller
	StdLogger.disableStacktrace = opts.DisableStacktrace

	// Set output paths
	if len(opts.OutputPaths) > 0 {
		StdLogger.outputs = opts.OutputPaths
	}

	// 配置日志轮转
	if opts.RotateConfig != nil {
		GetRotateManager().Configure(opts.RotateConfig)
	}

	return nil
}

// Sync flushes any buffered log entries.
func Sync() error {
	// Since we're using a simple logger without buffering, this is a no-op
	return nil
}

// C 返回带有请求上下文的日志记录器，用于追踪请求
func C(ctx interface{}) *Logger {
	if StdLogger == nil {
		// 默认初始化
		Init(nil)
	}

	// 在实际项目中，这里会将请求ID等信息添加到日志记录器中
	return StdLogger
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
	if StdLogger == nil {
		// 默认初始化
		Init(nil)
	}

	logMessage(os.Stdout, "INFO", msg, keysAndValues...)
}

// Errorw 记录 error 级别的结构化日志
func Errorw(msg string, keysAndValues ...interface{}) {
	if StdLogger == nil {
		// 默认初始化
		Init(nil)
	}

	logMessage(os.Stderr, "ERROR", msg, keysAndValues...)
}

// Fatalw 记录 fatal 级别的结构化日志，并终止程序
func Fatalw(msg string, keysAndValues ...interface{}) {
	if StdLogger == nil {
		// 默认初始化
		Init(nil)
	}

	logMessage(os.Stderr, "FATAL", msg, keysAndValues...)
	os.Exit(1)
}

// Debugw 记录 debug 级别的结构化日志
func Debugw(msg string, keysAndValues ...interface{}) {
	if StdLogger == nil {
		// 默认初始化
		Init(nil)
	}

	// 只有在 debug 级别时才输出
	if StdLogger.level == "debug" {
		logMessage(os.Stdout, "DEBUG", msg, keysAndValues...)
	}
}

// Warnw 记录 warn 级别的结构化日志
func Warnw(msg string, keysAndValues ...interface{}) {
	if StdLogger == nil {
		// 默认初始化
		Init(nil)
	}

	logMessage(os.Stdout, "WARN", msg, keysAndValues...)
}

// logMessage logs a message with the given level and key-value pairs
func logMessage(output *os.File, level, msg string, keysAndValues ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 获取调用者信息
	var caller string
	if !StdLogger.disableCaller {
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
	logEntry := fmt.Sprintf("[%s] [%s]%s %s", timestamp, level, caller, msg)

	// 处理键值对参数
	if len(keysAndValues) > 0 {
		// 打印键值对
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				logEntry += fmt.Sprintf(" %v=%v", keysAndValues[i], keysAndValues[i+1])
			} else {
				logEntry += fmt.Sprintf(" %v=MISSING_VALUE", keysAndValues[i])
			}
		}
	}

	logEntry += "\n"

	// 写入标准输出/错误
	fmt.Fprint(output, logEntry)

	// 创建堆栈跟踪(如果需要)
	var stackTrace string
	if (level == "ERROR" || level == "FATAL") && !StdLogger.disableStacktrace {
		// 简单的堆栈跟踪实现
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		stackTrace = fmt.Sprintf("Stack trace:\n%s\n", buf[:n])

		// 在控制台打印堆栈
		fmt.Fprint(output, stackTrace)
	}

	// 写入文件
	for _, path := range StdLogger.outputs {
		if path != "stdout" && path != "stderr" {
			// 确保日志目录存在
			dir := path[:strings.LastIndex(path, "/")]
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating log directory: %v\n", err)
				continue
			}

			// 准备写入的完整日志内容
			fileContent := logEntry
			if stackTrace != "" {
				fileContent += stackTrace
			}

			// 使用轮转管理器写入日志
			if err := GetRotateManager().Write(path, []byte(fileContent)); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
			}
		}
	}
}
