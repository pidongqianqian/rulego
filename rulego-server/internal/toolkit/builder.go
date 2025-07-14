package toolkit

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	
	"gopkg.in/yaml.v3"
)

// BuildConfig 构建配置
type BuildConfig struct {
	Server struct {
		Name    string `yaml:"name"`
		Port    int    `yaml:"port"`
		Version string `yaml:"version"`
	} `yaml:"server"`
	
	Components struct {
		Builtin    []string `yaml:"builtin"`
		Extensions []string `yaml:"extensions"`
	} `yaml:"components"`
	
	Plugins []struct {
		Path string `yaml:"path"`
		Name string `yaml:"name"`
	} `yaml:"plugins"`
	
	DAO struct {
		Type       string `yaml:"type"`
		Connection string `yaml:"connection"`
	} `yaml:"dao"`
	
	Auth struct {
		Enabled   bool   `yaml:"enabled"`
		JWTSecret string `yaml:"jwt_secret"`
	} `yaml:"auth"`
	
	Monitoring struct {
		MetricsEnabled      bool   `yaml:"metrics_enabled"`
		PrometheusEndpoint  string `yaml:"prometheus_endpoint"`
	} `yaml:"monitoring"`
}

// BuildOptions 构建选项
type BuildOptions struct {
	ConfigFile string
	OutputDir  string
	Tags       []string
	Ldflags    string
	Race       bool
	Verbose    bool
}

// Builder 构建器
type Builder struct {
	workDir string
	config  *BuildConfig
}

// NewBuilder 创建构建器
func NewBuilder() *Builder {
	return &Builder{
		workDir: "./temp_build",
	}
}

// LoadConfig 加载配置
func (b *Builder) LoadConfig(configFile string) (*BuildConfig, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	var config BuildConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	
	b.config = &config
	return &config, nil
}

// Build 构建项目
func (b *Builder) Build(config *BuildConfig, options *BuildOptions) error {
	// 创建工作目录
	if err := os.MkdirAll(b.workDir, 0755); err != nil {
		return fmt.Errorf("创建工作目录失败: %w", err)
	}
	defer os.RemoveAll(b.workDir)
	
	// 生成main.go
	if err := b.generateMainFile(config, options); err != nil {
		return fmt.Errorf("生成主文件失败: %w", err)
	}
	
	// 生成go.mod
	if err := b.generateGoMod(config, options); err != nil {
		return fmt.Errorf("生成go.mod失败: %w", err)
	}
	
	// 生成配置文件
	if err := b.generateConfigFile(config, options); err != nil {
		return fmt.Errorf("生成配置文件失败: %w", err)
	}
	
	// 编译项目
	if err := b.compileProject(config, options); err != nil {
		return fmt.Errorf("编译项目失败: %w", err)
	}
	
	// 复制文件到输出目录
	if err := b.copyToOutput(config, options); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}
	
	return nil
}

// generateMainFile 生成主文件
func (b *Builder) generateMainFile(config *BuildConfig, options *BuildOptions) error {
	tmpl := `package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/rulego/rulego-server/pkg/server"
	"github.com/rulego/rulego-server/pkg/config"
	{{- range .Components.Extensions }}
	_ "{{ . }}"
	{{- end }}
	{{- range .Plugins }}
	_ "{{ .Path }}"
	{{- end }}
)

func main() {
	// 加载配置
	cfg := config.Default()
	cfg.Server.Name = "{{ .Server.Name }}"
	cfg.Server.Port = {{ .Server.Port }}
	cfg.Server.Version = "{{ .Server.Version }}"
	
	// 设置DAO配置
	cfg.DAO.Type = "{{ .DAO.Type }}"
	cfg.DAO.Connection = "{{ .DAO.Connection }}"
	
	// 设置认证配置
	cfg.Auth.Enabled = {{ .Auth.Enabled }}
	cfg.Auth.JWTSecret = "{{ .Auth.JWTSecret }}"
	
	// 设置监控配置
	cfg.Monitoring.MetricsEnabled = {{ .Monitoring.MetricsEnabled }}
	cfg.Monitoring.PrometheusEndpoint = "{{ .Monitoring.PrometheusEndpoint }}"
	
	// 创建服务器
	srv := server.New(cfg)
	
	// 处理信号
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("接收到关闭信号")
		cancel()
	}()
	
	// 启动服务器
	log.Printf("启动 %s 服务器...", cfg.Server.Name)
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
`
	
	t, err := template.New("main").Parse(tmpl)
	if err != nil {
		return err
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, config); err != nil {
		return err
	}
	
	// 格式化代码
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(filepath.Join(b.workDir, "main.go"), formatted, 0644)
}

// generateGoMod 生成go.mod文件
func (b *Builder) generateGoMod(config *BuildConfig, options *BuildOptions) error {
	tmpl := `module {{ .Server.Name }}

go 1.19

require (
	github.com/rulego/rulego-server v1.0.0
	github.com/rulego/rulego v1.0.0
	{{- range .Components.Extensions }}
	{{ . }} v1.0.0
	{{- end }}
)
`
	
	t, err := template.New("gomod").Parse(tmpl)
	if err != nil {
		return err
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, config); err != nil {
		return err
	}
	
	return ioutil.WriteFile(filepath.Join(b.workDir, "go.mod"), buf.Bytes(), 0644)
}

// generateConfigFile 生成配置文件
func (b *Builder) generateConfigFile(config *BuildConfig, options *BuildOptions) error {
	configContent := fmt.Sprintf(`# %s 配置文件
# 自动生成于 %s

# 服务器配置
server:
  name: %s
  port: %d
  version: %s

# 数据库配置
database:
  type: %s
  connection: %s

# 认证配置
auth:
  enabled: %t
  jwt_secret: %s

# 监控配置
monitoring:
  metrics_enabled: %t
  prometheus_endpoint: %s

# 日志配置
logging:
  level: info
  format: json
  output: stdout

# 组件配置
components:
  builtin: %v
  extensions: %v

# 插件配置
plugins:
  enabled: true
  directory: ./plugins
`,
		config.Server.Name,
		time.Now().Format("2006-01-02 15:04:05"),
		config.Server.Name,
		config.Server.Port,
		config.Server.Version,
		config.DAO.Type,
		config.DAO.Connection,
		config.Auth.Enabled,
		config.Auth.JWTSecret,
		config.Monitoring.MetricsEnabled,
		config.Monitoring.PrometheusEndpoint,
		config.Components.Builtin,
		config.Components.Extensions,
	)
	
	configDir := filepath.Join(b.workDir, "configs")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	return ioutil.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0644)
}

// compileProject 编译项目
func (b *Builder) compileProject(config *BuildConfig, options *BuildOptions) error {
	// 构建命令
	args := []string{"build"}
	
	// 添加标签
	if len(options.Tags) > 0 {
		args = append(args, "-tags", strings.Join(options.Tags, ","))
	}
	
	// 添加ldflags
	if options.Ldflags != "" {
		args = append(args, "-ldflags", options.Ldflags)
	}
	
	// 添加race检测
	if options.Race {
		args = append(args, "-race")
	}
	
	// 添加输出文件
	outputFile := filepath.Join(b.workDir, "bin", config.Server.Name)
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return err
	}
	args = append(args, "-o", outputFile)
	
	args = append(args, ".")
	
	// 执行编译
	cmd := exec.Command("go", args...)
	cmd.Dir = b.workDir
	
	if options.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	
	return cmd.Run()
}

// copyToOutput 复制文件到输出目录
func (b *Builder) copyToOutput(config *BuildConfig, options *BuildOptions) error {
	outputDir := options.OutputDir
	if outputDir == "" {
		outputDir = "./dist"
	}
	
	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	
	// 复制二进制文件
	binSrc := filepath.Join(b.workDir, "bin", config.Server.Name)
	binDst := filepath.Join(outputDir, "bin", config.Server.Name)
	if err := b.copyFile(binSrc, binDst); err != nil {
		return err
	}
	
	// 复制配置文件
	configSrc := filepath.Join(b.workDir, "configs")
	configDst := filepath.Join(outputDir, "configs")
	if err := b.copyDir(configSrc, configDst); err != nil {
		return err
	}
	
	// 创建其他目录
	dirs := []string{"data", "logs", "plugins"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(outputDir, dir), 0755); err != nil {
			return err
		}
	}
	
	// 生成启动脚本
	if err := b.generateStartScript(config, outputDir); err != nil {
		return err
	}
	
	return nil
}

// generateStartScript 生成启动脚本
func (b *Builder) generateStartScript(config *BuildConfig, outputDir string) error {
	script := fmt.Sprintf(`#!/bin/bash

# %s 启动脚本

SERVER_NAME="%s"
SERVER_DIR="$(cd "$(dirname "$0")" && pwd)"
BIN_DIR="$SERVER_DIR/bin"
CONFIG_DIR="$SERVER_DIR/configs"
DATA_DIR="$SERVER_DIR/data"
LOG_DIR="$SERVER_DIR/logs"
PID_FILE="$SERVER_DIR/$SERVER_NAME.pid"

# 检查是否已经运行
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null 2>&1; then
        echo "$SERVER_NAME is already running (PID: $PID)"
        exit 1
    else
        rm -f "$PID_FILE"
    fi
fi

# 启动服务器
echo "Starting $SERVER_NAME..."
cd "$SERVER_DIR"
nohup "$BIN_DIR/$SERVER_NAME" --config="$CONFIG_DIR/config.yaml" > "$LOG_DIR/server.log" 2>&1 &
PID=$!

# 保存PID
echo $PID > "$PID_FILE"
echo "$SERVER_NAME started (PID: $PID)"

# 等待服务器启动
sleep 3
if ps -p $PID > /dev/null 2>&1; then
    echo "$SERVER_NAME is running successfully"
else
    echo "$SERVER_NAME failed to start"
    rm -f "$PID_FILE"
    exit 1
fi
`, config.Server.Name, config.Server.Name)
	
	scriptFile := filepath.Join(outputDir, "start.sh")
	if err := ioutil.WriteFile(scriptFile, []byte(script), 0755); err != nil {
		return err
	}
	
	// 生成停止脚本
	stopScript := fmt.Sprintf(`#!/bin/bash

# %s 停止脚本

SERVER_NAME="%s"
SERVER_DIR="$(cd "$(dirname "$0")" && pwd)"
PID_FILE="$SERVER_DIR/$SERVER_NAME.pid"

if [ ! -f "$PID_FILE" ]; then
    echo "$SERVER_NAME is not running"
    exit 1
fi

PID=$(cat "$PID_FILE")
if ! ps -p $PID > /dev/null 2>&1; then
    echo "$SERVER_NAME is not running"
    rm -f "$PID_FILE"
    exit 1
fi

echo "Stopping $SERVER_NAME (PID: $PID)..."
kill $PID

# 等待进程停止
for i in {1..10}; do
    if ! ps -p $PID > /dev/null 2>&1; then
        echo "$SERVER_NAME stopped"
        rm -f "$PID_FILE"
        exit 0
    fi
    sleep 1
done

# 强制停止
echo "Force stopping $SERVER_NAME..."
kill -9 $PID
rm -f "$PID_FILE"
echo "$SERVER_NAME stopped"
`, config.Server.Name, config.Server.Name)
	
	stopScriptFile := filepath.Join(outputDir, "stop.sh")
	return ioutil.WriteFile(stopScriptFile, []byte(stopScript), 0755)
}

// copyFile 复制文件
func (b *Builder) copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(dst, data, 0644)
}

// copyDir 复制目录
func (b *Builder) copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		
		dstPath := filepath.Join(dst, relPath)
		
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		
		return b.copyFile(path, dstPath)
	})
}

// ProjectGenerator 项目生成器
type ProjectGenerator struct{}

// NewProjectGenerator 创建项目生成器
func NewProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{}
}

// Init 初始化项目
func (g *ProjectGenerator) Init(name, template string) error {
	// 创建项目目录
	if err := os.MkdirAll(name, 0755); err != nil {
		return err
	}
	
	// 生成基础配置文件
	configContent := fmt.Sprintf(`# %s 配置文件
server:
  name: %s
  port: 8080
  version: "1.0.0"

components:
  builtin:
    - "filter"
    - "transform"
    - "action"
  extensions: []

plugins: []

dao:
  type: "file"
  connection: "./data"

auth:
  enabled: false
  jwt_secret: "your-secret-key"

monitoring:
  metrics_enabled: true
  prometheus_endpoint: "/metrics"
`, name, name)
	
	configFile := filepath.Join(name, "server.yaml")
	if err := ioutil.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		return err
	}
	
	// 生成README文件
	readmeContent := fmt.Sprintf(`# %s

基于RuleGo的规则引擎服务器

## 编译

使用工具链编译：

`, name) + "```bash\n" + `toolkit build --config server.yaml --output ./dist
` + "```\n\n" + `## 运行

` + "```bash\n" + `cd dist
./start.sh
` + "```\n\n" + `## 停止

` + "```bash\n" + `./stop.sh
` + "```\n"
	
	readmeFile := filepath.Join(name, "README.md")
	return ioutil.WriteFile(readmeFile, []byte(readmeContent), 0644)
}

// CodeGenerator 代码生成器
type CodeGenerator struct{}

// NewCodeGenerator 创建代码生成器
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

// Generate 生成代码
func (g *CodeGenerator) Generate(genType, name, outputDir string) error {
	switch genType {
	case "component":
		return g.generateComponent(name, outputDir)
	case "plugin":
		return g.generatePlugin(name, outputDir)
	case "middleware":
		return g.generateMiddleware(name, outputDir)
	default:
		return fmt.Errorf("不支持的生成类型: %s", genType)
	}
}

// generateComponent 生成组件
func (g *CodeGenerator) generateComponent(name, outputDir string) error {
	tmpl := `package components

import (
	"context"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/components/base"
)

// {{ .Name }}Component {{ .Name }}组件
type {{ .Name }}Component struct {
	base.BaseNode
	Config {{ .Name }}Config
}

// {{ .Name }}Config {{ .Name }}配置
type {{ .Name }}Config struct {
	// 在这里添加配置字段
}

// Type 返回组件类型
func (c *{{ .Name }}Component) Type() string {
	return "{{ .Type }}"
}

// New 创建新实例
func (c *{{ .Name }}Component) New() types.Node {
	return &{{ .Name }}Component{
		Config: {{ .Name }}Config{},
	}
}

// Init 初始化组件
func (c *{{ .Name }}Component) Init(ruleConfig types.Config, configuration types.Configuration) error {
	// 初始化配置
	err := configuration.UnmarshalConfiguration(c.Config)
	if err != nil {
		return err
	}
	
	// 在这里添加初始化逻辑
	return nil
}

// OnMsg 处理消息
func (c *{{ .Name }}Component) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	// 在这里实现消息处理逻辑
	
	// 发送到下一个节点
	ctx.TellSuccess(msg)
}

// Destroy 销毁组件
func (c *{{ .Name }}Component) Destroy() {
	// 在这里添加清理逻辑
}

// 注册组件
func init() {
	rulego.Registry.Register(&{{ .Name }}Component{})
}
`
	
	data := struct {
		Name string
		Type string
	}{
		Name: strings.Title(name),
		Type: fmt.Sprintf("custom/%s", name),
	}
	
	t, err := template.New("component").Parse(tmpl)
	if err != nil {
		return err
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}
	
	// 格式化代码
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	
	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_component.go", name))
	return ioutil.WriteFile(outputFile, formatted, 0644)
}

// generatePlugin 生成插件
func (g *CodeGenerator) generatePlugin(name, outputDir string) error {
	tmpl := `package main

import (
	"context"
	"github.com/rulego/rulego-server/pkg/plugin"
)

// {{ .Name }}Plugin {{ .Name }}插件
type {{ .Name }}Plugin struct {
	config map[string]interface{}
}

// Name 返回插件名称
func (p *{{ .Name }}Plugin) Name() string {
	return "{{ .Name }}"
}

// Version 返回插件版本
func (p *{{ .Name }}Plugin) Version() string {
	return "1.0.0"
}

// Init 初始化插件
func (p *{{ .Name }}Plugin) Init(config map[string]interface{}) error {
	p.config = config
	// 在这里添加初始化逻辑
	return nil
}

// Start 启动插件
func (p *{{ .Name }}Plugin) Start(ctx context.Context) error {
	// 在这里添加启动逻辑
	return nil
}

// Stop 停止插件
func (p *{{ .Name }}Plugin) Stop(ctx context.Context) error {
	// 在这里添加停止逻辑
	return nil
}

// Health 健康检查
func (p *{{ .Name }}Plugin) Health(ctx context.Context) error {
	// 在这里添加健康检查逻辑
	return nil
}

// 导出插件
var Plugin = &{{ .Name }}Plugin{}
`
	
	data := struct {
		Name string
	}{
		Name: strings.Title(name),
	}
	
	t, err := template.New("plugin").Parse(tmpl)
	if err != nil {
		return err
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}
	
	// 格式化代码
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	
	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_plugin.go", name))
	return ioutil.WriteFile(outputFile, formatted, 0644)
}

// generateMiddleware 生成中间件
func (g *CodeGenerator) generateMiddleware(name, outputDir string) error {
	tmpl := `package middleware

import (
	"context"
	"github.com/rulego/rulego-server/pkg/controller"
)

// {{ .Name }}Middleware {{ .Name }}中间件
func {{ .Name }}Middleware() controller.Middleware {
	return func(next controller.HandlerFunc) controller.HandlerFunc {
		return func(ctx context.Context, req *controller.Request) (*controller.Response, error) {
			// 在这里添加中间件逻辑
			
			// 调用下一个处理器
			return next(ctx, req)
		}
	}
}
`
	
	data := struct {
		Name string
	}{
		Name: strings.Title(name),
	}
	
	t, err := template.New("middleware").Parse(tmpl)
	if err != nil {
		return err
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}
	
	// 格式化代码
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	
	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_middleware.go", name))
	return ioutil.WriteFile(outputFile, formatted, 0644)
}