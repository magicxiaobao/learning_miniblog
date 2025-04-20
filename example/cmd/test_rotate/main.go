package main

import (
	"fmt"
	"os"
	"time"

	"example/internal/pkg/log"
)

func main() {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
		os.Exit(1)
	}

	// 配置日志
	logOpts := &log.Options{
		Level:             "debug",
		Format:            "console",
		OutputPaths:       []string{"stdout", "./logs/rotate-test.log"},
		DisableCaller:     false,
		DisableStacktrace: false,
		RotateConfig: &log.RotateConfig{
			MaxSize:    1, // 设置为1MB以便快速触发轮转
			MaxAge:     7,
			MaxBackups: 5,
			Compress:   false,
		},
	}

	if err := log.Init(logOpts); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("开始测试日志轮转功能...")
	fmt.Println("将生成约5MB的日志数据，触发多次日志轮转")

	// 生成足够大的日志以触发多次轮转
	for i := 0; i < 50000; i++ {
		log.Infow(fmt.Sprintf("测试日志消息 #%d", i),
			"timestamp", time.Now().UnixNano(),
			"data", fmt.Sprintf("这是一些测试数据，用于填充日志文件 %d", i),
			"iteration", i)

		// 每1000条日志暂停一下，让用户能看到进度
		if i%1000 == 0 && i > 0 {
			fmt.Printf("已生成 %d 条日志记录\n", i)
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("日志生成完成，请检查logs目录中的日志文件")

	// 列出日志文件
	files, err := os.ReadDir("logs")
	if err != nil {
		fmt.Printf("读取日志目录失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n日志文件列表:")
	for _, file := range files {
		info, _ := file.Info()
		fmt.Printf("- %s (%.2f KB)\n", file.Name(), float64(info.Size())/1024)
	}
}
