package router

import (
	"github.com/gin-gonic/gin"
	"ovs-manager/api"
)

// RegisterPortRoutes 注册 OVS 端口管理相关路由
func RegisterPortRoutes(rg *gin.RouterGroup) {
	rg.POST("/port/list", api.ListPortsHandler)        // 获取端口列表
	rg.POST("/port/add", api.AddPortHandler)           // 新增端口
	rg.POST("/port/delete", api.DeletePortHandler)     // 删除端口

	// 端口与命名空间绑定/解绑
	rg.POST("/port/bind-netns", api.BindPortToNetnsHandler)     // 端口绑定命名空间
	rg.POST("/port/unbind-netns", api.UnbindPortFromNetnsHandler) // 端口解绑到主命名空间

	// 端口 up/down、分配 IP
	rg.POST("/port/updown", api.SetPortUpDownHandler)           // 端口 up/down
	rg.POST("/port/addr", api.SetPortAddrHandler)               // 端口分配 IP

	// 端口 VLAN tag 设置
	rg.POST("/port/set-vlan", api.SetPortVlanTagHandler)         // 端口 VLAN tag 设置
	rg.POST("/tunnel/add", api.AddTunnelPortHandler) // 添加 GRE/Geneve Tunnel Port
	rg.POST("/set-bfd", api.SetBfdHandler) // 设置 BFD
	rg.POST("/set-cfm", api.SetCfmHandler) // 设置 CFM
	rg.POST("/bridge/set-mcast-snooping", api.SetMcastSnoopingHandler) // 设置组播监听
	rg.POST("/port/set-hfsc-qos", api.SetHfscQosHandler) // 设置 HFSC QoS
	rg.POST("/bridge/set-datapath-type", api.SetDatapathTypeHandler) // 设置 datapath_type
} 