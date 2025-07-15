package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	
	"gopkg.in/yaml.v3"
)

// Config 服务器配置
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Auth       AuthConfig       `yaml:"auth"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Logging    LoggingConfig    `yaml:"logging"`
	Components ComponentsConfig `yaml:"components"`
	Plugins    PluginsConfig    `yaml:"plugins"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name         string        `yaml:"name"`
	Port         int           `yaml:"port"`
	Version      string        `yaml:"version"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	MaxHeaderMB  int           `yaml:"max_header_mb"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type         string            `yaml:"type"`
	Connection   string            `yaml:"connection"`
	MaxOpenConns int               `yaml:"max_open_conns"`
	MaxIdleConns int               `yaml:"max_idle_conns"`
	MaxLifetime  time.Duration     `yaml:"max_lifetime"`
	Options      map[string]string `yaml:"options"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Enabled      bool          `yaml:"enabled"`
	JWTSecret    string        `yaml:"jwt_secret"`
	TokenExpiry  time.Duration `yaml:"token_expiry"`
	RefreshToken bool          `yaml:"refresh_token"`
	Algorithm    string        `yaml:"algorithm"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	MetricsEnabled     bool   `yaml:"metrics_enabled"`
	PrometheusEndpoint string `yaml:"prometheus_endpoint"`
	HealthEndpoint     string `yaml:"health_endpoint"`
	PProfEnabled       bool   `yaml:"pprof_enabled"`
	PProfEndpoint      string `yaml:"pprof_endpoint"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// ComponentsConfig 组件配置
type ComponentsConfig struct {
	Builtin    []string `yaml:"builtin"`
	Extensions []string `yaml:"extensions"`
}

// PluginsConfig 插件配置
type PluginsConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Directory string `yaml:"directory"`
}

// Default 返回默认配置
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Name:         "rulego-server",
			Port:         8080,
			Version:      "1.0.0",
			Host:         "0.0.0.0",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
			MaxHeaderMB:  1,
		},
		Database: DatabaseConfig{
			Type:         "file",
			Connection:   "./data",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
			MaxLifetime:  time.Hour,
			Options:      make(map[string]string),
		},
		Auth: AuthConfig{
			Enabled:      false,
			JWTSecret:    "your-secret-key",
			TokenExpiry:  24 * time.Hour,
			RefreshToken: true,
			Algorithm:    "HS256",
		},
		Monitoring: MonitoringConfig{
			MetricsEnabled:     true,
			PrometheusEndpoint: "/metrics",
			HealthEndpoint:     "/health",
			PProfEnabled:       false,
			PProfEndpoint:      "/debug/pprof",
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			Filename:   "logs/server.log",
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		},
		Components: ComponentsConfig{
			Builtin:    []string{"filter", "transform", "action"},
			Extensions: []string{},
		},
		Plugins: PluginsConfig{
			Enabled:   true,
			Directory: "./plugins",
		},
	}
}

// LoadFromFile 从文件加载配置
func LoadFromFile(filename string) (*Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", filename)
	}
	
	// 读取文件内容
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	// 解析配置
	config := Default()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}
	
	return config, nil
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() (*Config, error) {
	config := Default()
	
	// 从环境变量覆盖配置
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := parsePort(port); err == nil {
			config.Server.Port = p
		}
	}
	
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	
	if dbType := os.Getenv("DB_TYPE"); dbType != "" {
		config.Database.Type = dbType
	}
	
	if dbConn := os.Getenv("DB_CONNECTION"); dbConn != "" {
		config.Database.Connection = dbConn
	}
	
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.Auth.JWTSecret = jwtSecret
	}
	
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}
	
	return config, nil
}

// SaveToFile 保存配置到文件
func (c *Config) SaveToFile(filename string) error {
	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	
	// 序列化配置
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	
	// 写入文件
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	
	return nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证服务器配置
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", c.Server.Port)
	}
	
	if c.Server.Name == "" {
		return fmt.Errorf("服务器名称不能为空")
	}
	
	// 验证数据库配置
	if c.Database.Type == "" {
		return fmt.Errorf("数据库类型不能为空")
	}
	
	if c.Database.Connection == "" {
		return fmt.Errorf("数据库连接不能为空")
	}
	
	// 验证认证配置
	if c.Auth.Enabled && c.Auth.JWTSecret == "" {
		return fmt.Errorf("启用认证时JWT密钥不能为空")
	}
	
	// 验证日志配置
	validLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLevels, c.Logging.Level) {
		return fmt.Errorf("无效的日志级别: %s", c.Logging.Level)
	}
	
	validFormats := []string{"json", "text"}
	if !contains(validFormats, c.Logging.Format) {
		return fmt.Errorf("无效的日志格式: %s", c.Logging.Format)
	}
	
	return nil
}

// 辅助函数
func parsePort(s string) (int, error) {
	// 简单的端口解析实现
	port := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("无效的端口号: %s", s)
		}
		port = port*10 + int(c-'0')
	}
	if port <= 0 || port > 65535 {
		return 0, fmt.Errorf("端口号超出范围: %d", port)
	}
	return port, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}