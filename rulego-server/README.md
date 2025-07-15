# RuleGo Server

基于 [RuleGo](https://github.com/rulego/rulego) 的高性能规则引擎服务器包，提供类似 Node-RED 的可视化流程编排功能。

## 特性

- 🚀 **高性能**: 基于 RuleGo 的轻量级规则引擎
- 🎨 **可视化编辑**: 支持 RuleGo-Editor 可视化流程编排
- 🔧 **零代码编译**: 通过 YAML 配置文件即可构建完整服务器
- 🏗️ **模块化设计**: 清晰的三层架构 (Controller/Service/DAO)
- 🔌 **插件系统**: 支持动态加载组件和插件
- 📊 **监控支持**: 内置 Prometheus 监控和健康检查
- 🔐 **认证授权**: 支持 JWT 认证和多租户隔离
- 🐳 **容器化**: 支持 Docker 和 Kubernetes 部署

## 快速开始

### 1. 安装依赖

```bash
# 克隆项目
git clone https://github.com/pidongqianqian/rulego-server.git
cd rulego-server

# 安装依赖
make deps
```

### 2. 编译项目

```bash
# 编译服务器
make build

# 编译工具链
make build-toolkit

# 编译所有
make build-all
```

### 3. 运行服务器

```bash
# 使用默认配置运行
make run

# 使用自定义配置运行
make run-config

# 或者直接运行
./bin/rulego-server -config configs/config.yaml
```

### 4. 使用工具链

```bash
# 初始化新项目
./bin/toolkit init my-project

# 构建项目
./bin/toolkit build -c my-project/server.yaml -o ./dist

# 生成组件代码
./bin/toolkit generate component my-component -o ./components

# 验证配置
./bin/toolkit validate -c server.yaml
```

## 项目结构

```
rulego-server/
├── cmd/                    # 命令行工具
│   ├── server/            # 服务器主程序
│   └── toolkit/           # 工具链CLI
├── pkg/                   # 公共包
│   ├── server/           # 服务器核心
│   ├── service/          # 服务接口
│   ├── model/            # 数据模型
│   └── config/           # 配置管理
├── internal/             # 内部包
│   └── toolkit/          # 工具链实现
├── configs/              # 配置文件
├── docs/                 # 文档
├── examples/             # 示例
├── Makefile             # 构建脚本
└── go.mod               # 依赖管理
```

## 配置文件

### 基础配置 (server.yaml)

```yaml
# 服务器配置
server:
  name: "my-server"
  port: 8080
  version: "1.0.0"

# 组件配置
components:
  builtin:
    - "filter"
    - "transform"
    - "action"
  extensions: []

# 插件配置
plugins: []

# 数据库配置
dao:
  type: "file"
  connection: "./data"

# 认证配置
auth:
  enabled: false
  jwt_secret: "your-secret-key"

# 监控配置
monitoring:
  metrics_enabled: true
  prometheus_endpoint: "/metrics"
```

## API 接口

### 规则链管理

```bash
# 创建规则链
curl -X POST http://localhost:8080/api/v1/rule-chains \
  -H "Content-Type: application/json" \
  -d '{"name": "test-chain", "definition": {...}}'

# 获取规则链
curl http://localhost:8080/api/v1/rule-chains/test-chain

# 执行规则链
curl -X POST http://localhost:8080/api/v1/rule-chains/test-chain/execute \
  -H "Content-Type: application/json" \
  -d '{"data": "test message"}'
```

### 组件管理

```bash
# 获取组件列表
curl http://localhost:8080/api/v1/components

# 获取组件详情
curl http://localhost:8080/api/v1/components/filter
```

### 系统管理

```bash
# 健康检查
curl http://localhost:8080/health

# 监控指标
curl http://localhost:8080/metrics

# 系统状态
curl http://localhost:8080/api/v1/system/status
```

## 部署方式

### 1. 独立部署

```bash
# 编译
make build

# 运行
./bin/rulego-server -config configs/config.yaml
```

### 2. Docker 部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

### 3. Kubernetes 部署

```bash
# 应用配置
kubectl apply -f deploy/kubernetes/
```

## 开发指南

### 创建自定义组件

```bash
# 生成组件模板
./bin/toolkit generate component my-component -o ./components

# 实现组件逻辑
# 编辑 components/my_component.go
```

### 创建插件

```bash
# 生成插件模板
./bin/toolkit generate plugin my-plugin -o ./plugins

# 实现插件逻辑
# 编辑 plugins/my_plugin.go
```

### 创建中间件

```bash
# 生成中间件模板
./bin/toolkit generate middleware my-middleware -o ./middleware

# 实现中间件逻辑
# 编辑 middleware/my_middleware.go
```

## 工具链使用

### 项目初始化

```bash
# 创建新项目
toolkit init my-project --template basic

# 进入项目目录
cd my-project

# 编辑配置文件
vim server.yaml
```

### 零代码编译

```bash
# 从配置文件构建
toolkit build --config server.yaml --output ./dist

# 构建结果
./dist/
├── bin/my-project        # 二进制文件
├── configs/              # 配置文件
├── data/                 # 数据目录
├── logs/                 # 日志目录
├── plugins/              # 插件目录
├── start.sh             # 启动脚本
└── stop.sh              # 停止脚本
```

### 代码生成

```bash
# 生成组件
toolkit generate component transform-data -o ./components

# 生成插件
toolkit generate plugin data-processor -o ./plugins

# 生成中间件
toolkit generate middleware rate-limiter -o ./middleware
```

## 监控和日志

### Prometheus 监控

访问 `http://localhost:8080/metrics` 获取监控指标

### 健康检查

访问 `http://localhost:8080/health` 检查服务状态

### 日志配置

```yaml
logging:
  level: "info"
  format: "json"
  output: "stdout"
  filename: "logs/server.log"
  max_size: 100
  max_backups: 5
  max_age: 30
  compress: true
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 相关项目

- [RuleGo](https://github.com/rulego/rulego) - 轻量级规则引擎
- [RuleGo-Editor](https://github.com/rulego/rulego-editor) - 可视化流程编辑器

## 支持

- 📖 [文档](https://github.com/rulego/rulego-server/wiki)
- 🐛 [问题反馈](https://github.com/rulego/rulego-server/issues)
- 💬 [讨论](https://github.com/rulego/rulego-server/discussions)

## 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本更新历史