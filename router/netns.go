package router

import (
	"github.com/gin-gonic/gin"
	"ovs-manager/api"
)

// RegisterNetnsRoutes 注册网络命名空间相关路由
func RegisterNetnsRoutes(r *gin.Engine) {
	r.POST("/api/netns/create", api.CreateNetnsHandler)   // 创建命名空间
	r.POST("/api/netns/delete", api.DeleteNetnsHandler)   // 删除命名空间
	r.POST("/api/netns/list", api.ListNetnsHandler)       // 获取命名空间列表
} 