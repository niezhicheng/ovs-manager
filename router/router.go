package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化 Gin 路由，汇总各功能模块
func InitRouter() *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // 允许所有源，也可以指定多个源，如：[]string{"http://localhost:8080", "http://example.com"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

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
