package router

import (
	"github.com/gin-gonic/gin"
	"ovs-manager/api"
)

// RegisterMirrorRoutes 注册端口镜像相关路由
func RegisterMirrorRoutes(rg *gin.RouterGroup) {
	rg.GET("/mirror", api.GetMirrorsHandler)           // 获取端口镜像列表 (GET)
	rg.POST("/mirror/add", api.AddMirrorHandler)       // 新增端口镜像
	rg.POST("/mirror/delete", api.DeleteMirrorHandler) // 删除端口镜像
	rg.POST("/mirror/list", api.ListMirrorsHandler)    // 查询端口镜像
} 