# RuleGo Server API 状态监控功能

## 概述

本项目为 RuleGo Server 实现了完整的 API 状态监控功能，允许通过 REST API 获取规则链和节点的实时状态信息。

## 功能特性

### 1. 规则链管理
- ✅ 创建规则链
- ✅ 获取规则链列表
- ✅ 获取规则链详情
- ✅ 更新规则链
- ✅ 删除规则链
- ✅ 部署/取消部署规则链
- ✅ 执行规则链

### 2. 状态监控
- ✅ 获取规则链状态（包含执行指标）
- ✅ 获取节点状态（包含路由信息）
- ✅ 实时指标收集（成功/失败/活跃数）
- ✅ 节点关系映射

### 3. 指标统计
- ✅ 总处理消息数
- ✅ 成功/失败计数
- ✅ 当前活跃执行数
- ✅ 成功率计算
- ✅ 执行时间统计

## 快速开始

### 1. 启动服务器

```bash
cd rulego-server
make build
./bin/rulego-server -config configs/config.yaml
```

### 2. 运行测试

```bash
cd test
./api_test.sh
```

### 3. 使用API

#### 创建规则链
```bash
curl -X POST http://localhost:8080/api/v1/rule-chains \
  -H "Content-Type: application/json" \
  -d '{
    "name": "temperature-monitor",
    "description": "温度监控规则链",
    "definition": {
      "ruleChain": {
        "id": "temperature-monitor",
        "name": "温度监控",
        "root": true
      },
      "metadata": {
        "nodes": [
          {
            "id": "filter",
            "type": "filter",
            "name": "温度过滤",
            "configuration": {
              "script": "return msg.temperature > 30;"
            }
          }
        ],
        "connections": []
      }
    }
  }'
```

#### 获取规则链状态
```bash
curl http://localhost:8080/api/v1/rule-chains/temperature-monitor/status
```

#### 获取节点状态
```bash
curl http://localhost:8080/api/v1/rule-chains/temperature-monitor/nodes/filter/status
```

## API 端点

### 规则链管理
- `POST /api/v1/rule-chains` - 创建规则链
- `GET /api/v1/rule-chains` - 获取规则链列表
- `GET /api/v1/rule-chains/{id}` - 获取规则链详情
- `PUT /api/v1/rule-chains/{id}` - 更新规则链
- `DELETE /api/v1/rule-chains/{id}` - 删除规则链

### 规则链操作
- `POST /api/v1/rule-chains/{id}/deploy` - 部署规则链
- `POST /api/v1/rule-chains/{id}/undeploy` - 取消部署
- `POST /api/v1/rule-chains/{id}/execute` - 执行规则链

### 状态监控
- `GET /api/v1/rule-chains/{id}/status` - 获取规则链状态
- `GET /api/v1/rule-chains/{chainId}/nodes/{nodeId}/status` - 获取节点状态

## 状态信息结构

### 规则链状态
```json
{
  "id": "temperature-monitor",
  "name": "温度监控",
  "status": "deployed",
  "totalMessages": 1000,
  "successCount": 950,
  "failedCount": 50,
  "currentActive": 5,
  "successRate": 95.0,
  "nodes": [
    {
      "id": "filter",
      "type": "filter",
      "name": "温度过滤",
      "status": "active",
      "routes": [
        {
          "toId": "alert",
          "relationType": "True"
        }
      ],
      "parents": [],
      "metadata": {}
    }
  ],
  "lastUpdated": "2024-01-01T00:00:00Z"
}
```

### 节点状态
```json
{
  "id": "filter",
  "type": "filter",
  "name": "温度过滤",
  "status": "active",
  "routes": [
    {
      "toId": "alert",
      "relationType": "True"
    }
  ],
  "parents": [],
  "metadata": {
    "processedMessages": 1000,
    "filteredMessages": 200
  }
}
```

## 实现原理

### 1. 服务层架构
- **RuleServiceImpl**: 核心业务逻辑实现
- **状态管理**: 基于 RuleGo 引擎的指标收集
- **并发安全**: 使用读写锁保证线程安全

### 2. 状态收集机制
- **引擎指标**: 通过 `engine.GetMetrics()` 获取执行统计
- **节点遍历**: 遍历规则链上下文获取节点信息
- **路由映射**: 构建节点间的连接关系

### 3. 数据模型
- **RuleChainStatusInfo**: 规则链状态信息
- **NodeStatus**: 节点状态信息
- **NodeRoute**: 节点路由信息

## 核心代码

### 状态获取实现
```go
func (s *RuleServiceImpl) GetRuleChainStatus(ctx context.Context, chainId string) (*model.RuleChainStatusInfo, error) {
    engine, exists := s.ruleEnginePool.Get(chainId)
    if !exists {
        return nil, fmt.Errorf("rule engine %s not found", chainId)
    }

    // 获取引擎指标
    engineMetrics := engine.GetMetrics()
    metrics := engineMetrics.Get()

    // 获取规则链定义和节点信息
    definition := engine.Definition()
    rootCtx := engine.RootRuleChainCtx()

    // 构建节点状态列表
    var nodeStatuses []model.NodeStatus
    if rootCtx != nil {
        for _, nodeId := range rootCtx.(*engine.RuleChainCtx).nodeIds {
            node, exists := rootCtx.(*engine.RuleChainCtx).GetNodeById(nodeId)
            if exists {
                nodeStatus := model.NodeStatus{
                    ID:       string(nodeId),
                    Type:     node.Type(),
                    Name:     node.GetNodeId(),
                    Status:   "active",
                    Metadata: make(map[string]interface{}),
                }
                nodeStatuses = append(nodeStatuses, nodeStatus)
            }
        }
    }

    status := &model.RuleChainStatusInfo{
        ID:            chainId,
        Name:          definition.Name,
        Status:        s.getEngineStatus(engine),
        TotalMessages: metrics.Total,
        SuccessCount:  metrics.Success,
        FailedCount:   metrics.Failed,
        CurrentActive: metrics.Current,
        SuccessRate:   float64(metrics.Success) / float64(metrics.Total) * 100,
        Nodes:         nodeStatuses,
        LastUpdated:   time.Now(),
    }

    return status, nil
}
```

## 监控和告警

### 1. 关键指标
- **成功率**: 低于 95% 时告警
- **活跃执行数**: 超过阈值时告警
- **执行时间**: 超过预期时告警

### 2. 监控集成
- **Prometheus**: 通过 `/metrics` 端点暴露指标
- **健康检查**: 通过 `/health` 端点检查服务状态
- **日志记录**: 详细的操作日志

## 最佳实践

### 1. 性能优化
- 使用缓存减少重复计算
- 异步更新状态信息
- 批量获取节点状态

### 2. 可靠性保证
- 错误重试机制
- 状态一致性检查
- 优雅降级处理

### 3. 扩展性设计
- 插件化状态收集器
- 可配置的指标聚合
- 支持自定义状态字段

## 故障排除

### 常见问题

1. **状态获取失败**
   - 检查规则引擎是否已初始化
   - 验证规则链ID是否正确
   - 查看服务器日志

2. **指标不准确**
   - 确认指标切面已启用
   - 检查并发访问是否正常
   - 验证时间窗口设置

3. **节点状态缺失**
   - 检查节点ID是否存在
   - 验证规则链定义完整性
   - 确认节点类型支持

### 调试方法

1. **启用调试日志**
```yaml
logging:
  level: "debug"
```

2. **检查健康状态**
```bash
curl http://localhost:8080/health
```

3. **查看指标数据**
```bash
curl http://localhost:8080/metrics
```

## 未来规划

### 1. 功能增强
- [ ] 实时状态推送（WebSocket）
- [ ] 历史状态查询
- [ ] 状态变更通知
- [ ] 自定义指标收集

### 2. 性能优化
- [ ] 状态缓存机制
- [ ] 增量状态更新
- [ ] 分布式状态同步

### 3. 监控增强
- [ ] 告警规则配置
- [ ] 状态趋势分析
- [ ] 性能瓶颈检测

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交代码变更
4. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证。