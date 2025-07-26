package router

import (
	"github.com/gin-gonic/gin"
	"ovs-manager/api"
)

// RegisterPortRoutes 注册 OVS 端口管理相关路由
func RegisterPortRoutes(rg *gin.RouterGroup) {
	rg.POST("/port/list", api.ListPortsHandler)        // 获取端口列表
	rg.POST("/port/add", api.AddPortHandler)           // 新增端口（通用）
	rg.POST("/port/delete", api.DeletePortHandler)     // 删除端口

	// 端口类型专用API
	rg.POST("/port/add-normal", api.AddNormalPortHandler)     // 新增普通端口
	rg.POST("/port/add-internal", api.AddInternalPortHandler) // 新增内部端口
	rg.POST("/port/add-gre", api.AddGrePortHandler)           // 新增GRE隧道端口
	rg.POST("/port/add-vxlan", api.AddVxlanPortHandler)       // 新增VXLAN隧道端口
	rg.POST("/port/add-bond", api.AddBondPortHandler)         // 新增Bond端口

	// 端口与命名空间绑定/解绑
	rg.POST("/port/bind-netns", api.BindPortToNetnsHandler)     // 端口绑定命名空间
	rg.POST("/port/unbind-netns", api.UnbindPortFromNetnsHandler) // 端口解绑到主命名空间

	// 端口 up/down、分配 IP
	rg.POST("/port/updown", api.SetPortUpDownHandler)           // 端口 up/down
	rg.POST("/port/addr", api.SetPortAddrHandler)               // 端口分配 IP
	rg.POST("/port/get-addrs", api.GetPortAddrsHandler)         // 获取端口IP地址列表
	rg.POST("/port/delete-addr", api.DeletePortAddrHandler)     // 删除端口IP地址

	// 端口 VLAN tag 设置
	rg.POST("/port/set-vlan", api.SetPortVlanTagHandler)         // 端口 VLAN tag 设置
	rg.POST("/tunnel/add", api.AddTunnelPortHandler) // 添加 GRE/Geneve Tunnel Port
	rg.POST("/set-bfd", api.SetBfdHandler) // 设置 BFD
	rg.POST("/set-cfm", api.SetCfmHandler) // 设置 CFM
	rg.POST("/bridge/set-mcast-snooping", api.SetMcastSnoopingHandler) // 设置组播监听
	rg.POST("/port/set-hfsc-qos", api.SetHfscQosHandler) // 设置 HFSC QoS
	rg.POST("/bridge/set-datapath-type", api.SetDatapathTypeHandler) // 设置 datapath_type
	rg.POST("/patch/add", api.AddPatchPortHandler) // 添加 patch 端口（peer可选）
	rg.POST("/patch/add-without-peer", api.AddPatchPortWithoutPeerHandler) // 添加不设置对端的 patch 端口
	rg.POST("/patch/set-peer", api.SetPatchPortPeerHandler) // 设置 patch 端口对端
	rg.POST("/patch/add-pair", api.AddPatchPortPairHandler) // 一键成对创建 patch 端口
	rg.POST("/port/patch-list", api.ListAllPatchPortsHandler) // 全局 patch 端口列表
	rg.POST("/tap/add", api.AddTapPortHandler) // 添加 tap 端口
	rg.POST("/tun/add", api.AddTunPortHandler) // 添加 tun 端口
	rg.POST("/port/set-type-peer", api.SetPortTypePeerHandler) // 设置端口类型和 peer
	rg.POST("/port/set-alias", api.SetPortAliasHandler) // 设置端口别名
	rg.POST("/port/set-route", api.SetPortRouteHandler) // 设置端口路由
	rg.POST("/port/delete-route", api.DeletePortRouteHandler) // 删除端口路由
	rg.POST("/port/get-routes", api.GetPortRoutesHandler) // 获取端口路由列表
} 