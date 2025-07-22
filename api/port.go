package api

import (
	"net/http"
	"ovs-manager/service"
	"github.com/gin-gonic/gin"
)

// ListPortsHandler 端口列表接口
// @Summary 获取指定交换机的所有端口
// @Description 获取指定 OVS 交换机下的所有端口
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body ListPortsRequest true "交换机名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/list [post]
type ListPortsRequest struct {
	Bridge string `json:"bridge" binding:"required"`
}
func ListPortsHandler(c *gin.Context) {
	var req ListPortsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ports, err := service.ListPorts(req.Bridge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ports": ports})
}

// AddPortRequest 新增端口请求结构体
// @Summary 新增端口
// @Description 向指定 OVS 交换机添加端口，可指定端口类型（如 internal）
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body AddPortRequest true "交换机名称、端口名称、端口类型"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/add [post]
type AddPortRequest struct {
	Bridge string `json:"bridge" binding:"required"`
	PortName string `json:"portName" binding:"required"`
	Type string `json:"type"`
}
func AddPortHandler(c *gin.Context) {
	var req AddPortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddPort(req.Bridge, req.PortName, req.Type); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeletePortRequest 删除端口请求结构体
// @Summary 删除端口
// @Description 从指定 OVS 交换机删除端口
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body DeletePortRequest true "交换机名称和端口名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/delete [post]
type DeletePortRequest struct {
	Bridge string `json:"bridge" binding:"required"`
	PortName string `json:"portName" binding:"required"`
}
func DeletePortHandler(c *gin.Context) {
	var req DeletePortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeletePort(req.Bridge, req.PortName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// BindPortToNetnsRequest 端口绑定命名空间请求结构体
// @Summary 端口绑定命名空间
// @Description 将端口绑定到指定网络命名空间
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body BindPortToNetnsRequest true "端口名称和命名空间名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/bind-netns [post]
type BindPortToNetnsRequest struct {
	PortName string `json:"portName" binding:"required"`
	Netns string `json:"netns" binding:"required"`
}
func BindPortToNetnsHandler(c *gin.Context) {
	var req BindPortToNetnsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.BindPortToNetns(req.PortName, req.Netns); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// UnbindPortFromNetnsRequest 端口解绑命名空间请求结构体
// @Summary 端口解绑到主命名空间
// @Description 将端口解绑到主命名空间
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body UnbindPortFromNetnsRequest true "端口名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/unbind-netns [post]
type UnbindPortFromNetnsRequest struct {
	PortName string `json:"portName" binding:"required"`
}
func UnbindPortFromNetnsHandler(c *gin.Context) {
	var req UnbindPortFromNetnsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.UnbindPortFromNetns(req.PortName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetPortUpDownRequest 端口 up/down 请求结构体
// @Summary 设置端口 up/down
// @Description 设置指定命名空间下端口 up/down
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body SetPortUpDownRequest true "命名空间、端口名称、up/down"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/updown [post]
type SetPortUpDownRequest struct {
	Netns string `json:"netns" binding:"required"`
	PortName string `json:"portName" binding:"required"`
	Up bool `json:"up"`
}
func SetPortUpDownHandler(c *gin.Context) {
	var req SetPortUpDownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetPortUpDown(req.Netns, req.PortName, req.Up); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetPortAddrRequest 端口分配 IP 请求结构体
// @Summary 端口分配 IP
// @Description 给指定命名空间下端口分配 IP 地址
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body SetPortAddrRequest true "命名空间、端口名称、IP地址"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/addr [post]
type SetPortAddrRequest struct {
	Netns string `json:"netns" binding:"required"`
	PortName string `json:"portName" binding:"required"`
	IP string `json:"ip" binding:"required"`
}
func SetPortAddrHandler(c *gin.Context) {
	var req SetPortAddrRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetPortAddr(req.Netns, req.PortName, req.IP); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetPortVlanTagRequest 端口 VLAN tag 设置请求结构体
// @Summary 设置端口 VLAN tag
// @Description 设置端口 VLAN tag
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body SetPortVlanTagRequest true "端口名称、VLAN tag"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/set-vlan [post]
type SetPortVlanTagRequest struct {
	PortName string `json:"portName" binding:"required"`
	Tag int `json:"tag" binding:"required"`
}
func SetPortVlanTagHandler(c *gin.Context) {
	var req SetPortVlanTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetPortVlanTag(req.PortName, req.Tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetPortVlanModeRequest 设置端口 VLAN mode 请求结构体
// @Summary 设置端口 VLAN mode
// @Description 设置端口 VLAN mode（trunk/access/native-tagged/native-untagged）
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body SetPortVlanModeRequest true "端口名称、VLAN mode"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/set-vlan-mode [post]
type SetPortVlanModeRequest struct {
	PortName string `json:"portName" binding:"required"`
	VlanMode string `json:"vlanMode" binding:"required"`
}
func SetPortVlanModeHandler(c *gin.Context) {
	var req SetPortVlanModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetPortVlanMode(req.PortName, req.VlanMode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetPortTrunksRequest 设置端口 trunks 请求结构体
// @Summary 设置端口 trunks
// @Description 设置端口允许通过的 VLAN trunk 列表
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body SetPortTrunksRequest true "端口名称、trunks"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/set-trunks [post]
type SetPortTrunksRequest struct {
	PortName string `json:"portName" binding:"required"`
	Trunks   []int  `json:"trunks" binding:"required"`
}
func SetPortTrunksHandler(c *gin.Context) {
	var req SetPortTrunksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetPortTrunks(req.PortName, req.Trunks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// RemovePortPropertyRequest 移除端口属性请求结构体
// @Summary 移除端口属性
// @Description 移除端口的某个属性（如 tag、trunks）
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body RemovePortPropertyRequest true "端口名称、属性、值"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/remove-property [post]
type RemovePortPropertyRequest struct {
	PortName string      `json:"portName" binding:"required"`
	Property string      `json:"property" binding:"required"`
	Value    interface{} `json:"value" binding:"required"`
}
func RemovePortPropertyHandler(c *gin.Context) {
	var req RemovePortPropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.RemovePortProperty(req.PortName, req.Property, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// AddPatchPortRequest 添加 Patch Port 请求结构体
// @Summary 添加 Patch Port
// @Description 添加 Patch Port，支持 peer 配置
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body AddPatchPortRequest true "网桥、端口、peer"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/patch/add [post]
type AddPatchPortRequest struct {
	Bridge   string `json:"bridge" binding:"required"`
	PortName string `json:"portName" binding:"required"`
	Peer     string `json:"peer" binding:"required"`
}
func AddPatchPortHandler(c *gin.Context) {
	var req AddPatchPortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddPatchPort(req.Bridge, req.PortName, req.Peer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// AddTunnelPortRequest 添加 Tunnel Port 请求结构体
// @Summary 添加 Tunnel Port（GRE/Geneve）
// @Description 添加 GRE/Geneve Tunnel Port，支持 options 配置
// @Tags OVS-Tunnel
// @Accept json
// @Produce json
// @Param data body AddTunnelPortRequest true "网桥、端口、类型、options"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/tunnel/add [post]
type AddTunnelPortRequest struct {
	Bridge   string            `json:"bridge" binding:"required"`
	PortName string            `json:"portName" binding:"required"`
	Type     string            `json:"type" binding:"required"`
	Options  map[string]string `json:"options"`
}
func AddTunnelPortHandler(c *gin.Context) {
	var req AddTunnelPortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddTunnelPort(req.Bridge, req.PortName, req.Type, req.Options); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// PortInfoRequest 查询端口属性请求结构体
// @Summary 查询端口属性
// @Description 查询端口/interface 详细属性
// @Tags OVS-Port
// @Accept json
// @Produce json
// @Param data body PortInfoRequest true "端口名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/info [post]
type PortInfoRequest struct {
	PortName string `json:"portName" binding:"required"`
}
func PortInfoHandler(c *gin.Context) {
	var req PortInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	info, err := service.PortInfo(req.PortName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"info": info})
}

// SetBfdRequest 设置 BFD 请求结构体
// @Summary 设置 BFD
// @Description 设置端口的 BFD 配置
// @Tags OVS-BFD
// @Accept json
// @Produce json
// @Param data body SetBfdRequest true "端口名称、BFD 配置"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/set-bfd [post]
type SetBfdRequest struct {
	PortName string            `json:"portName" binding:"required"`
	Bfd      map[string]string `json:"bfd" binding:"required"`
}
func SetBfdHandler(c *gin.Context) {
	var req SetBfdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetBfd(req.PortName, req.Bfd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetCfmRequest 设置 CFM 请求结构体
// @Summary 设置 CFM
// @Description 设置端口的 CFM (802.1ag) 配置
// @Tags OVS-CFM
// @Accept json
// @Produce json
// @Param data body SetCfmRequest true "端口名称、CFM 配置"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/set-cfm [post]
type SetCfmRequest struct {
	PortName string            `json:"portName" binding:"required"`
	Cfm      map[string]string `json:"cfm" binding:"required"`
}
func SetCfmHandler(c *gin.Context) {
	var req SetCfmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetCfm(req.PortName, req.Cfm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetMcastSnoopingRequest 设置组播监听请求结构体
// @Summary 设置组播监听
// @Description 设置网桥的组播监听开关
// @Tags OVS-Multicast
// @Accept json
// @Produce json
// @Param data body SetMcastSnoopingRequest true "网桥、enable"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/set-mcast-snooping [post]
type SetMcastSnoopingRequest struct {
	Bridge string `json:"bridge" binding:"required"`
	Enable bool   `json:"enable" binding:"required"`
}
func SetMcastSnoopingHandler(c *gin.Context) {
	var req SetMcastSnoopingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetMcastSnooping(req.Bridge, req.Enable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetHfscQosRequest 设置 HFSC QoS 请求结构体
// @Summary 设置 HFSC QoS
// @Description 设置端口的 HFSC QoS 配置
// @Tags OVS-QoS
// @Accept json
// @Produce json
// @Param data body SetHfscQosRequest true "端口、max-rate、queues等"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/set-hfsc-qos [post]
type SetHfscQosRequest struct {
	PortName string            `json:"portName" binding:"required"`
	MaxRate  string            `json:"maxRate"`
	Queues   map[string]string `json:"queues"`
}
func SetHfscQosHandler(c *gin.Context) {
	var req SetHfscQosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetHfscQos(req.PortName, req.MaxRate, req.Queues); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetDatapathTypeRequest 设置 datapath 类型请求结构体
// @Summary 设置 datapath 类型
// @Description 设置网桥的 datapath_type（如 system、netdev）
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Param data body SetDatapathTypeRequest true "网桥、datapath_type"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/set-datapath-type [post]
type SetDatapathTypeRequest struct {
	Bridge        string `json:"bridge" binding:"required"`
	DatapathType  string `json:"datapathType" binding:"required"`
}
func SetDatapathTypeHandler(c *gin.Context) {
	var req SetDatapathTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetDatapathType(req.Bridge, req.DatapathType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
} 