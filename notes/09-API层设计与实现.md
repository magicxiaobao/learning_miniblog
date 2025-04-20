# API层设计与实现

API层是[[Web应用架构]]中连接前端和后端的重要桥梁，它定义了客户端如何与服务器通信的接口规范。在Miniblog项目中，我们采用了RESTful API设计风格，以实现清晰、可扩展的接口。

## RESTful API设计原则

[[RESTful API]]是一种基于HTTP的架构风格，具有以下特点：

1. **资源导向**：使用URI（统一资源标识符）表示资源
2. **HTTP方法语义**：使用HTTP方法表示对资源的操作
   - GET：获取资源
   - POST：创建资源
   - PUT：更新资源
   - DELETE：删除资源
3. **无状态**：服务器不保存客户端状态
4. **统一接口**：标准化的请求和响应格式

### RESTful API命名示例

```
# 用户资源
GET    /v1/users       - 获取用户列表
POST   /v1/users       - 创建用户
GET    /v1/users/:name - 获取指定用户详情
PUT    /v1/users/:name - 更新指定用户
DELETE /v1/users/:name - 删除指定用户

# 博客文章资源
GET    /v1/posts       - 获取文章列表
POST   /v1/posts       - 创建文章
GET    /v1/posts/:id   - 获取指定文章详情
PUT    /v1/posts/:id   - 更新指定文章
DELETE /v1/posts/:id   - 删除指定文章

# 关联资源
GET    /v1/users/:name/posts - 获取指定用户的所有博客
```

## Gin框架基础

[[Gin]]是Go语言中流行的Web框架，它具有高性能、轻量级的特点，非常适合构建API服务。

### 路由定义

Gin使用路由组和中间件来组织API路由：

```go
func (r *Router) Load(g *gin.Engine) {
    // API 版本 v1
    v1 := g.Group("/v1")
    {
        // 无需认证的接口
        v1.POST("/users", r.userCtrl.Create)
        v1.POST("/login", r.userCtrl.Login)

        // 需要认证的接口
        auth := v1.Group("")
        auth.Use(middleware.Authn())
        {
            // 用户相关路由
            users := auth.Group("/users")
            {
                users.GET("/:username", r.userCtrl.Get)
            }

            // 博客文章相关路由
            posts := auth.Group("/posts")
            {
                posts.POST("", r.postCtrl.Create)
            }
        }
    }
}
```

### 控制器实现

控制器负责处理HTTP请求，并返回相应的响应：

```go
// 创建用户控制器
func (ctrl *Controller) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        core.WriteResponse(c, errno.ErrBind, nil)
        return
    }
    
    // 业务逻辑处理...
    
    core.WriteResponse(c, nil, nil)
}
```

## 请求参数绑定与验证

Gin提供了便捷的请求参数绑定和验证机制：

### JSON请求绑定

```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Password string `json:"password" binding:"required,min=6,max=30"`
    Nickname string `json:"nickname" binding:"required,min=1,max=30"`
    Email    string `json:"email" binding:"required,email"`
    Phone    string `json:"phone" binding:"omitempty"`
}

func Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 处理绑定错误
        return
    }
    // ...处理请求
}
```

### URL参数获取

```go
func Get(c *gin.Context) {
    username := c.Param("username")
    if username == "" {
        // 处理参数缺失
        return
    }
    // ...处理请求
}
```

### 查询参数获取

```go
func List(c *gin.Context) {
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    // ...处理请求
}
```

## 统一响应格式

为确保API的一致性，我们定义了统一的响应格式：

```go
// Response 定义统一的API响应格式
type Response struct {
    Code    int         `json:"code"`    // 错误码，0表示成功
    Message string      `json:"message"` // 错误信息
    Data    interface{} `json:"data"`    // 响应数据
}

// WriteResponse 写入响应
func WriteResponse(c *gin.Context, err error, data interface{}) {
    if err != nil {
        // 处理错误响应
        c.JSON(getStatusCode(err), Response{
            Code:    getErrCode(err),
            Message: err.Error(),
            Data:    nil,
        })
        return
    }

    // 成功响应
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "",
        Data:    data,
    })
}
```

## 错误码设计

良好的错误码设计能够帮助客户端理解API返回的错误，便于调试和处理：

```go
const (
    // OK 成功
    OK = 0

    // 通用错误
    ErrUnknown           = 10001 // 未知错误
    ErrBind              = 10002 // 参数绑定错误
    ErrValidation        = 10003 // 参数验证错误
    ErrInternalServer    = 10004 // 服务器内部错误

    // 认证授权相关错误
    ErrTokenInvalid      = 20001 // 无效的token
    ErrUnauthorized      = 20002 // 未授权
    ErrForbidden         = 20003 // 禁止访问

    // 用户相关错误
    ErrUserNotFound      = 30001 // 用户不存在
    ErrUserExists        = 30002 // 用户已存在
    ErrPasswordIncorrect = 30003 // 密码错误
)
```

## API版本控制

API版本控制是确保向后兼容性的重要策略：

1. **URL路径版本控制**：如 `/v1/users`、`/v2/users`
2. **HTTP头版本控制**：如 `Accept: application/vnd.miniblog.v1+json`
3. **查询参数版本控制**：如 `/users?version=1`

在Miniblog项目中，我们采用URL路径版本控制，这是最直观和简单的方式。

## 中间件实现

中间件是Gin框架提供的强大功能，用于处理请求和响应的拦截器：

```go
// 日志中间件
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        // 处理请求
        c.Next()
        
        // 记录请求日志
        latency := time.Since(start)
        statusCode := c.Writer.Status()
        log.Infow("API request",
            "method", c.Request.Method,
            "path", path,
            "status", statusCode,
            "latency", latency,
            "client_ip", c.ClientIP(),
        )
    }
}
```

## 常用中间件

在Miniblog项目中，我们实现了以下中间件：

1. **Logger**：记录请求日志
2. **Recovery**：恢复panic，避免服务器崩溃
3. **RequestID**：为每个请求生成唯一ID
4. **CORS**：处理跨域请求
5. **Secure**：设置安全相关的HTTP头
6. **NoCache**：禁止客户端缓存
7. **Authn**：用户认证
8. **Authz**：用户授权

## 文档生成

API文档对于前后端协作和API维护至关重要。常用的API文档工具包括：

1. **Swagger/OpenAPI**：最流行的API文档标准
2. **API Blueprint**：简洁的API文档格式
3. **RAML**：RESTful API建模语言

在Go生态中，常用swagger-ui和swaggo生成Swagger文档：

```go
// @title Miniblog API
// @version 1.0
// @description 这是一个简单的博客API服务

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
func main() {
    // ...
}
```

## API测试

为确保API质量，我们需要进行多层次的测试：

1. **单元测试**：测试控制器逻辑
2. **集成测试**：测试API端到端功能
3. **性能测试**：测试API性能和并发处理能力

可以使用HTTP客户端工具（如Postman、IntelliJ HTTP Client）创建API测试集合：

```http
### 用户登录
POST http://localhost:8080/v1/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "testpassword123"
}

> {%
    if (response.body.data && response.body.data.token) {
        client.global.set("token", response.body.data.token);
    }
%}

### 创建博客文章
POST http://localhost:8080/v1/posts
Content-Type: application/json
Authorization: Bearer {{token}}

{
  "title": "测试文章",
  "content": "这是一篇测试文章"
}
```

## 最佳实践

1. **遵循RESTful设计原则**：使用正确的HTTP方法和状态码
2. **资源命名使用复数**：如`/users`而非`/user`
3. **版本化API**：在URL路径中包含版本号
4. **正确处理错误**：返回有意义的错误码和错误信息
5. **使用JSON作为数据交换格式**：简单、跨平台、易于阅读
6. **分页处理大量数据**：使用limit/offset或cursor分页
7. **实现API限流**：保护API免受滥用
8. **提供完整文档**：帮助API使用者理解和使用API

## 相关链接

- [[Gin框架文档|https://gin-gonic.com/docs/]]
- [[RESTful API设计最佳实践]]
- [[swagger文档生成工具|https://github.com/swaggo/swag]]
- [[前一篇：模型层与数据库操作]] 