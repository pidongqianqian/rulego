# RuleGo Server API 使用指南

## 概述

RuleGo Server 提供了完整的 REST API 来管理规则链和节点状态。本文档介绍如何使用这些 API。

## 基础信息

- **基础URL**: `http://localhost:8080/api/v1`
- **内容类型**: `application/json`
- **认证**: 目前未启用，生产环境建议启用

## API 端点

### 1. 规则链管理

#### 1.1 创建规则链

```bash
POST /api/v1/rule-chains
```

**请求示例**:
```json
{
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
        },
        {
          "id": "alert",
          "type": "restApiCall",
          "name": "发送告警",
          "configuration": {
            "restEndpointUrlPattern": "http://alert-service/api/alerts",
            "requestMethod": "POST"
          }
        }
      ],
      "connections": [
        {
          "fromId": "filter",
          "toId": "alert",
          "type": "True"
        }
      ]
    }
  },
  "tags": ["monitoring", "temperature"],
  "metadata": {
    "owner": "admin",
    "environment": "production"
  }
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "Rule chain created successfully",
  "data": {
    "id": "temperature-monitor",
    "name": "temperature-monitor",
    "description": "温度监控规则链",
    "status": "draft",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "version": 1,
    "tags": ["monitoring", "temperature"]
  }
}
```

#### 1.2 获取规则链列表

```bash
GET /api/v1/rule-chains?page=1&size=10&keyword=temperature&status=deployed&tags=monitoring
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `size`: 每页大小 (默认: 10)
- `keyword`: 搜索关键词
- `status`: 状态过滤 (draft, deployed, disabled)
- `tags`: 标签过滤 (可多个)

**响应示例**:
```json
{
  "success": true,
  "message": "Rule chains retrieved successfully",
  "data": {
    "items": [
      {
        "id": "temperature-monitor",
        "name": "temperature-monitor",
        "description": "温度监控规则链",
        "status": "deployed",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z",
        "version": 1,
        "tags": ["monitoring", "temperature"]
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

#### 1.3 获取规则链详情

```bash
GET /api/v1/rule-chains/{id}
```

**响应示例**:
```json
{
  "success": true,
  "message": "Rule chain retrieved successfully",
  "data": {
    "id": "temperature-monitor",
    "name": "temperature-monitor",
    "description": "温度监控规则链",
    "definition": {
      "ruleChain": {
        "id": "temperature-monitor",
        "name": "温度监控",
        "root": true
      },
      "metadata": {
        "nodes": [...],
        "connections": [...]
      }
    },
    "status": "deployed",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "version": 1,
    "tags": ["monitoring", "temperature"]
  }
}
```

#### 1.4 更新规则链

```bash
PUT /api/v1/rule-chains/{id}
```

**请求示例**:
```json
{
  "description": "更新后的温度监控规则链",
  "tags": ["monitoring", "temperature", "updated"]
}
```

#### 1.5 删除规则链

```bash
DELETE /api/v1/rule-chains/{id}
```

### 2. 规则链部署

#### 2.1 部署规则链

```bash
POST /api/v1/rule-chains/{id}/deploy
```

#### 2.2 取消部署

```bash
POST /api/v1/rule-chains/{id}/undeploy
```

### 3. 规则链执行

#### 3.1 执行规则链

```bash
POST /api/v1/rule-chains/{id}/execute
```

**请求示例**:
```json
{
  "id": "msg-001",
  "type": "TELEMETRY",
  "dataType": "JSON",
  "data": "{\"temperature\": 35, \"humidity\": 60}",
  "metadata": {
    "deviceId": "sensor-001",
    "location": "room-101"
  },
  "ts": 1640995200000
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "Rule chain executed successfully",
  "data": {
    "success": true,
    "message": "Rule chain executed successfully",
    "data": "{\"temperature\": 35, \"humidity\": 60}",
    "duration": "1.234ms",
    "traceId": "msg-001",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

### 4. 状态监控

#### 4.1 获取规则链状态

```bash
GET /api/v1/rule-chains/{id}/status
```

**响应示例**:
```json
{
  "success": true,
  "message": "Rule chain status retrieved successfully",
  "data": {
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
      },
      {
        "id": "alert",
        "type": "restApiCall",
        "name": "发送告警",
        "status": "active",
        "routes": [],
        "parents": ["filter"],
        "metadata": {}
      }
    ],
    "lastUpdated": "2024-01-01T00:00:00Z"
  }
}
```

#### 4.2 获取节点状态

```bash
GET /api/v1/rule-chains/{chainId}/nodes/{nodeId}/status
```

**响应示例**:
```json
{
  "success": true,
  "message": "Node status retrieved successfully",
  "data": {
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
}
```

## 错误处理

所有 API 都返回统一的错误格式：

```json
{
  "success": false,
  "message": "错误描述",
  "error": {
    "code": "ERROR_CODE",
    "message": "详细错误信息",
    "details": "额外错误详情"
  }
}
```

### 常见错误码

- `INVALID_REQUEST`: 请求格式错误
- `MISSING_ID`: 缺少必需的ID参数
- `NOT_FOUND`: 资源不存在
- `CREATE_FAILED`: 创建失败
- `UPDATE_FAILED`: 更新失败
- `DELETE_FAILED`: 删除失败
- `DEPLOY_FAILED`: 部署失败
- `EXECUTE_FAILED`: 执行失败

## 使用示例

### 完整工作流程

1. **创建规则链**
```bash
curl -X POST http://localhost:8080/api/v1/rule-chains \
  -H "Content-Type: application/json" \
  -d @rule-chain.json
```

2. **部署规则链**
```bash
curl -X POST http://localhost:8080/api/v1/rule-chains/temperature-monitor/deploy
```

3. **执行规则链**
```bash
curl -X POST http://localhost:8080/api/v1/rule-chains/temperature-monitor/execute \
  -H "Content-Type: application/json" \
  -d '{
    "id": "msg-001",
    "type": "TELEMETRY",
    "dataType": "JSON",
    "data": "{\"temperature\": 35}"
  }'
```

4. **监控状态**
```bash
curl http://localhost:8080/api/v1/rule-chains/temperature-monitor/status
```

5. **查看节点状态**
```bash
curl http://localhost:8080/api/v1/rule-chains/temperature-monitor/nodes/filter/status
```

## 监控和指标

### 健康检查

```bash
curl http://localhost:8080/health
```

### Prometheus 指标

```bash
curl http://localhost:8080/metrics
```

## 最佳实践

1. **规则链命名**: 使用有意义的名称，包含业务领域信息
2. **标签管理**: 合理使用标签进行分类和管理
3. **状态监控**: 定期检查规则链和节点状态
4. **错误处理**: 实现适当的错误处理和重试机制
5. **性能优化**: 监控执行时间和成功率，及时优化规则链

## 注意事项

1. 规则链创建后默认为 `draft` 状态，需要手动部署
2. 部署后的规则链才能接收和执行消息
3. 节点状态信息会实时更新
4. 建议在生产环境中启用认证和HTTPS
5. 定期备份规则链配置