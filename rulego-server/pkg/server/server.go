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
	"github.com/rulego/rulego-server/pkg/service"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	router := gin.Default()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	s := &Server{
		config: cfg,
		router: router,
		server: server,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// 创建服务实例
	ruleService := service.NewRuleService()
	
	// 创建控制器实例
	ruleController := controller.NewRuleController(ruleService)

	// Health check endpoint
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Rule chain management endpoints
		v1.POST("/rule-chains", ruleController.CreateRuleChain)
		v1.GET("/rule-chains", ruleController.ListRuleChains)
		v1.GET("/rule-chains/:id", ruleController.GetRuleChain)
		v1.PUT("/rule-chains/:id", ruleController.UpdateRuleChain)
		v1.DELETE("/rule-chains/:id", ruleController.DeleteRuleChain)
		
		// Rule chain deployment endpoints
		v1.POST("/rule-chains/:id/deploy", ruleController.DeployRuleChain)
		v1.POST("/rule-chains/:id/undeploy", ruleController.UndeployRuleChain)
		
		// Rule chain execution endpoint
		v1.POST("/rule-chains/:id/execute", ruleController.ExecuteRuleChain)
		
		// Rule chain status endpoints
		v1.GET("/rule-chains/:id/status", ruleController.GetRuleChainStatus)
		v1.GET("/rule-chains/:chainId/nodes/:nodeId/status", ruleController.GetNodeStatus)

		// Legacy endpoints for backward compatibility
		v1.POST("/rules", s.createRule)
		v1.GET("/rules", s.listRules)
		v1.GET("/rules/:id", s.getRule)
		v1.PUT("/rules/:id", s.updateRule)
		v1.DELETE("/rules/:id", s.deleteRule)

		// Component management endpoints
		v1.GET("/components", s.listComponents)
		v1.GET("/components/:type", s.getComponent)
	}
}

// Placeholder handlers - these would be implemented with actual business logic
func (s *Server) createRule(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Rule created successfully"})
}

func (s *Server) listRules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"rules": []string{}})
}

func (s *Server) getRule(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Rule details"})
}

func (s *Server) updateRule(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Rule updated successfully"})
}

func (s *Server) deleteRule(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Rule deleted successfully"})
}

func (s *Server) listComponents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"components": []string{}})
}

func (s *Server) getComponent(c *gin.Context) {
	componentType := c.Param("type")
	c.JSON(http.StatusOK, gin.H{"type": componentType, "message": "Component details"})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting server on port %d", s.config.Server.Port)

	// Start server in a goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	return s.Stop()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server stopped")
	return nil
}
