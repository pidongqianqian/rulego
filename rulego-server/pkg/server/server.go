package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rulego/rulego-server/pkg/config"
	"github.com/rulego/rulego-server/pkg/controller"
	"github.com/rulego/rulego-server/pkg/dao"
	"github.com/rulego/rulego-server/pkg/service"
)

// Server 规则引擎服务器
type Server struct {
	config     *config.Config
	router     *gin.Engine
	httpServer *http.Server
	
	// 服务层
	ruleService      service.RuleService
	componentService service.ComponentService
	pluginService    service.PluginService
	
	// 数据访问层
	ruleDAO dao.RuleDAO
	
	// 控制器
	ruleController      *controller.RuleController
	componentController *controller.ComponentController
	systemController    *controller.SystemController
}

// New 创建新的服务器实例
func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// SetRuleDAO 设置自定义的规则DAO
func (s *Server) SetRuleDAO(dao dao.RuleDAO) {
	s.ruleDAO = dao
}

// Start 启动服务器
func (s *Server) Start(ctx context.Context) error {
	// 初始化组件
	if err := s.initialize(); err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}
	
	// 设置路由
	s.setupRoutes()
	
	// 创建HTTP服务器
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Port),
		Handler: s.router,
	}
	
	// 启动服务器
	go func() {
		log.Printf("Starting server on port %d", s.config.Server.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	
	// 等待信号
	s.waitForShutdown(ctx)
	
	return nil
}

// Stop 停止服务器
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// initialize 初始化服务器组件
func (s *Server) initialize() error {
	// 初始化DAO层
	if err := s.initializeDAO(); err != nil {
		return fmt.Errorf("failed to initialize DAO: %w", err)
	}
	
	// 初始化服务层
	if err := s.initializeServices(); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}
	
	// 初始化控制器
	if err := s.initializeControllers(); err != nil {
		return fmt.Errorf("failed to initialize controllers: %w", err)
	}
	
	return nil
}

// initializeDAO 初始化数据访问层
func (s *Server) initializeDAO() error {
	if s.ruleDAO == nil {
		// 使用默认的文件DAO
		s.ruleDAO = dao.NewFileRuleDAO(s.config.Data.Directory)
	}
	return nil
}

// initializeServices 初始化服务层
func (s *Server) initializeServices() error {
	// 创建规则服务
	s.ruleService = service.NewRuleService(s.ruleDAO, s.config)
	
	// 创建组件服务
	s.componentService = service.NewComponentService(s.config)
	
	// 创建插件服务
	s.pluginService = service.NewPluginService(s.config)
	
	return nil
}

// initializeControllers 初始化控制器
func (s *Server) initializeControllers() error {
	// 创建规则控制器
	s.ruleController = controller.NewRuleController(s.ruleService)
	
	// 创建组件控制器
	s.componentController = controller.NewComponentController(s.componentService)
	
	// 创建系统控制器
	s.systemController = controller.NewSystemController(s.config)
	
	return nil
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 设置gin模式
	if s.config.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	
	s.router = gin.New()
	
	// 添加中间件
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.router.Use(s.corsMiddleware())
	
	// 添加认证中间件
	if s.config.Auth.Enabled {
		s.router.Use(s.authMiddleware())
	}
	
	// API路由组
	api := s.router.Group("/api/v1")
	{
		// 规则链管理
		rules := api.Group("/rule-chains")
		{
			rules.POST("", s.ruleController.CreateRuleChain)
			rules.GET("/:chainId", s.ruleController.GetRuleChain)
			rules.PUT("/:chainId", s.ruleController.UpdateRuleChain)
			rules.DELETE("/:chainId", s.ruleController.DeleteRuleChain)
			rules.POST("/:chainId/deploy", s.ruleController.DeployRuleChain)
			rules.POST("/:chainId/undeploy", s.ruleController.UndeployRuleChain)
			rules.POST("/:chainId/execute", s.ruleController.ExecuteRuleChain)
			rules.GET("", s.ruleController.ListRuleChains)
		}
		
		// 组件管理
		components := api.Group("/components")
		{
			components.GET("", s.componentController.GetComponents)
			components.GET("/:componentType", s.componentController.GetComponent)
			components.POST("", s.componentController.RegisterComponent)
			components.DELETE("/:componentType", s.componentController.UnregisterComponent)
		}
		
		// 系统管理
		system := api.Group("/system")
		{
			system.GET("/status", s.systemController.GetStatus)
			system.GET("/config", s.systemController.GetConfig)
			system.PUT("/config", s.systemController.UpdateConfig)
			system.POST("/reload", s.systemController.Reload)
		}
	}
	
	// 健康检查
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
	
	// 静态文件服务
	if s.config.Static.Enabled {
		s.router.Static("/static", s.config.Static.Directory)
	}
}

// corsMiddleware CORS中间件
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// authMiddleware 认证中间件
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过健康检查和登录接口
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/api/v1/auth/login" {
			c.Next()
			return
		}
		
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		
		// 验证JWT token
		token := authHeader[7:] // 移除 "Bearer " 前缀
		if !s.isValidToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// isValidToken 验证JWT token
func (s *Server) isValidToken(token string) bool {
	// 这里应该实现JWT token验证逻辑
	// 为了简化示例，这里只是简单检查
	return token != ""
}

// waitForShutdown 等待关闭信号
func (s *Server) waitForShutdown(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	select {
	case <-quit:
		log.Println("Shutdown signal received")
	case <-ctx.Done():
		log.Println("Context cancelled")
	}
	
	// 优雅关闭
	log.Println("Shutting down server...")
	
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := s.Stop(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	
	log.Println("Server stopped")
}