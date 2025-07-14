package main

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rulego/rulego-server/internal/toolkit"
)

var (
	version = "1.0.0"
	cfgFile string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "toolkit",
		Short: "RuleGo Server Toolkit",
		Long:  `RuleGo Server Toolkit - 无代码编译工具链`,
		Version: version,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径")
	
	// 添加子命令
	rootCmd.AddCommand(
		newInitCommand(),
		newBuildCommand(),
		newGenerateCommand(),
		newValidateCommand(),
		newUpgradeCommand(),
		newServeCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initCommand 初始化项目
func newInitCommand() *cobra.Command {
	var projectName string
	var template string
	
	cmd := &cobra.Command{
		Use:   "init",
		Short: "初始化新项目",
		Long:  `初始化一个新的RuleGo Server项目`,
		Run: func(cmd *cobra.Command, args []string) {
			if projectName == "" {
				fmt.Println("项目名称不能为空")
				os.Exit(1)
			}
			
			generator := toolkit.NewProjectGenerator()
			if err := generator.Init(projectName, template); err != nil {
				fmt.Printf("初始化项目失败: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("项目 %s 初始化成功\n", projectName)
		},
	}
	
	cmd.Flags().StringVarP(&projectName, "name", "n", "", "项目名称")
	cmd.Flags().StringVarP(&template, "template", "t", "default", "项目模板")
	cmd.MarkFlagRequired("name")
	
	return cmd
}

// buildCommand 编译项目
func newBuildCommand() *cobra.Command {
	var configFile string
	var outputDir string
	var tags []string
	var ldflags string
	var race bool
	var verbose bool
	
	cmd := &cobra.Command{
		Use:   "build",
		Short: "编译服务器",
		Long:  `根据配置文件编译RuleGo Server`,
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				fmt.Println("配置文件不能为空")
				os.Exit(1)
			}
			
			builder := toolkit.NewBuilder()
			config, err := builder.LoadConfig(configFile)
			if err != nil {
				fmt.Printf("加载配置失败: %v\n", err)
				os.Exit(1)
			}
			
			buildOptions := &toolkit.BuildOptions{
				ConfigFile: configFile,
				OutputDir:  outputDir,
				Tags:       tags,
				Ldflags:    ldflags,
				Race:       race,
				Verbose:    verbose,
			}
			
			if err := builder.Build(config, buildOptions); err != nil {
				fmt.Printf("编译失败: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("编译成功，输出目录: %s\n", outputDir)
		},
	}
	
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "配置文件路径")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "./dist", "输出目录")
	cmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "构建标签")
	cmd.Flags().StringVar(&ldflags, "ldflags", "", "链接器标志")
	cmd.Flags().BoolVar(&race, "race", false, "启用竞态检测")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")
	cmd.MarkFlagRequired("config")
	
	return cmd
}

// generateCommand 生成代码
func newGenerateCommand() *cobra.Command {
	var genType string
	var name string
	var outputDir string
	
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "生成代码",
		Long:  `生成组件、插件等代码模板`,
		Run: func(cmd *cobra.Command, args []string) {
			if genType == "" {
				fmt.Println("生成类型不能为空")
				os.Exit(1)
			}
			
			if name == "" {
				fmt.Println("名称不能为空")
				os.Exit(1)
			}
			
			generator := toolkit.NewCodeGenerator()
			if err := generator.Generate(genType, name, outputDir); err != nil {
				fmt.Printf("生成失败: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("%s %s 生成成功\n", genType, name)
		},
	}
	
	cmd.Flags().StringVar(&genType, "type", "", "生成类型 (component/plugin/middleware)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "名称")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "./generated", "输出目录")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("name")
	
	return cmd
}

// validateCommand 验证配置
func newValidateCommand() *cobra.Command {
	var configFile string
	
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "验证配置",
		Long:  `验证配置文件的正确性`,
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				fmt.Println("配置文件不能为空")
				os.Exit(1)
			}
			
			validator := toolkit.NewValidator()
			if err := validator.Validate(configFile); err != nil {
				fmt.Printf("配置验证失败: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println("配置验证通过")
		},
	}
	
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "配置文件路径")
	cmd.MarkFlagRequired("config")
	
	return cmd
}

// upgradeCommand 升级
func newUpgradeCommand() *cobra.Command {
	var targetVersion string
	var force bool
	
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "升级服务器",
		Long:  `升级RuleGo Server到指定版本`,
		Run: func(cmd *cobra.Command, args []string) {
			upgrader := toolkit.NewUpgrader()
			if err := upgrader.Upgrade(targetVersion, force); err != nil {
				fmt.Printf("升级失败: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("升级到版本 %s 成功\n", targetVersion)
		},
	}
	
	cmd.Flags().StringVarP(&targetVersion, "version", "v", "", "目标版本")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "强制升级")
	cmd.MarkFlagRequired("version")
	
	return cmd
}

// serveCommand 启动服务
func newServeCommand() *cobra.Command {
	var configFile string
	var daemon bool
	var pidFile string
	
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启动服务",
		Long:  `启动RuleGo Server服务`,
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				fmt.Println("配置文件不能为空")
				os.Exit(1)
			}
			
			server := toolkit.NewServer()
			if err := server.Start(configFile, daemon, pidFile); err != nil {
				fmt.Printf("启动服务失败: %v\n", err)
				os.Exit(1)
			}
		},
	}
	
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "配置文件路径")
	cmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "后台运行")
	cmd.Flags().StringVar(&pidFile, "pid", "", "PID文件路径")
	cmd.MarkFlagRequired("config")
	
	return cmd
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rulego-toolkit")
	}
	
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("使用配置文件:", viper.ConfigFileUsed())
	}
}