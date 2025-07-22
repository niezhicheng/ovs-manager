package router

import (
	"github.com/gin-gonic/gin"
	"ovs-manager/api"
)

// RegisterBridgeRoutes 注册 OVS 交换机管理相关路由
func RegisterBridgeRoutes(rg *gin.RouterGroup) {
	rg.POST("/bridge/list", api.ListBridgesHandler)   // 获取交换机列表
	rg.POST("/bridge/add", api.AddBridgeHandler)      // 新增交换机
	rg.POST("/bridge/delete", api.DeleteBridgeHandler) // 删除交换机
	rg.POST("/set-netflow", api.SetNetFlowHandler) // 设置 NetFlow
	rg.POST("/set-sflow", api.SetSFlowHandler)     // 设置 sFlow
	rg.POST("/set-stp", api.SetStpHandler)         // 设置 STP
	rg.POST("/port/set-qos", api.SetQosHandler)    // 设置 QoS
	rg.POST("/set-rstp", api.SetRstpHandler) // 设置 RSTP
	rg.POST("/set-ipfix", api.SetIpfixHandler) // 设置 IPFIX
	rg.POST("/dump-flows", api.DumpFlowsHandler) // 查询流缓存
} 