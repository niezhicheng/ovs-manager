package router

import (
	"github.com/gin-gonic/gin"
	"ovs-manager/api"
)

// RegisterBondRoutes 注册 Bond 管理相关路由
func RegisterBondRoutes(rg *gin.RouterGroup) {
	rg.POST("/bond/add", api.AddBondHandler)   // 新增 Bond 端口
	rg.POST("/bond/set", api.SetBondHandler)   // 设置 Bond 属性
	rg.POST("/bond/show", api.ShowBondHandler) // 查询 Bond 状态
	rg.POST("/bond/delete", api.DeleteBondHandler) // 删除 Bond 端口
} 