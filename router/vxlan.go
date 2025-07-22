package router

import (
	"github.com/gin-gonic/gin"
	"ovs-manager/api"
)

// RegisterVxlanRoutes 注册 VXLAN 端口相关路由
func RegisterVxlanRoutes(rg *gin.RouterGroup) {
	rg.POST("/vxlan/add", api.AddVxlanPortHandler) // 新增 VXLAN 端口
	rg.POST("/vxlan/delete", api.DeleteVxlanPortHandler) // 删除 VXLAN 端口
} 