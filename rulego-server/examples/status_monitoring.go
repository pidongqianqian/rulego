package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego-server/pkg/model"
	"github.com/rulego/rulego-server/pkg/service"
)

func main() {
	// 创建规则服务
	ruleService := service.NewRuleService()

	// 创建上下文
	ctx := context.Background()

	// 示例1: 创建规则链
	fmt.Println("=== 创建规则链 ===")
	ruleChain, err := createRuleChain(ctx, ruleService)
	if err != nil {
		log.Fatalf("创建规则链失败: %v", err)
	}
	fmt.Printf("规则链创建成功: %s\n", ruleChain.ID)

	// 示例2: 部署规则链
	fmt.Println("\n=== 部署规则链 ===")
	err = ruleService.DeployRuleChain(ctx, ruleChain.ID)
	if err != nil {
		log.Fatalf("部署规则链失败: %v", err)
	}
	fmt.Println("规则链部署成功")

	// 示例3: 执行规则链
	fmt.Println("\n=== 执行规则链 ===")
	for i := 0; i < 5; i++ {
		err = executeRuleChain(ctx, ruleService, ruleChain.ID, i)
		if err != nil {
			log.Printf("执行规则链失败: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 示例4: 获取规则链状态
	fmt.Println("\n=== 获取规则链状态 ===")
	status, err := ruleService.GetRuleChainStatus(ctx, ruleChain.ID)
	if err != nil {
		log.Fatalf("获取规则链状态失败: %v", err)
	}
	printRuleChainStatus(status)

	// 示例5: 获取节点状态
	fmt.Println("\n=== 获取节点状态 ===")
	for _, node := range status.Nodes {
		nodeStatus, err := ruleService.GetNodeStatus(ctx, ruleChain.ID, node.ID)
		if err != nil {
			log.Printf("获取节点状态失败: %v", err)
			continue
		}
		printNodeStatus(nodeStatus)
	}

	// 示例6: 清理
	fmt.Println("\n=== 清理资源 ===")
	err = ruleService.DeleteRuleChain(ctx, ruleChain.ID)
	if err != nil {
		log.Printf("删除规则链失败: %v", err)
	} else {
		fmt.Println("规则链删除成功")
	}
}

// createRuleChain 创建示例规则链
func createRuleChain(ctx context.Context, ruleService *service.RuleServiceImpl) (*model.RuleChain, error) {
	req := &model.CreateRuleChainRequest{
		Name:        "demo-chain",
		Description: "演示规则链",
		Definition: &model.RuleChainDefinition{
			RuleChain: types.RuleChain{
				Id:   "demo-chain",
				Name: "演示规则链",
				Root: true,
			},
			Metadata: types.Metadata{
				Nodes: []types.RuleNode{
					{
						Id:   "filter",
						Type: "filter",
						Name: "温度过滤",
						Configuration: types.Configuration{
							"script": "return msg.temperature > 30;",
						},
					},
					{
						Id:   "transform",
						Type: "jsTransform",
						Name: "数据转换",
						Configuration: types.Configuration{
							"jsScript": `
msg.temperature = msg.temperature + 1;
msg.processed = true;
return {msg: msg, metadata: metadata, msgType: msgType};
`,
						},
					},
					{
						Id:   "log",
						Type: "log",
						Name: "日志记录",
						Configuration: types.Configuration{
							"script": "print('Processed message:', JSON.stringify(msg));",
						},
					},
				},
				Connections: []types.RuleNodeConnection{
					{
						FromId: "filter",
						ToId:   "transform",
						Type:   "True",
					},
					{
						FromId: "transform",
						ToId:   "log",
						Type:   "Success",
					},
				},
			},
		},
		Tags: []string{"demo", "monitoring"},
		Metadata: map[string]interface{}{
			"owner":       "demo",
			"environment": "test",
		},
	}

	return ruleService.CreateRuleChain(ctx, req)
}

// executeRuleChain 执行规则链
func executeRuleChain(ctx context.Context, ruleService *service.RuleServiceImpl, chainId string, index int) error {
	msg := &model.RuleMessage{
		ID:       fmt.Sprintf("msg-%d", index),
		Type:     "TELEMETRY",
		DataType: "JSON",
		Data:     fmt.Sprintf(`{"temperature": %d, "humidity": 60}`, 25+index*5),
		Metadata: map[string]string{
			"deviceId": "sensor-001",
			"location": "room-101",
		},
		Ts: time.Now().UnixMilli(),
	}

	result, err := ruleService.ExecuteRuleChain(ctx, chainId, msg)
	if err != nil {
		return err
	}

	fmt.Printf("消息 %s 执行成功，耗时: %v\n", msg.ID, result.Duration)
	return nil
}

// printRuleChainStatus 打印规则链状态
func printRuleChainStatus(status *model.RuleChainStatusInfo) {
	fmt.Printf("规则链状态:\n")
	fmt.Printf("  ID: %s\n", status.ID)
	fmt.Printf("  名称: %s\n", status.Name)
	fmt.Printf("  状态: %s\n", status.Status)
	fmt.Printf("  总消息数: %d\n", status.TotalMessages)
	fmt.Printf("  成功数: %d\n", status.SuccessCount)
	fmt.Printf("  失败数: %d\n", status.FailedCount)
	fmt.Printf("  当前活跃: %d\n", status.CurrentActive)
	fmt.Printf("  成功率: %.2f%%\n", status.SuccessRate)
	fmt.Printf("  节点数: %d\n", len(status.Nodes))
	fmt.Printf("  最后更新: %s\n", status.LastUpdated.Format("2006-01-02 15:04:05"))
}

// printNodeStatus 打印节点状态
func printNodeStatus(status *model.NodeStatus) {
	fmt.Printf("节点状态:\n")
	fmt.Printf("  ID: %s\n", status.ID)
	fmt.Printf("  类型: %s\n", status.Type)
	fmt.Printf("  名称: %s\n", status.Name)
	fmt.Printf("  状态: %s\n", status.Status)
	fmt.Printf("  路由数: %d\n", len(status.Routes))
	fmt.Printf("  父节点数: %d\n", len(status.Parents))
	
	if len(status.Routes) > 0 {
		fmt.Printf("  路由信息:\n")
		for _, route := range status.Routes {
			fmt.Printf("    -> %s (%s)\n", route.ToID, route.RelationType)
		}
	}
	
	if len(status.Parents) > 0 {
		fmt.Printf("  父节点: %v\n", status.Parents)
	}
	
	if len(status.Metadata) > 0 {
		fmt.Printf("  元数据: %v\n", status.Metadata)
	}
	fmt.Println()
}

// 辅助函数：格式化JSON输出
func prettyPrint(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}