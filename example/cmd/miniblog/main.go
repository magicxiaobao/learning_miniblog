package main

import (
	"fmt"
	"os"

	"example/internal/miniblog"
)

// main 函数是整个程序的入口点
func main() {
	// 创建命令行应用实例
	command := miniblog.NewMiniBlogCommand()

	// 执行命令行程序
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
