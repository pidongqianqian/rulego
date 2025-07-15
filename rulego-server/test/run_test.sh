#!/bin/bash

# IoT数据处理规则链测试脚本
# 此脚本用于测试完整的IoT数据处理流程

echo "=== RuleGo IoT数据处理规则链测试 ==="
echo ""

# 检查必要的服务是否运行
check_service() {
    local service=$1
    local port=$2
    local name=$3
    
    if nc -z localhost $port 2>/dev/null; then
        echo "✅ $name 服务正在运行 (端口: $port)"
    else
        echo "❌ $name 服务未运行 (端口: $port)"
        echo "   请启动 $name 服务后再运行测试"
        return 1
    fi
}

echo "检查必要服务状态..."
check_service "mqtt" 1883 "MQTT Broker" || exit 1
check_service "postgresql" 5432 "PostgreSQL" || exit 1
check_service "rabbitmq" 5672 "RabbitMQ" || exit 1

echo ""
echo "=== 准备数据库 ==="
echo "正在创建数据库表结构..."
psql -h localhost -U postgres -d postgres -f ../configs/database_schema.sql

echo ""
echo "=== 启动测试客户端 ==="
echo "正在编译测试客户端..."
cd "$(dirname "$0")"
go mod tidy
go build -o iot_test_client iot_test_client.go

if [ $? -eq 0 ]; then
    echo "✅ 测试客户端编译成功"
    echo ""
    echo "=== 开始测试 ==="
    echo "测试客户端即将启动，请按照提示选择测试场景："
    echo "  1. normal - 发送正常传感器数据"
    echo "  2. warning - 发送警告级别数据"
    echo "  3. critical - 发送严重级别数据"
    echo "  4. mixed - 发送混合正常/异常数据"
    echo "  5. invalid - 发送无效数据测试过滤"
    echo "  6. continuous - 运行连续混合数据模拟"
    echo "  7. exit - 退出程序"
    echo ""
    echo "启动测试客户端..."
    ./iot_test_client
else
    echo "❌ 测试客户端编译失败"
    exit 1
fi 