package router

import (
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化 Gin 路由，汇总各功能模块
func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	ovs := r.Group("/api/ovs")
	RegisterBridgeRoutes(ovs)
	RegisterPortRoutes(ovs)
	RegisterMirrorRoutes(ovs)
	RegisterFlowRoutes(ovs)
	RegisterVxlanRoutes(ovs)
	RegisterBondRoutes(ovs)
	RegisterScenarioRoutes(r)

	RegisterNetnsRoutes(r)

	r.POST("/api/ovs/show", func(c *gin.Context) {
		// 详细状态接口仍在主路由注册
		// 若需拆分可移至新文件
		c.JSON(200, gin.H{"message": "TODO: OVSShowHandler"})
	})

	return r
} 