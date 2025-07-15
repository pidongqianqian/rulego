# RuleGo IoT数据处理规则链测试指南

本测试演示了一个完整的IoT数据处理流程，包含MQTT输入、JavaScript过滤和解析、PostgreSQL数据存储以及RabbitMQ消息发布。

## 系统架构

```
MQTT输入 → JS过滤器 → JS转换器 → 数据分发 → PostgreSQL + RabbitMQ
   ↓           ↓          ↓         ↓            ↓
 传感器数据   数据验证   数据富化   并行处理    数据存储 + 消息队列
```

## 数据流程

1. **MQTT输入**: 接收来自IoT传感器的JSON数据
2. **JS过滤器**: 验证数据格式和数值范围
3. **JS转换器**: 数据标准化、单位转换、告警级别计算
4. **数据分发**: 将处理后的数据同时发送到多个目标
5. **PostgreSQL**: 存储历史数据和统计信息
6. **RabbitMQ**: 发布实时消息和告警通知

## 前置条件

### 1. 安装必要服务

```bash
# 安装PostgreSQL
brew install postgresql
brew services start postgresql

# 安装RabbitMQ
brew install rabbitmq
brew services start rabbitmq

# 安装MQTT Broker (Mosquitto)
brew install mosquitto
brew services start mosquitto
```

### 2. 配置数据库

```bash
# 创建数据库和表结构
psql -h localhost -U postgres -d postgres -f configs/database_schema.sql
```

### 3. 配置RabbitMQ

```bash
# 启用管理插件
rabbitmq-plugins enable rabbitmq_management

# 创建交换机和队列
rabbitmqctl add_vhost /
rabbitmqctl add_user guest guest
rabbitmqctl set_permissions -p / guest ".*" ".*" ".*"

# 访问管理界面: http://localhost:15672 (guest/guest)
```

## 测试数据格式

### 输入数据格式 (MQTT)

```json
{
  "deviceId": "sensor_001",
  "timestamp": "2024-01-20T10:30:00Z",
  "data": {
    "temperature": 25.5,
    "humidity": 60.0,
    "pressure": 1013.25
  }
}
```

### 处理后数据格式 (PostgreSQL)

```json
{
  "deviceId": "sensor_001",
  "timestamp": "2024-01-20T10:30:00Z",
  "originalData": { "temperature": 25.5, "humidity": 60.0, "pressure": 1013.25 },
  "processedAt": "2024-01-20T10:30:01Z",
  "dataType": "sensor_data",
  "temperature": {
    "celsius": 25.5,
    "fahrenheit": 77.9,
    "kelvin": 298.65
  },
  "humidity": {
    "percentage": 60.0,
    "absoluteHumidity": 130.0
  },
  "pressure": {
    "hPa": 1013.25,
    "mmHg": 760.0,
    "psi": 14.7
  },
  "alertLevel": "normal",
  "location": "机房A-机柜1"
}
```

## 告警级别规则

- **normal**: 温度 ≤ 35°C 且 湿度 ≤ 80%
- **warning**: 温度 35-45°C 或 湿度 80-90%
- **critical**: 温度 > 45°C 或 湿度 > 90%

## 运行测试

### 1. 启动RuleGo服务器

```bash
cd rulego-server
go run cmd/server/main.go -config configs/config.yaml
```

### 2. 加载规则链

```bash
# 通过API加载规则链
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d @configs/iot_processing_chain.json
```

### 3. 运行测试客户端

```bash
cd test
./run_test.sh
```

或者手动运行：

```bash
cd test
go mod tidy
go build -o iot_test_client iot_test_client.go
./iot_test_client
```

## 测试场景

### 1. 正常数据测试 (normal)
- 发送正常范围内的传感器数据
- 验证数据成功存储到PostgreSQL
- 验证消息发布到RabbitMQ

### 2. 警告级别测试 (warning)
- 发送温度36-44°C或湿度81-89%的数据
- 验证告警级别被正确设置为"warning"
- 验证告警消息发布到专门的告警队列

### 3. 严重级别测试 (critical)
- 发送温度>45°C或湿度>90%的数据
- 验证告警级别被正确设置为"critical"
- 验证严重告警消息处理

### 4. 混合数据测试 (mixed)
- 随机发送正常和异常数据
- 验证系统能正确处理各种情况

### 5. 无效数据测试 (invalid)
- 发送缺少必要字段的数据
- 发送超出范围的数值
- 验证过滤器正确拒绝无效数据

### 6. 连续测试 (continuous)
- 长时间运行混合数据模拟
- 验证系统稳定性和性能

## 验证结果

### 1. 检查PostgreSQL数据

```sql
-- 查看传感器数据
SELECT * FROM sensor_data ORDER BY timestamp DESC LIMIT 10;

-- 查看设备状态
SELECT * FROM device_latest_data;

-- 查看告警日志
SELECT * FROM alert_logs ORDER BY created_at DESC LIMIT 10;

-- 查看今日统计
SELECT * FROM today_stats;
```

### 2. 检查RabbitMQ消息

```bash
# 查看队列状态
rabbitmqctl list_queues

# 查看交换机
rabbitmqctl list_exchanges

# 查看绑定关系
rabbitmqctl list_bindings
```

### 3. 监控日志

```bash
# 查看RuleGo服务器日志
tail -f logs/server.log

# 查看MQTT代理日志
tail -f /usr/local/var/log/mosquitto/mosquitto.log
```

## 性能指标

- **处理延迟**: < 100ms per message
- **吞吐量**: > 1000 messages/second
- **数据准确性**: 100% (无数据丢失)
- **告警响应时间**: < 1 second

## 故障排除

### 1. MQTT连接问题
```bash
# 测试MQTT连接
mosquitto_pub -h localhost -p 1883 -t "test/topic" -m "hello"
mosquitto_sub -h localhost -p 1883 -t "test/topic"
```

### 2. PostgreSQL连接问题
```bash
# 测试数据库连接
psql -h localhost -U postgres -d rulego_test -c "SELECT version();"
```

### 3. RabbitMQ连接问题
```bash
# 检查RabbitMQ状态
rabbitmqctl status
```

### 4. 规则链调试
- 在规则链配置中启用 `debugMode: true`
- 查看详细的处理日志
- 使用RuleGo的调试工具

## 扩展功能

### 1. 添加新的传感器类型
- 修改JS过滤器以支持新的数据格式
- 更新数据库表结构
- 调整告警规则

### 2. 集成其他消息队列
- 添加Kafka支持
- 集成Redis发布/订阅
- 支持NATS消息传递

### 3. 增强数据处理
- 添加数据聚合功能
- 实现时间窗口统计
- 支持复杂事件处理

### 4. 监控和告警
- 集成Prometheus监控
- 添加Grafana仪表板
- 实现邮件/短信告警

## 注意事项

1. **数据格式**: 确保MQTT消息格式正确
2. **时间戳**: 使用ISO 8601格式的时间戳
3. **设备ID**: 设备ID应该唯一且有意义
4. **数据范围**: 传感器数据应在合理范围内
5. **网络延迟**: 考虑网络延迟对实时性的影响

## 联系支持

如果遇到问题，请：
1. 查看日志文件
2. 检查服务状态
3. 验证配置文件
4. 联系技术支持团队 