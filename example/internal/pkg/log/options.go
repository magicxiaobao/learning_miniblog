package log

// Options 包含配置日志的选项
type Options struct {
	// 是否开启 caller，如果开启会在日志中显示调用日志所在的文件和行号
	DisableCaller bool
	// 是否禁止在 panic 及以上级别打印堆栈信息
	DisableStacktrace bool
	// 日志级别：debug, info, warn, error, dpanic, panic, fatal
	Level string
	// 日志格式：console, json
	Format string
	// 日志输出位置
	OutputPaths []string
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
