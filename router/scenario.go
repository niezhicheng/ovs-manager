package router

import (
	"github.com/gin-gonic/gin"
	"ovs-manager/api"
)

// RegisterScenarioRoutes 注册场景引导相关路由
func RegisterScenarioRoutes(r *gin.Engine) {
	r.POST("/api/ovs/scenario/apply", api.ScenarioApplyHandler)
} 