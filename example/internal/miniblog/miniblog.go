package miniblog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"example/internal/pkg/log"
	"example/internal/pkg/middleware"
	"example/internal/routers"
)

var cfgFile string

const (
	// recommendedHomeDir 定义放置配置的默认目录
	recommendedHomeDir = ".miniexample"

	// defaultConfigName 指定默认配置文件名
	defaultConfigName = "config.yaml"
)

// NewMiniBlogCommand 创建命令行应用程序
func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miniexample",
		Short: "A simplified version of miniblog",
		Long: `A simplified version that demonstrates the key concepts of miniblog.
This is part of the learning project for marmotedu/miniblog.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 初始化日志
			log.Init(logOptions())
			defer log.Sync()

			log.Infow("Starting miniexample application", "version", "v0.1.0")
			return run()
		},
	}

	// 设置命令行参数解析回调
	cobra.OnInitialize(initConfig)

	// 添加命令行参数
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file path")

	return cmd
}

// initConfig 加载配置文件和环境变量
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigName)
	}

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MINIEXAMPLE")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 读取配置文件
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Warning: No config file found. Using defaults.")
		// 设置默认值
		viper.SetDefault("server.addr", ":8080")
		viper.SetDefault("server.timeout", 10)
		viper.SetDefault("runmode", "debug")
	}
}

// logOptions 构建日志配置
func logOptions() *log.Options {
	return &log.Options{
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
	}
}

// run 实现主程序逻辑
func run() error {
	// 设置 Gin 模式
	gin.SetMode(viper.GetString("runmode"))

	// 创建 Gin 引擎
	g := gin.New()

	// 添加中间件
	mws := []gin.HandlerFunc{
		gin.Recovery(),
		middleware.NoCache(),
		middleware.Cors(),
		middleware.Secure(),
		middleware.RequestID(),
		middleware.Logger(),
	}
	g.Use(mws...)

	// 注册路由
	if err := routers.InstallRouters(g); err != nil {
		return err
	}

	// 启动 HTTP 服务器
	srv := startServer(g)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("server.timeout"))*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorw("Server forced to shutdown", "err", err)
		return err
	}

	log.Infow("Server exiting")
	return nil
}

// startServer 启动 HTTP 服务器
func startServer(g *gin.Engine) *http.Server {
	addr := viper.GetString("server.addr")
	srv := &http.Server{
		Addr:    addr,
		Handler: g,
	}

	// 在 goroutine 中启动服务
	log.Infow("Starting HTTP server", "addr", addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw("Server failed", "err", err)
		}
	}()

	return srv
}
