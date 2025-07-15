package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego-server/pkg/model"
)

// RuleServiceImpl 规则服务实现
type RuleServiceImpl struct {
	ruleEnginePool *rulego.RuleGo
	ruleChains     map[string]*model.RuleChain
	mu             sync.RWMutex
}

// NewRuleService 创建规则服务实例
func NewRuleService() *RuleServiceImpl {
	return &RuleServiceImpl{
		ruleEnginePool: rulego.NewRuleGo(),
		ruleChains:     make(map[string]*model.RuleChain),
	}
}

// CreateRuleChain 创建规则链
func (s *RuleServiceImpl) CreateRuleChain(ctx context.Context, req *model.CreateRuleChainRequest) (*model.RuleChain, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否已存在
	if _, exists := s.ruleChains[req.Name]; exists {
		return nil, fmt.Errorf("rule chain with name %s already exists", req.Name)
	}

	// 创建规则引擎
	dsl, err := json.Marshal(req.Definition)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rule chain definition: %w", err)
	}
	
	_, err = s.ruleEnginePool.New(req.Name, dsl)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule engine: %w", err)
	}

	// 创建规则链记录
	ruleChain := &model.RuleChain{
		ID:          req.Name,
		Name:        req.Name,
		Description: req.Description,
		Definition:  req.Definition,
		Status:      model.RuleChainStatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   "system", // TODO: 从认证上下文获取
		UpdatedBy:   "system",
		Version:     1,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
	}

	s.ruleChains[req.Name] = ruleChain
	return ruleChain, nil
}

// GetRuleChain 获取规则链
func (s *RuleServiceImpl) GetRuleChain(ctx context.Context, chainId string) (*model.RuleChain, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ruleChain, exists := s.ruleChains[chainId]
	if !exists {
		return nil, fmt.Errorf("rule chain %s not found", chainId)
	}

	// 获取引擎实例
	engine, exists := s.ruleEnginePool.Get(chainId)
	if !exists {
		return nil, fmt.Errorf("rule engine %s not found", chainId)
	}

	// 更新状态信息
	ruleChain.Status = s.getEngineStatus(engine)
	ruleChain.UpdatedAt = time.Now()

	return ruleChain, nil
}

// UpdateRuleChain 更新规则链
func (s *RuleServiceImpl) UpdateRuleChain(ctx context.Context, chainId string, req *model.UpdateRuleChainRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ruleChain, exists := s.ruleChains[chainId]
	if !exists {
		return fmt.Errorf("rule chain %s not found", chainId)
	}

	// 更新字段
	if req.Name != nil {
		ruleChain.Name = *req.Name
	}
	if req.Description != nil {
		ruleChain.Description = *req.Description
	}
	if req.Definition != nil {
		ruleChain.Definition = req.Definition
		ruleChain.Version++
	}
	if req.Tags != nil {
		ruleChain.Tags = req.Tags
	}
	if req.Metadata != nil {
		ruleChain.Metadata = req.Metadata
	}

	ruleChain.UpdatedAt = time.Now()
	ruleChain.UpdatedBy = "system" // TODO: 从认证上下文获取

	// 重新加载规则引擎
	if req.Definition != nil {
		engine, exists := s.ruleEnginePool.Get(chainId)
		if exists {
			dsl, err := json.Marshal(req.Definition)
			if err != nil {
				return fmt.Errorf("failed to marshal rule chain definition: %w", err)
			}
			err = engine.ReloadSelf(dsl)
			if err != nil {
				return fmt.Errorf("failed to reload rule engine: %w", err)
			}
		}
	}

	return nil
}

// DeleteRuleChain 删除规则链
func (s *RuleServiceImpl) DeleteRuleChain(ctx context.Context, chainId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 删除规则引擎
	s.ruleEnginePool.Del(chainId)

	// 删除记录
	delete(s.ruleChains, chainId)

	return nil
}

// DeployRuleChain 部署规则链
func (s *RuleServiceImpl) DeployRuleChain(ctx context.Context, chainId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ruleChain, exists := s.ruleChains[chainId]
	if !exists {
		return fmt.Errorf("rule chain %s not found", chainId)
	}

	// 检查引擎是否存在
	engine, exists := s.ruleEnginePool.Get(chainId)
	if !exists {
		return fmt.Errorf("rule engine %s not found", chainId)
	}

	// 检查引擎是否已初始化
	if !engine.Initialized() {
		return fmt.Errorf("rule engine %s not initialized", chainId)
	}

	ruleChain.Status = model.RuleChainStatusDeployed
	ruleChain.UpdatedAt = time.Now()
	ruleChain.UpdatedBy = "system"

	return nil
}

// UndeployRuleChain 取消部署
func (s *RuleServiceImpl) UndeployRuleChain(ctx context.Context, chainId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ruleChain, exists := s.ruleChains[chainId]
	if !exists {
		return fmt.Errorf("rule chain %s not found", chainId)
	}

	ruleChain.Status = model.RuleChainStatusDisabled
	ruleChain.UpdatedAt = time.Now()
	ruleChain.UpdatedBy = "system"

	return nil
}

// ExecuteRuleChain 执行规则链
func (s *RuleServiceImpl) ExecuteRuleChain(ctx context.Context, chainId string, msg *model.RuleMessage) (*model.ExecuteResult, error) {
	engine, exists := s.ruleEnginePool.Get(chainId)
	if !exists {
		return nil, fmt.Errorf("rule engine %s not found", chainId)
	}

	// 转换消息格式
	ruleMsg := types.NewMsg(0, msg.Type, types.DataType(msg.DataType), types.NewMetadata(), msg.Data)
	for k, v := range msg.Metadata {
		ruleMsg.Metadata.PutValue(k, v)
	}

	// 执行规则链
	startTime := time.Now()
	engine.OnMsg(ruleMsg)
	duration := time.Since(startTime)

	// 转换元数据格式
	metadata := make(map[string]interface{})
	for k, v := range msg.Metadata {
		metadata[k] = v
	}

	// 获取执行结果
	result := &model.ExecuteResult{
		Success:   true, // TODO: 根据实际执行结果判断
		Message:   "Rule chain executed successfully",
		Data:      msg.Data,
		Metadata:  metadata,
		Duration:  duration,
		TraceID:   msg.ID,
		Timestamp: startTime,
	}

	return result, nil
}

// ListRuleChains 获取规则链列表
func (s *RuleServiceImpl) ListRuleChains(ctx context.Context, req *model.ListRuleChainsRequest) (*model.ListRuleChainsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var items []model.RuleChain
	for _, ruleChain := range s.ruleChains {
		// 应用过滤条件
		if req.Status != nil && ruleChain.Status != *req.Status {
			continue
		}
		if req.Keyword != "" && !contains(ruleChain.Name, req.Keyword) && !contains(ruleChain.Description, req.Keyword) {
			continue
		}
		if len(req.Tags) > 0 && !hasCommonTag(ruleChain.Tags, req.Tags) {
			continue
		}

		// 更新状态信息
		engine, exists := s.ruleEnginePool.Get(ruleChain.ID)
		if exists {
			ruleChain.Status = s.getEngineStatus(engine)
		}

		items = append(items, *ruleChain)
	}

	// 应用分页
	total := len(items)
	start := (req.Page - 1) * req.Size
	end := start + req.Size
	if start >= total {
		items = []model.RuleChain{}
	} else if end > total {
		items = items[start:]
	} else {
		items = items[start:end]
	}

	return &model.ListRuleChainsResponse{
		Items: items,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// GetRuleChainStatus 获取规则链状态
func (s *RuleServiceImpl) GetRuleChainStatus(ctx context.Context, chainId string) (*model.RuleChainStatusInfo, error) {
	engine, exists := s.ruleEnginePool.Get(chainId)
	if !exists {
		return nil, fmt.Errorf("rule engine %s not found", chainId)
	}

	// 获取引擎指标
	engineMetrics := engine.GetMetrics()
	metrics := engineMetrics.Get()

	// 获取规则链定义
	definition := engine.Definition()
	rootCtx := engine.RootRuleChainCtx()

	// 构建节点状态
	var nodeStatuses []model.NodeStatus
	if rootCtx != nil {
		for _, node := range definition.Metadata.Nodes {
			// 获取路由和父节点
			var routes []model.NodeRoute
			var parentIds []string
			for _, conn := range definition.Metadata.Connections {
				if conn.FromId == node.Id {
					routes = append(routes, model.NodeRoute{ToID: conn.ToId, RelationType: conn.Type})
				}
				if conn.ToId == node.Id {
					parentIds = append(parentIds, conn.FromId)
				}
			}
			nodeStatus := model.NodeStatus{
				ID:       node.Id,
				Type:     node.Type,
				Name:     node.Name,
				Status:   "active", // TODO: 实现更详细的状态检查
				Routes:   routes,
				Parents:  parentIds,
				Metadata: make(map[string]interface{}),
			}
			nodeStatuses = append(nodeStatuses, nodeStatus)
		}
	}

	status := &model.RuleChainStatusInfo{
		ID:            chainId,
		Name:          definition.RuleChain.Name, // 使用 RuleChain.Name
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

// GetNodeStatus 获取节点状态
func (s *RuleServiceImpl) GetNodeStatus(ctx context.Context, chainId string, nodeId string) (*model.NodeStatus, error) {
	engine, exists := s.ruleEnginePool.Get(chainId)
	if !exists {
		return nil, fmt.Errorf("rule engine %s not found", chainId)
	}

	rootCtx := engine.RootRuleChainCtx()
	if rootCtx == nil {
		return nil, fmt.Errorf("root context not found")
	}

	definition := engine.Definition()
	var nodeDef *types.RuleNode
	for _, n := range definition.Metadata.Nodes {
		if n.Id == nodeId {
			nodeDef = n
			break
		}
	}
	if nodeDef == nil {
		return nil, fmt.Errorf("node %s not found in chain %s", nodeId, chainId)
	}

	// 获取节点路由信息
	var routes []model.NodeRoute
	var parentIds []string
	for _, conn := range definition.Metadata.Connections {
		if conn.FromId == nodeId {
			routes = append(routes, model.NodeRoute{ToID: conn.ToId, RelationType: conn.Type})
		}
		if conn.ToId == nodeId {
			parentIds = append(parentIds, conn.FromId)
		}
	}

	// 构建节点状态
	nodeStatus := &model.NodeStatus{
		ID:       nodeId,
		Type:     nodeDef.Type,
		Name:     nodeDef.Name,
		Status:   "active", // TODO: 实现更详细的状态检查
		Routes:   routes,
		Parents:  parentIds,
		Metadata: make(map[string]interface{}),
	}

	return nodeStatus, nil
}

// 辅助方法

// getEngineStatus 获取引擎状态
func (s *RuleServiceImpl) getEngineStatus(engine types.RuleEngine) model.RuleChainStatus {
	if !engine.Initialized() {
		return model.RuleChainStatusDraft
	}
	return model.RuleChainStatusDeployed
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	}()))
}

// hasCommonTag 检查是否有共同标签
func hasCommonTag(tags1, tags2 []string) bool {
	for _, tag1 := range tags1 {
		for _, tag2 := range tags2 {
			if tag1 == tag2 {
				return true
			}
		}
	}
	return false
}