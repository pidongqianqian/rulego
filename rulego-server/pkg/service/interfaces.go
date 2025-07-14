package service

import (
	"context"
	"github.com/rulego/rulego-server/pkg/model"
)

// RuleService 规则管理服务接口
type RuleService interface {
	// 创建规则链
	CreateRuleChain(ctx context.Context, req *model.CreateRuleChainRequest) (*model.RuleChain, error)
	
	// 获取规则链
	GetRuleChain(ctx context.Context, chainId string) (*model.RuleChain, error)
	
	// 更新规则链
	UpdateRuleChain(ctx context.Context, chainId string, req *model.UpdateRuleChainRequest) error
	
	// 删除规则链
	DeleteRuleChain(ctx context.Context, chainId string) error
	
	// 部署规则链
	DeployRuleChain(ctx context.Context, chainId string) error
	
	// 取消部署
	UndeployRuleChain(ctx context.Context, chainId string) error
	
	// 执行规则链
	ExecuteRuleChain(ctx context.Context, chainId string, msg *model.RuleMessage) (*model.ExecuteResult, error)
	
	// 获取规则链列表
	ListRuleChains(ctx context.Context, req *model.ListRuleChainsRequest) (*model.ListRuleChainsResponse, error)
}

// ComponentService 组件管理服务接口
type ComponentService interface {
	// 获取所有组件
	GetComponents(ctx context.Context) ([]model.Component, error)
	
	// 获取组件详情
	GetComponent(ctx context.Context, componentType string) (*model.Component, error)
	
	// 注册组件
	RegisterComponent(ctx context.Context, component model.Component) error
	
	// 注销组件
	UnregisterComponent(ctx context.Context, componentType string) error
	
	// 获取组件表单配置
	GetComponentForm(ctx context.Context, componentType string) (*model.ComponentForm, error)
}

// PluginService 插件管理服务接口
type PluginService interface {
	// 加载插件
	LoadPlugin(ctx context.Context, pluginPath string) error
	
	// 卸载插件
	UnloadPlugin(ctx context.Context, pluginName string) error
	
	// 获取插件列表
	ListPlugins(ctx context.Context) ([]model.Plugin, error)
	
	// 获取插件详情
	GetPlugin(ctx context.Context, pluginName string) (*model.Plugin, error)
}

// SystemService 系统服务接口
type SystemService interface {
	// 获取系统状态
	GetStatus(ctx context.Context) (*model.SystemStatus, error)
	
	// 获取系统配置
	GetConfig(ctx context.Context) (*model.SystemConfig, error)
	
	// 更新系统配置
	UpdateConfig(ctx context.Context, config *model.SystemConfig) error
	
	// 重新加载配置
	Reload(ctx context.Context) error
}

// EventService 事件服务接口
type EventService interface {
	// 发布事件
	PublishEvent(ctx context.Context, event *model.Event) error
	
	// 订阅事件
	SubscribeEvent(ctx context.Context, eventType string, callback func(*model.Event)) error
	
	// 取消订阅
	UnsubscribeEvent(ctx context.Context, eventType string) error
}

// MetricsService 指标服务接口
type MetricsService interface {
	// 记录指标
	RecordMetric(ctx context.Context, metric *model.Metric) error
	
	// 获取指标
	GetMetrics(ctx context.Context, query *model.MetricsQuery) (*model.MetricsResponse, error)
}

// LogService 日志服务接口
type LogService interface {
	// 记录日志
	Log(ctx context.Context, level model.LogLevel, message string, fields map[string]interface{}) error
	
	// 查询日志
	QueryLogs(ctx context.Context, query *model.LogQuery) (*model.LogResponse, error)
}

// AuthService 认证服务接口
type AuthService interface {
	// 登录
	Login(ctx context.Context, username, password string) (*model.LoginResponse, error)
	
	// 验证token
	ValidateToken(ctx context.Context, token string) (*model.User, error)
	
	// 刷新token
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error)
	
	// 登出
	Logout(ctx context.Context, token string) error
}

// UserService 用户服务接口
type UserService interface {
	// 创建用户
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	
	// 获取用户
	GetUser(ctx context.Context, userId string) (*model.User, error)
	
	// 更新用户
	UpdateUser(ctx context.Context, userId string, req *model.UpdateUserRequest) error
	
	// 删除用户
	DeleteUser(ctx context.Context, userId string) error
	
	// 获取用户列表
	ListUsers(ctx context.Context, req *model.ListUsersRequest) (*model.ListUsersResponse, error)
}