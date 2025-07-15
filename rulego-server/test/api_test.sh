#!/bin/bash

# RuleGo Server API 测试脚本
# 使用方法: ./api_test.sh [server_url]
# 默认服务器地址: http://localhost:8080

SERVER_URL=${1:-"http://localhost:8080"}
API_BASE="$SERVER_URL/api/v1"

echo "=== RuleGo Server API 测试 ==="
echo "服务器地址: $SERVER_URL"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "${YELLOW}测试: $description${NC}"
    echo "请求: $method $API_BASE$endpoint"
    
    if [ -n "$data" ]; then
        echo "数据: $data"
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_BASE$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            "$API_BASE$endpoint")
    fi
    
    # 分离响应体和状态码
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    echo "状态码: $http_code"
    echo "响应: $response_body"
    
    if [[ $http_code -ge 200 && $http_code -lt 300 ]]; then
        echo -e "${GREEN}✓ 成功${NC}"
    else
        echo -e "${RED}✗ 失败${NC}"
    fi
    echo ""
}

# 健康检查
echo "=== 1. 健康检查 ==="
test_api "GET" "/health" "" "健康检查"

# 创建规则链
echo "=== 2. 创建规则链 ==="
RULE_CHAIN_DATA='{
  "name": "test-chain",
  "description": "测试规则链",
  "definition": {
    "ruleChain": {
      "id": "test-chain",
      "name": "测试规则链",
      "root": true
    },
    "metadata": {
      "nodes": [
        {
          "id": "filter",
          "type": "filter",
          "name": "过滤节点",
          "configuration": {
            "script": "return msg.temperature > 30;"
          }
        },
        {
          "id": "transform",
          "type": "jsTransform",
          "name": "转换节点",
          "configuration": {
            "jsScript": "msg.temperature = msg.temperature + 1; return {msg: msg, metadata: metadata, msgType: msgType};"
          }
        }
      ],
      "connections": [
        {
          "fromId": "filter",
          "toId": "transform",
          "type": "True"
        }
      ]
    }
  },
  "tags": ["test", "demo"],
  "metadata": {
    "owner": "tester",
    "environment": "test"
  }
}'

test_api "POST" "/rule-chains" "$RULE_CHAIN_DATA" "创建测试规则链"

# 获取规则链列表
echo "=== 3. 获取规则链列表 ==="
test_api "GET" "/rule-chains" "" "获取规则链列表"

# 获取规则链详情
echo "=== 4. 获取规则链详情 ==="
test_api "GET" "/rule-chains/test-chain" "" "获取规则链详情"

# 部署规则链
echo "=== 5. 部署规则链 ==="
test_api "POST" "/rule-chains/test-chain/deploy" "" "部署规则链"

# 获取规则链状态
echo "=== 6. 获取规则链状态 ==="
test_api "GET" "/rule-chains/test-chain/status" "" "获取规则链状态"

# 执行规则链
echo "=== 7. 执行规则链 ==="
EXECUTE_DATA='{
  "id": "msg-001",
  "type": "TELEMETRY",
  "dataType": "JSON",
  "data": "{\"temperature\": 35, \"humidity\": 60}",
  "metadata": {
    "deviceId": "sensor-001",
    "location": "room-101"
  },
  "ts": 1640995200000
}'

test_api "POST" "/rule-chains/test-chain/execute" "$EXECUTE_DATA" "执行规则链"

# 获取节点状态
echo "=== 8. 获取节点状态 ==="
test_api "GET" "/rule-chains/test-chain/nodes/filter/status" "" "获取过滤节点状态"
test_api "GET" "/rule-chains/test-chain/nodes/transform/status" "" "获取转换节点状态"

# 更新规则链
echo "=== 9. 更新规则链 ==="
UPDATE_DATA='{
  "description": "更新后的测试规则链",
  "tags": ["test", "demo", "updated"]
}'

test_api "PUT" "/rule-chains/test-chain" "$UPDATE_DATA" "更新规则链"

# 再次获取规则链详情
echo "=== 10. 再次获取规则链详情 ==="
test_api "GET" "/rule-chains/test-chain" "" "获取更新后的规则链详情"

# 取消部署
echo "=== 11. 取消部署 ==="
test_api "POST" "/rule-chains/test-chain/undeploy" "" "取消部署规则链"

# 删除规则链
echo "=== 12. 删除规则链 ==="
test_api "DELETE" "/rule-chains/test-chain" "" "删除规则链"

# 验证删除
echo "=== 13. 验证删除 ==="
test_api "GET" "/rule-chains/test-chain" "" "验证规则链已删除"

echo "=== 测试完成 ==="
echo ""
echo "测试总结:"
echo "- 如果所有测试都显示绿色 ✓，说明API工作正常"
echo "- 如果有红色 ✗，请检查服务器状态和API实现"
echo "- 建议查看服务器日志以获取更多调试信息"