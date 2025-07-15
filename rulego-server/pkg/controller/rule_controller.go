package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulego/rulego-server/pkg/model"
	"github.com/rulego/rulego-server/pkg/service"
)

// RuleController 规则控制器
type RuleController struct {
	ruleService service.RuleService
}

// NewRuleController 创建规则控制器
func NewRuleController(ruleService service.RuleService) *RuleController {
	return &RuleController{
		ruleService: ruleService,
	}
}

// CreateRuleChain 创建规则链
func (c *RuleController) CreateRuleChain(ctx *gin.Context) {
	var req model.CreateRuleChainRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error: &model.APIError{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	ruleChain, err := c.ruleService.CreateRuleChain(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Message: "Failed to create rule chain",
			Error: &model.APIError{
				Code:    "CREATE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusCreated, model.APIResponse{
		Success: true,
		Message: "Rule chain created successfully",
		Data:    ruleChain,
	})
}

// GetRuleChain 获取规则链
func (c *RuleController) GetRuleChain(ctx *gin.Context) {
	chainId := ctx.Param("id")
	if chainId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Chain ID is required",
			Error: &model.APIError{
				Code:    "MISSING_ID",
				Message: "Chain ID parameter is required",
			},
		})
		return
	}

	ruleChain, err := c.ruleService.GetRuleChain(ctx, chainId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Message: "Rule chain not found",
			Error: &model.APIError{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Rule chain retrieved successfully",
		Data:    ruleChain,
	})
}

// UpdateRuleChain 更新规则链
func (c *RuleController) UpdateRuleChain(ctx *gin.Context) {
	chainId := ctx.Param("id")
	if chainId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Chain ID is required",
			Error: &model.APIError{
				Code:    "MISSING_ID",
				Message: "Chain ID parameter is required",
			},
		})
		return
	}

	var req model.UpdateRuleChainRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error: &model.APIError{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	err := c.ruleService.UpdateRuleChain(ctx, chainId, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Message: "Failed to update rule chain",
			Error: &model.APIError{
				Code:    "UPDATE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Rule chain updated successfully",
	})
}

// DeleteRuleChain 删除规则链
func (c *RuleController) DeleteRuleChain(ctx *gin.Context) {
	chainId := ctx.Param("id")
	if chainId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Chain ID is required",
			Error: &model.APIError{
				Code:    "MISSING_ID",
				Message: "Chain ID parameter is required",
			},
		})
		return
	}

	err := c.ruleService.DeleteRuleChain(ctx, chainId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Message: "Failed to delete rule chain",
			Error: &model.APIError{
				Code:    "DELETE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Rule chain deleted successfully",
	})
}

// DeployRuleChain 部署规则链
func (c *RuleController) DeployRuleChain(ctx *gin.Context) {
	chainId := ctx.Param("id")
	if chainId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Chain ID is required",
			Error: &model.APIError{
				Code:    "MISSING_ID",
				Message: "Chain ID parameter is required",
			},
		})
		return
	}

	err := c.ruleService.DeployRuleChain(ctx, chainId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Message: "Failed to deploy rule chain",
			Error: &model.APIError{
				Code:    "DEPLOY_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Rule chain deployed successfully",
	})
}

// UndeployRuleChain 取消部署规则链
func (c *RuleController) UndeployRuleChain(ctx *gin.Context) {
	chainId := ctx.Param("id")
	if chainId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Chain ID is required",
			Error: &model.APIError{
				Code:    "MISSING_ID",
				Message: "Chain ID parameter is required",
			},
		})
		return
	}

	err := c.ruleService.UndeployRuleChain(ctx, chainId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Message: "Failed to undeploy rule chain",
			Error: &model.APIError{
				Code:    "UNDEPLOY_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Rule chain undeployed successfully",
	})
}

// ExecuteRuleChain 执行规则链
func (c *RuleController) ExecuteRuleChain(ctx *gin.Context) {
	chainId := ctx.Param("id")
	if chainId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Chain ID is required",
			Error: &model.APIError{
				Code:    "MISSING_ID",
				Message: "Chain ID parameter is required",
			},
		})
		return
	}

	var msg model.RuleMessage
	if err := ctx.ShouldBindJSON(&msg); err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid message format",
			Error: &model.APIError{
				Code:    "INVALID_MESSAGE",
				Message: err.Error(),
			},
		})
		return
	}

	result, err := c.ruleService.ExecuteRuleChain(ctx, chainId, &msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Message: "Failed to execute rule chain",
			Error: &model.APIError{
				Code:    "EXECUTE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Rule chain executed successfully",
		Data:    result,
	})
}

// ListRuleChains 获取规则链列表
func (c *RuleController) ListRuleChains(ctx *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")
	tags := ctx.QueryArray("tags")

	req := &model.ListRuleChainsRequest{
		Page:    page,
		Size:    size,
		Keyword: keyword,
		Tags:    tags,
	}

	// 解析状态参数
	if status != "" {
		statusEnum := model.RuleChainStatus(status)
		req.Status = &statusEnum
	}

	response, err := c.ruleService.ListRuleChains(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Message: "Failed to list rule chains",
			Error: &model.APIError{
				Code:    "LIST_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Rule chains retrieved successfully",
		Data:    response,
	})
}

// GetRuleChainStatus 获取规则链状态
func (c *RuleController) GetRuleChainStatus(ctx *gin.Context) {
	chainId := ctx.Param("id")
	if chainId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Chain ID is required",
			Error: &model.APIError{
				Code:    "MISSING_ID",
				Message: "Chain ID parameter is required",
			},
		})
		return
	}

	status, err := c.ruleService.GetRuleChainStatus(ctx, chainId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Message: "Rule chain status not found",
			Error: &model.APIError{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Rule chain status retrieved successfully",
		Data:    status,
	})
}

// GetNodeStatus 获取节点状态
func (c *RuleController) GetNodeStatus(ctx *gin.Context) {
	chainId := ctx.Param("chainId")
	nodeId := ctx.Param("nodeId")
	
	if chainId == "" || nodeId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Chain ID and Node ID are required",
			Error: &model.APIError{
				Code:    "MISSING_PARAMS",
				Message: "Chain ID and Node ID parameters are required",
			},
		})
		return
	}

	status, err := c.ruleService.GetNodeStatus(ctx, chainId, nodeId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Message: "Node status not found",
			Error: &model.APIError{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Node status retrieved successfully",
		Data:    status,
	})
}