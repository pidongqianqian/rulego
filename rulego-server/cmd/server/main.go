package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rulego/rulego-server/pkg/config"
	"github.com/rulego/rulego-server/pkg/server"
)

var (
	configFile = flag.String("config", "configs/config.yaml", "配置文件路径")
	port       = flag.Int("port", 0, "服务器端口 (覆盖配置文件)")
	version    = flag.Bool("version", false, "显示版本信息")
	help       = flag.Bool("help", false, "显示帮助信息")
)

// 版本信息 (编译时注入)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	flag.Parse()

	// 显示版本信息
	if *version {
		fmt.Printf("RuleGo Server %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		return
	}

	// 显示帮助信息
	if *help {
		printHelp()
		return
	}

	// 加载配置
	cfg, err := config.LoadFromFile(*configFile)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 命令行参数覆盖配置
	if *port > 0 {
		cfg.Server.Port = *port
	}

	// 设置版本信息
	cfg.Server.Version = Version

	// 创建服务器
	srv := server.New(cfg)

	// 设置信号处理
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动信号监听
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		log.Printf("接收到信号 %s，开始关闭服务器...", sig)
		cancel()
	}()

	// 启动服务器
	log.Printf("启动 RuleGo Server v%s", Version)
	log.Printf("配置文件: %s", *configFile)
	log.Printf("监听端口: %d", cfg.Server.Port)

	if err := srv.Start(); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}

	// 等待关闭信号
	<-ctx.Done()
	log.Println("正在关闭服务器...")

	// 优雅关闭服务器
	if err := srv.Stop(); err != nil {
		log.Printf("关闭服务器时出错: %v", err)
	}

	log.Println("服务器已关闭")
}

func printHelp() {
	fmt.Println("RuleGo Server - 基于RuleGo的规则引擎服务器")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  rulego-server [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -config string")
	fmt.Println("        配置文件路径 (默认: configs/config.yaml)")
	fmt.Println("  -port int")
	fmt.Println("        服务器端口，覆盖配置文件中的设置")
	fmt.Println("  -version")
	fmt.Println("        显示版本信息")
	fmt.Println("  -help")
	fmt.Println("        显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  rulego-server")
	fmt.Println("  rulego-server -config /path/to/config.yaml")
	fmt.Println("  rulego-server -port 9090")
	fmt.Println("  rulego-server -config custom.yaml -port 8888")
	fmt.Println()
	fmt.Println("更多信息请访问: https://github.com/rulego/rulego-server")
}
