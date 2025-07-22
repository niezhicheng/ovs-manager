package router

import (
	"github.com/gin-gonic/gin"
	"ovs-manager/api"
)

// RegisterFlowRoutes 注册流表规则相关路由
func RegisterFlowRoutes(rg *gin.RouterGroup) {
	rg.POST("/flow/list-v2", api.ListFlowsV2Handler)     // 查询流表规则
	rg.POST("/flow/add-v2", api.AddFlowV2Handler)       // 添加流表规则
	rg.POST("/flow/delete-v2", api.DeleteFlowV2Handler) // 删除流表规则
} 