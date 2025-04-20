package miniblog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"example/internal/pkg/auth"
	"example/internal/pkg/log"
	"example/internal/pkg/middleware"
	"example/internal/pkg/schema"
	"example/internal/router"
	"example/internal/routers"
)

var cfgFile string

const (
	// recommendedHomeDir 定义放置配置的默认目录
	recommendedHomeDir = ".miniexample"

	// defaultConfigName 指定默认配置文件名
	defaultConfigName = "config.yaml"
)

// HTTPServer 定义了一个HTTP服务器实例
type HTTPServer struct {
	db     *gorm.DB
	authz  *auth.Authz
	router *router.Router
	server *http.Server
}

// NewHTTPServer 创建一个HTTP服务器实例
func NewHTTPServer() (*HTTPServer, error) {
	// 初始化数据库连接
	dsn := viper.GetString("db.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 初始化数据库表结构
	if err := schema.InitTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database tables: %v", err)
	}

	// 初始化授权组件
	authz, err := auth.NewAuthz(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create authz: %v", err)
	}

	// 初始化JWT配置
	auth.InitToken(auth.TokenConfig{
		SigningKey: viper.GetString("jwt.key"),
		ExpireTime: viper.GetInt("jwt.expire"),
	})

	// 创建路由器
	r := router.New(db, authz)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    viper.GetString("server.addr"),
		Handler: gin.New(),
	}

	return &HTTPServer{
		db:     db,
		authz:  authz,
		router: r,
		server: server,
	}, nil
}

// Run 运行HTTP服务器
func (s *HTTPServer) Run() error {
	// 设置Gin模式
	gin.SetMode(viper.GetString("server.mode"))

	// 创建Gin引擎
	g := gin.New()

	// 加载路由
	s.router.Load(g)

	// 设置HTTP服务器处理程序
	s.server.Handler = g

	// 监听HTTP请求
	log.Infow("Starting HTTP server", "addr", s.server.Addr)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw("Failed to start HTTP server", "err", err)
		}
	}()

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("Shutting down server...")

	// 创建一个5秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorw("Server forced to shutdown", "err", err)
		return err
	}

	log.Infow("Server exiting")
	return nil
}

// NewMiniBlogCommand 创建一个新的cobra命令
func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miniblog",
		Short: "一个博客系统",
		Long:  `一个基于Go语言的博客系统，支持注册、登录、发布文章等功能`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 创建HTTP服务器
			server, err := NewHTTPServer()
			if err != nil {
				return err
			}
			// 运行HTTP服务器
			return server.Run()
		},
	}

	// 绑定命令行标志
	cmd.Flags().StringP("config", "c", "config.yaml", "配置文件路径")

	// 绑定Viper
	viper.BindPFlag("config", cmd.Flags().Lookup("config"))
	cobra.OnInitialize(initConfig)

	return cmd
}

// initConfig 读取配置文件
func initConfig() {
	configFile := viper.GetString("config")
	viper.SetConfigFile(configFile)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log.Init(logOptions())
}

// logOptions 构建日志配置
func logOptions() *log.Options {
	return &log.Options{
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		RotateConfig: &log.RotateConfig{
			MaxSize:    viper.GetInt("log.rotate.max-size"),
			MaxAge:     viper.GetInt("log.rotate.max-age"),
			MaxBackups: viper.GetInt("log.rotate.max-backups"),
			Compress:   viper.GetBool("log.rotate.compress"),
		},
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
		middleware.Recovery(),                          // 1. 恢复中间件，捕获所有 panic
		middleware.RequestID(),                         // 2. 请求 ID 中间件
		middleware.Logger(),                            // 3. 日志中间件
		middleware.TimeoutMiddleware(30 * time.Second), // 4. 超时中间件
		middleware.RateLimiter(100, 200),               // 5. 限流中间件，每秒 100 个请求，突发最大 200
		middleware.NoCache(),                           // 6. 禁用缓存中间件
		middleware.Cors(),                              // 7. CORS 中间件
		middleware.Secure(),                            // 8. 安全中间件
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
