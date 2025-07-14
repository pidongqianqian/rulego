# RuleGo Server Package

基于rulego的类似node-red的server包，提供完整的规则引擎管理服务。

## 特性

- 🚀 **快速开发**：可以快速通过这个server开发各种应用
- 🔧 **模块化设计**：每个模块和扩展点可以灵活扩展和替换
- 🏗️ **分层架构**：控制层、服务层、DAO层接口化，保留默认实现
- 🔄 **无缝集成**：可以快速集成到现有系统
- 📦 **零依赖**：尽量不引入其他第三方依赖
- 🛠️ **工具链支持**：无需写代码即可编译出新的Server应用
- 🔀 **框架无关**：控制层可以无缝切换到现有web框架（如gin）

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    控制层 (Controller)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   HTTP API  │  │   gRPC API  │  │   WebSocket │          │
│  │ Controller  │  │ Controller  │  │ Controller  │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    服务层 (Service)                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │    Rule     │  │  Component  │  │   Plugin    │          │
│  │  Service    │  │  Service    │  │  Service    │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   数据层 (DAO)                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   File DAO  │  │   DB DAO    │  │   Cache DAO │          │
│  │ (default)   │  │ (pluggable) │  │ (pluggable) │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

## 项目结构

```
rulego-server/
├── cmd/                      # 命令行工具
│   ├── server/              # 服务端主程序
│   └── toolkit/             # 工具链
├── pkg/                     # 公共包
│   ├── controller/          # 控制层接口
│   ├── service/             # 服务层接口
│   ├── dao/                 # 数据访问层接口
│   ├── model/               # 数据模型
│   └── config/              # 配置管理
├── internal/                # 内部实现
│   ├── controller/          # 控制层实现
│   ├── service/             # 服务层实现
│   ├── dao/                 # 数据访问层实现
│   └── middleware/          # 中间件
├── configs/                 # 配置文件
├── examples/                # 示例代码
├── docs/                    # 文档
└── tools/                   # 工具脚本
```

## 核心接口设计

### 1. 规则管理接口

```go
// RuleService 规则管理服务接口
type RuleService interface {
    // 创建规则链
    CreateRuleChain(ctx context.Context, req *CreateRuleChainRequest) (*RuleChain, error)
    
    // 获取规则链
    GetRuleChain(ctx context.Context, chainId string) (*RuleChain, error)
    
    // 更新规则链
    UpdateRuleChain(ctx context.Context, chainId string, req *UpdateRuleChainRequest) error
    
    // 删除规则链
    DeleteRuleChain(ctx context.Context, chainId string) error
    
    // 部署规则链
    DeployRuleChain(ctx context.Context, chainId string) error
    
    // 取消部署
    UndeployRuleChain(ctx context.Context, chainId string) error
    
    // 执行规则链
    ExecuteRuleChain(ctx context.Context, chainId string, msg *RuleMessage) (*ExecuteResult, error)
    
    // 获取规则链列表
    ListRuleChains(ctx context.Context, req *ListRuleChainsRequest) (*ListRuleChainsResponse, error)
}
```

### 2. 组件管理接口

```go
// ComponentService 组件管理服务接口
type ComponentService interface {
    // 获取所有组件
    GetComponents(ctx context.Context) ([]Component, error)
    
    // 获取组件详情
    GetComponent(ctx context.Context, componentType string) (*Component, error)
    
    // 注册组件
    RegisterComponent(ctx context.Context, component Component) error
    
    // 注销组件
    UnregisterComponent(ctx context.Context, componentType string) error
}
```

### 3. 数据访问接口

```go
// RuleDAO 规则数据访问接口
type RuleDAO interface {
    // 保存规则链
    SaveRuleChain(ctx context.Context, chain *RuleChain) error
    
    // 获取规则链
    GetRuleChain(ctx context.Context, chainId string) (*RuleChain, error)
    
    // 删除规则链
    DeleteRuleChain(ctx context.Context, chainId string) error
    
    // 查询规则链
    QueryRuleChains(ctx context.Context, query *RuleChainQuery) ([]*RuleChain, error)
}
```

## API接口设计

### 1. 规则链管理API

```http
# 创建规则链
POST /api/v1/rule-chains
Content-Type: application/json

{
  "name": "测试规则链",
  "description": "这是一个测试规则链",
  "definition": {
    "ruleChain": {
      "name": "测试规则链"
    },
    "metadata": {
      "nodes": [...],
      "connections": [...]
    }
  }
}

# 获取规则链
GET /api/v1/rule-chains/{chainId}

# 更新规则链
PUT /api/v1/rule-chains/{chainId}

# 删除规则链
DELETE /api/v1/rule-chains/{chainId}

# 部署规则链
POST /api/v1/rule-chains/{chainId}/deploy

# 取消部署
POST /api/v1/rule-chains/{chainId}/undeploy

# 执行规则链
POST /api/v1/rule-chains/{chainId}/execute
Content-Type: application/json

{
  "msgType": "POST",
  "dataType": "JSON",
  "data": "{\"temperature\": 25}",
  "metadata": {
    "deviceId": "device001"
  }
}

# 获取规则链列表
GET /api/v1/rule-chains?page=1&size=10&keyword=test
```

### 2. 组件管理API

```http
# 获取所有组件
GET /api/v1/components

# 获取组件详情
GET /api/v1/components/{componentType}

# 注册组件
POST /api/v1/components

# 注销组件
DELETE /api/v1/components/{componentType}
```

### 3. 系统管理API

```http
# 获取系统状态
GET /api/v1/system/status

# 获取系统配置
GET /api/v1/system/config

# 更新系统配置
PUT /api/v1/system/config

# 重新加载配置
POST /api/v1/system/reload
```

## 使用示例

### 1. 基本使用

```go
package main

import (
    "context"
    "log"
    
    "github.com/rulego/rulego-server/pkg/server"
    "github.com/rulego/rulego-server/pkg/config"
)

func main() {
    // 加载配置
    cfg := config.Default()
    
    // 创建服务器
    srv := server.New(cfg)
    
    // 启动服务器
    if err := srv.Start(context.Background()); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### 2. 自定义DAO实现

```go
package main

import (
    "context"
    "database/sql"
    
    "github.com/rulego/rulego-server/pkg/dao"
    "github.com/rulego/rulego-server/pkg/model"
    "github.com/rulego/rulego-server/pkg/server"
)

// 自定义数据库DAO实现
type MySQLRuleDAO struct {
    db *sql.DB
}

func (d *MySQLRuleDAO) SaveRuleChain(ctx context.Context, chain *model.RuleChain) error {
    // 实现保存逻辑
    return nil
}

func (d *MySQLRuleDAO) GetRuleChain(ctx context.Context, chainId string) (*model.RuleChain, error) {
    // 实现获取逻辑
    return nil, nil
}

// 更多方法实现...

func main() {
    // 创建数据库连接
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/rulego")
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建自定义DAO
    customDAO := &MySQLRuleDAO{db: db}
    
    // 配置服务器使用自定义DAO
    cfg := config.Default()
    cfg.DAO.Type = "mysql"
    
    srv := server.New(cfg)
    srv.SetRuleDAO(customDAO)
    
    // 启动服务器
    srv.Start(context.Background())
}
```

### 3. 集成到Gin框架

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/rulego/rulego-server/pkg/controller"
    "github.com/rulego/rulego-server/pkg/service"
)

func main() {
    // 创建gin引擎
    r := gin.Default()
    
    // 创建服务
    ruleService := service.NewRuleService()
    
    // 创建控制器
    ruleController := controller.NewRuleController(ruleService)
    
    // 注册路由
    api := r.Group("/api/v1")
    {
        api.POST("/rule-chains", ruleController.CreateRuleChain)
        api.GET("/rule-chains/:chainId", ruleController.GetRuleChain)
        api.PUT("/rule-chains/:chainId", ruleController.UpdateRuleChain)
        api.DELETE("/rule-chains/:chainId", ruleController.DeleteRuleChain)
        api.POST("/rule-chains/:chainId/execute", ruleController.ExecuteRuleChain)
    }
    
    // 启动服务器
    r.Run(":8080")
}
```

## 工具链支持

### 1. 配置文件驱动编译

创建 `server.yaml` 配置文件：

```yaml
# 服务器配置
server:
  name: "my-rulego-server"
  port: 8080
  
# 组件配置
components:
  # 内置组件
  builtin:
    - "filter"
    - "transform"
    - "action"
  
  # 扩展组件
  extensions:
    - "github.com/rulego/rulego-components/external/kafka"
    - "github.com/rulego/rulego-components/external/redis"
    - "github.com/rulego/rulego-components/external/mongodb"
    
# 插件配置
plugins:
  - path: "./plugins/custom-plugin.so"
    name: "custom-plugin"
    
# DAO配置
dao:
  type: "mysql"
  connection: "user:password@tcp(localhost:3306)/rulego"
  
# 认证配置
auth:
  enabled: true
  jwt_secret: "your-secret-key"
  
# 监控配置
monitoring:
  metrics_enabled: true
  prometheus_endpoint: "/metrics"
```

### 2. 使用工具链编译

```bash
# 使用工具链编译服务器
./toolkit build --config server.yaml --output my-server

# 编译后的文件结构
my-server/
├── bin/
│   └── server              # 可执行文件
├── configs/
│   └── server.yaml         # 配置文件
├── plugins/                # 插件目录
└── data/                   # 数据目录
```

### 3. 工具链命令

```bash
# 初始化项目
./toolkit init --name my-project

# 编译服务器
./toolkit build --config server.yaml

# 生成插件模板
./toolkit generate plugin --name my-plugin

# 验证配置
./toolkit validate --config server.yaml

# 升级服务器
./toolkit upgrade --version v1.2.0
```

## 扩展示例

### 1. 自定义组件

```go
package main

import (
    "context"
    "github.com/rulego/rulego/api/types"
    "github.com/rulego/rulego/components/base"
)

// 自定义组件
type CustomComponent struct {
    base.BaseNode
    Config CustomConfig
}

type CustomConfig struct {
    // 组件配置
    Timeout int `json:"timeout"`
}

func (c *CustomComponent) Type() string {
    return "custom/myComponent"
}

func (c *CustomComponent) New() types.Node {
    return &CustomComponent{Config: CustomConfig{Timeout: 30}}
}

func (c *CustomComponent) Init(ruleConfig types.Config, configuration types.Configuration) error {
    // 初始化逻辑
    return nil
}

func (c *CustomComponent) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
    // 处理消息逻辑
    ctx.TellSuccess(msg)
}

// 注册组件
func init() {
    rulego.Registry.Register(&CustomComponent{})
}
```

### 2. 中间件扩展

```go
package middleware

import (
    "context"
    "time"
    
    "github.com/rulego/rulego-server/pkg/controller"
)

// 请求日志中间件
func RequestLogging() controller.Middleware {
    return func(next controller.HandlerFunc) controller.HandlerFunc {
        return func(ctx context.Context, req *controller.Request) (*controller.Response, error) {
            start := time.Now()
            
            resp, err := next(ctx, req)
            
            duration := time.Since(start)
            log.Printf("Request: %s %s - Duration: %v", req.Method, req.Path, duration)
            
            return resp, err
        }
    }
}

// 认证中间件
func Authentication() controller.Middleware {
    return func(next controller.HandlerFunc) controller.HandlerFunc {
        return func(ctx context.Context, req *controller.Request) (*controller.Response, error) {
            // 验证JWT token
            token := req.Header.Get("Authorization")
            if !isValidToken(token) {
                return &controller.Response{
                    Status: 401,
                    Message: "Unauthorized",
                }, nil
            }
            
            return next(ctx, req)
        }
    }
}
```

## 部署方式

### 1. 独立部署

```bash
# 编译
go build -o rulego-server cmd/server/main.go

# 运行
./rulego-server --config configs/server.yaml
```

### 2. Docker部署

```dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o rulego-server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/rulego-server .
COPY configs/ ./configs/

CMD ["./rulego-server", "--config", "configs/server.yaml"]
```

### 3. Kubernetes部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rulego-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: rulego-server
  template:
    metadata:
      labels:
        app: rulego-server
    spec:
      containers:
      - name: rulego-server
        image: rulego-server:latest
        ports:
        - containerPort: 8080
        env:
        - name: CONFIG_FILE
          value: "/etc/config/server.yaml"
        volumeMounts:
        - name: config
          mountPath: /etc/config
      volumes:
      - name: config
        configMap:
          name: rulego-server-config
```

## 性能优化

### 1. 规则链缓存

```go
type CachedRuleService struct {
    cache sync.Map
    delegate RuleService
}

func (s *CachedRuleService) GetRuleChain(ctx context.Context, chainId string) (*RuleChain, error) {
    if cached, ok := s.cache.Load(chainId); ok {
        return cached.(*RuleChain), nil
    }
    
    chain, err := s.delegate.GetRuleChain(ctx, chainId)
    if err == nil {
        s.cache.Store(chainId, chain)
    }
    
    return chain, err
}
```

### 2. 连接池优化

```go
type PooledRuleEngine struct {
    pool sync.Pool
}

func (p *PooledRuleEngine) GetEngine() *rulego.RuleEngine {
    if engine := p.pool.Get(); engine != nil {
        return engine.(*rulego.RuleEngine)
    }
    return rulego.New()
}

func (p *PooledRuleEngine) PutEngine(engine *rulego.RuleEngine) {
    engine.Reset()
    p.pool.Put(engine)
}
```

## 监控与日志

### 1. 性能监控

```go
type MetricsCollector struct {
    requestCount    prometheus.Counter
    requestDuration prometheus.Histogram
}

func (m *MetricsCollector) RecordRequest(duration time.Duration) {
    m.requestCount.Inc()
    m.requestDuration.Observe(duration.Seconds())
}
```

### 2. 结构化日志

```go
type StructuredLogger struct {
    logger *slog.Logger
}

func (l *StructuredLogger) Info(msg string, fields ...slog.Attr) {
    l.logger.Info(msg, fields...)
}

func (l *StructuredLogger) Error(msg string, err error, fields ...slog.Attr) {
    l.logger.Error(msg, slog.Any("error", err), fields...)
}
```

## 配置管理

### 1. 多环境配置

```yaml
# configs/base.yaml
server:
  name: "rulego-server"
  
---
# configs/dev.yaml
server:
  port: 8080
  debug: true
  
---
# configs/prod.yaml
server:
  port: 80
  debug: false
```

### 2. 动态配置

```go
type ConfigWatcher struct {
    config *Config
    watchers []func(*Config))
}

func (w *ConfigWatcher) Watch(callback func(*Config)) {
    w.watchers = append(w.watchers, callback)
}

func (w *ConfigWatcher) Reload() {
    newConfig := loadConfig()
    for _, callback := range w.watchers {
        callback(newConfig)
    }
}
```

## 最佳实践

1. **模块化设计**：将不同功能模块化，便于维护和扩展
2. **接口导向**：通过接口定义契约，提高代码的可测试性
3. **配置外部化**：将配置信息外部化，支持多环境部署
4. **错误处理**：统一错误处理机制，提供友好的错误信息
5. **性能优化**：合理使用缓存、连接池等技术提升性能
6. **监控告警**：建立完善的监控体系，及时发现问题
7. **安全防护**：实施认证授权、数据加密等安全措施

## 总结

这个RuleGo Server包提供了一个完整的、可扩展的规则引擎服务解决方案。它不仅满足了类似node-red的功能需求，还通过模块化设计和工具链支持，使得开发者可以快速构建定制化的规则引擎应用。

通过分层架构设计，各层之间职责清晰，便于维护和扩展。同时，通过接口化设计，使得各个组件可以灵活替换，满足不同场景的需求。

工具链的支持让开发者可以通过配置文件快速编译出符合需求的服务器应用，无需编写大量的样板代码。这大大提高了开发效率，降低了使用门槛。