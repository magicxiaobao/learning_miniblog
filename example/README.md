# miniblog 学习示例项目

这是一个用于学习 [marmotedu/miniblog](https://github.com/marmotedu/miniblog) 项目的示例代码仓库。该项目模拟实现了 miniblog 的核心功能和架构设计，便于理解原项目的设计理念和实现方式。

## 项目结构

```
.
├── cmd/                  # 应用程序入口
│   └── miniblog/         # miniblog 应用入口
├── internal/             # 私有代码
│   ├── miniblog/         # 核心应用代码
│   │   ├── biz/          # 业务逻辑层
│   │   ├── controller/   # 控制器层
│   │   └── store/        # 存储层
│   └── pkg/              # 内部共享包
│       ├── log/          # 日志包
│       └── model/        # 数据模型
├── logs/                 # 日志目录
├── config.yaml           # 配置文件
├── learning_outline.md   # 学习大纲
└── README.md             # 项目说明
```

## 学习内容

本示例项目按照 `learning_outline.md` 中的大纲逐步实现，包括：

1. 项目基础与架构设计
2. 项目启动流程（命令行解析、配置管理、服务启动）
3. 日志系统
4. 错误处理机制
5. 中间件设计与实现
6. 数据库访问层
7. 用户管理模块
8. 博客管理模块
9. API 设计与文档
10. 测试策略
11. 部署与发布

## 运行项目

```bash
# 运行项目
$ go run main.go

# 使用指定配置文件运行
$ go run main.go -c /path/to/config.yaml
```

## 测试 API

启动服务后，可以通过以下方式测试 API：

```bash
# 健康检查
$ curl http://localhost:8080/healthz
``` 