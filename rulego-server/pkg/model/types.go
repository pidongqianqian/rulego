package model

import (
	"time"

	"github.com/rulego/rulego/api/types"
)

// RuleChain 规则链
type RuleChain struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Definition  *RuleChainDefinition   `json:"definition"`
	Status      RuleChainStatus        `json:"status"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	CreatedBy   string                 `json:"createdBy"`
	UpdatedBy   string                 `json:"updatedBy"`
	Version     int                    `json:"version"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// RuleChainDefinition 规则链定义
type RuleChainDefinition struct {
	RuleChain types.RuleChain `json:"ruleChain"`
	Metadata  types.Metadata  `json:"metadata"`
}

// RuleChainStatus 规则链状态
type RuleChainStatus string

const (
	RuleChainStatusDraft    RuleChainStatus = "draft"
	RuleChainStatusDeployed RuleChainStatus = "deployed"
	RuleChainStatusDisabled RuleChainStatus = "disabled"
)

// CreateRuleChainRequest 创建规则链请求
type CreateRuleChainRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Definition  *RuleChainDefinition   `json:"definition" binding:"required"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateRuleChainRequest 更新规则链请求
type UpdateRuleChainRequest struct {
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Definition  *RuleChainDefinition   `json:"definition"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ListRuleChainsRequest 获取规则链列表请求
type ListRuleChainsRequest struct {
	Page     int              `json:"page" form:"page"`
	Size     int              `json:"size" form:"size"`
	Keyword  string           `json:"keyword" form:"keyword"`
	Status   *RuleChainStatus `json:"status" form:"status"`
	Tags     []string         `json:"tags" form:"tags"`
	SortBy   string           `json:"sortBy" form:"sortBy"`
	SortDesc bool             `json:"sortDesc" form:"sortDesc"`
}

// ListRuleChainsResponse 获取规则链列表响应
type ListRuleChainsResponse struct {
	Items []RuleChain `json:"items"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// RuleMessage 规则消息
type RuleMessage struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	DataType string            `json:"dataType"`
	Data     string            `json:"data"`
	Metadata map[string]string `json:"metadata"`
	Ts       int64             `json:"ts"`
}

// ExecuteResult 执行结果
type ExecuteResult struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	Duration  time.Duration          `json:"duration"`
	TraceID   string                 `json:"traceId"`
	Timestamp time.Time              `json:"timestamp"`
}

// Component 组件
type Component struct {
	Type        string             `json:"type"`
	Name        string             `json:"name"`
	Category    string             `json:"category"`
	Description string             `json:"description"`
	Icon        string             `json:"icon"`
	Version     string             `json:"version"`
	Author      string             `json:"author"`
	Config      *ComponentConfig   `json:"config"`
	Form        *ComponentForm     `json:"form"`
	Examples    []ComponentExample `json:"examples"`
	Tags        []string           `json:"tags"`
	Deprecated  bool               `json:"deprecated"`
}

// ComponentConfig 组件配置
type ComponentConfig struct {
	Fields      []ConfigField          `json:"fields"`
	Required    []string               `json:"required"`
	Defaults    map[string]interface{} `json:"defaults"`
	Validations map[string]interface{} `json:"validations"`
}

// ConfigField 配置字段
type ConfigField struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     interface{} `json:"default"`
	Required    bool        `json:"required"`
	Options     []Option    `json:"options"`
	Validation  *Validation `json:"validation"`
}

// Option 选项
type Option struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

// Validation 验证规则
type Validation struct {
	Min     *int    `json:"min"`
	Max     *int    `json:"max"`
	Pattern *string `json:"pattern"`
	Message string  `json:"message"`
}

// ComponentForm 组件表单
type ComponentForm struct {
	Schema   map[string]interface{} `json:"schema"`
	UISchema map[string]interface{} `json:"uiSchema"`
}

// ComponentExample 组件示例
type ComponentExample struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

// Plugin 插件
type Plugin struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	Path        string                 `json:"path"`
	Status      PluginStatus           `json:"status"`
	LoadedAt    time.Time              `json:"loadedAt"`
	Config      map[string]interface{} `json:"config"`
}

// PluginStatus 插件状态
type PluginStatus string

const (
	PluginStatusLoaded   PluginStatus = "loaded"
	PluginStatusUnloaded PluginStatus = "unloaded"
	PluginStatusError    PluginStatus = "error"
)

// SystemStatus 系统状态
type SystemStatus struct {
	Status        string                 `json:"status"`
	Uptime        time.Duration          `json:"uptime"`
	Version       string                 `json:"version"`
	Memory        MemoryUsage            `json:"memory"`
	CPU           CPUUsage               `json:"cpu"`
	Goroutines    int                    `json:"goroutines"`
	RuleChains    int                    `json:"ruleChains"`
	Components    int                    `json:"components"`
	Plugins       int                    `json:"plugins"`
	Connections   int                    `json:"connections"`
	LastHeartbeat time.Time              `json:"lastHeartbeat"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// MemoryUsage 内存使用情况
type MemoryUsage struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"totalAlloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"numGC"`
}

// CPUUsage CPU使用情况
type CPUUsage struct {
	Usage     float64 `json:"usage"`
	LoadAvg1  float64 `json:"loadAvg1"`
	LoadAvg5  float64 `json:"loadAvg5"`
	LoadAvg15 float64 `json:"loadAvg15"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	Server      ServerConfig      `json:"server"`
	Database    DatabaseConfig    `json:"database"`
	Auth        AuthConfig        `json:"auth"`
	Logging     LoggingConfig     `json:"logging"`
	Monitoring  MonitoringConfig  `json:"monitoring"`
	Components  ComponentsConfig  `json:"components"`
	Plugins     PluginsConfig     `json:"plugins"`
	Performance PerformanceConfig `json:"performance"`
	Security    SecurityConfig    `json:"security"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host          string        `json:"host"`
	Port          int           `json:"port"`
	ReadTimeout   time.Duration `json:"readTimeout"`
	WriteTimeout  time.Duration `json:"writeTimeout"`
	IdleTimeout   time.Duration `json:"idleTimeout"`
	MaxHeaderSize int           `json:"maxHeaderSize"`
	Debug         bool          `json:"debug"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type         string        `json:"type"`
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Database     string        `json:"database"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	MaxOpenConns int           `json:"maxOpenConns"`
	MaxIdleConns int           `json:"maxIdleConns"`
	MaxLifetime  time.Duration `json:"maxLifetime"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Enabled    bool          `json:"enabled"`
	JWTSecret  string        `json:"jwtSecret"`
	ExpireTime time.Duration `json:"expireTime"`
	Issuer     string        `json:"issuer"`
	Algorithm  string        `json:"algorithm"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	MaxSize    int    `json:"maxSize"`
	MaxAge     int    `json:"maxAge"`
	MaxBackups int    `json:"maxBackups"`
	Compress   bool   `json:"compress"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled           bool   `json:"enabled"`
	MetricsEnabled    bool   `json:"metricsEnabled"`
	PrometheusEnabled bool   `json:"prometheusEnabled"`
	PrometheusPath    string `json:"prometheusPath"`
	HealthCheckPath   string `json:"healthCheckPath"`
}

// ComponentsConfig 组件配置
type ComponentsConfig struct {
	LoadBuiltin    bool     `json:"loadBuiltin"`
	LoadExtensions bool     `json:"loadExtensions"`
	ExtensionPaths []string `json:"extensionPaths"`
	Whitelist      []string `json:"whitelist"`
	Blacklist      []string `json:"blacklist"`
}

// PluginsConfig 插件配置
type PluginsConfig struct {
	Enabled    bool     `json:"enabled"`
	Directory  string   `json:"directory"`
	AutoReload bool     `json:"autoReload"`
	Whitelist  []string `json:"whitelist"`
	Blacklist  []string `json:"blacklist"`
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	MaxWorkers      int           `json:"maxWorkers"`
	WorkerQueueSize int           `json:"workerQueueSize"`
	MaxConnections  int           `json:"maxConnections"`
	RequestTimeout  time.Duration `json:"requestTimeout"`
	CacheSize       int           `json:"cacheSize"`
	CacheTTL        time.Duration `json:"cacheTTL"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EnableHTTPS    bool     `json:"enableHTTPS"`
	CertFile       string   `json:"certFile"`
	KeyFile        string   `json:"keyFile"`
	AllowedOrigins []string `json:"allowedOrigins"`
	AllowedMethods []string `json:"allowedMethods"`
	AllowedHeaders []string `json:"allowedHeaders"`
	RateLimitRPS   int      `json:"rateLimitRPS"`
	RateLimitBurst int      `json:"rateLimitBurst"`
}

// Event 事件
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

// Metric 指标
type Metric struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Value     float64                `json:"value"`
	Labels    map[string]string      `json:"labels"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// MetricsQuery 指标查询
type MetricsQuery struct {
	Names       []string          `json:"names"`
	Labels      map[string]string `json:"labels"`
	StartTime   time.Time         `json:"startTime"`
	EndTime     time.Time         `json:"endTime"`
	Step        time.Duration     `json:"step"`
	Aggregation string            `json:"aggregation"`
}

// MetricsResponse 指标响应
type MetricsResponse struct {
	Metrics []Metric `json:"metrics"`
	Total   int      `json:"total"`
}

// LogLevel 日志级别
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogQuery 日志查询
type LogQuery struct {
	Level     *LogLevel `json:"level"`
	Message   string    `json:"message"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Page      int       `json:"page"`
	Size      int       `json:"size"`
}

// LogResponse 日志响应
type LogResponse struct {
	Logs  []LogEntry `json:"logs"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// LogEntry 日志条目
type LogEntry struct {
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp time.Time              `json:"timestamp"`
}

// User 用户
type User struct {
	ID        string                 `json:"id"`
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	Role      string                 `json:"role"`
	Status    UserStatus             `json:"status"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	LastLogin time.Time              `json:"lastLogin"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string                 `json:"username" binding:"required"`
	Email    string                 `json:"email" binding:"required,email"`
	Password string                 `json:"password" binding:"required,min=8"`
	Role     string                 `json:"role"`
	Metadata map[string]interface{} `json:"metadata"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    *string                `json:"email"`
	Role     *string                `json:"role"`
	Status   *UserStatus            `json:"status"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ListUsersRequest 获取用户列表请求
type ListUsersRequest struct {
	Page     int         `json:"page" form:"page"`
	Size     int         `json:"size" form:"size"`
	Keyword  string      `json:"keyword" form:"keyword"`
	Role     string      `json:"role" form:"role"`
	Status   *UserStatus `json:"status" form:"status"`
	SortBy   string      `json:"sortBy" form:"sortBy"`
	SortDesc bool        `json:"sortDesc" form:"sortDesc"`
}

// ListUsersResponse 获取用户列表响应
type ListUsersResponse struct {
	Items []User `json:"items"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	User         User      `json:"user"`
}

// TokenResponse token响应
type TokenResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

// APIResponse API响应
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Error     *APIError   `json:"error"`
	Timestamp time.Time   `json:"timestamp"`
	TraceID   string      `json:"traceId"`
}

// APIError API错误
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// PageInfo 分页信息
type PageInfo struct {
	Page    int  `json:"page"`
	Size    int  `json:"size"`
	Total   int  `json:"total"`
	HasNext bool `json:"hasNext"`
	HasPrev bool `json:"hasPrev"`
}

// ErrorCode 错误代码
type ErrorCode string

const (
	ErrorCodeInvalidRequest   ErrorCode = "INVALID_REQUEST"
	ErrorCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrorCodeConflict         ErrorCode = "CONFLICT"
	ErrorCodeInternalError    ErrorCode = "INTERNAL_ERROR"
	ErrorCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrorCodeRateLimited      ErrorCode = "RATE_LIMITED"
)
