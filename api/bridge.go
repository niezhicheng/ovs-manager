package api

import (
	"net/http"
	"ovs-manager/service"
	"github.com/gin-gonic/gin"
)

// ListBridgesHandler 交换机列表接口
// @Summary 获取所有 OVS 交换机
// @Description 获取所有 OVS 交换机列表，无需参数
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/list [post]
func ListBridgesHandler(c *gin.Context) {
	bridges, err := service.ListBridges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bridges": bridges})
}

// AddBridgeRequest 新增交换机请求结构体
// @Summary 新增 OVS 交换机
// @Description 新增一个 OVS 交换机
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Param data body AddBridgeRequest true "交换机名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/add [post]
type AddBridgeRequest struct {
	// 交换机名称
	Name string `json:"name" binding:"required"`
}
func AddBridgeHandler(c *gin.Context) {
	var req AddBridgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddBridge(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeleteBridgeRequest 删除交换机请求结构体
// @Summary 删除 OVS 交换机
// @Description 删除一个 OVS 交换机
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Param data body DeleteBridgeRequest true "交换机名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/delete [post]
type DeleteBridgeRequest struct {
	// 交换机名称
	Name string `json:"name" binding:"required"`
}
func DeleteBridgeHandler(c *gin.Context) {
	var req DeleteBridgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeleteBridge(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetNetFlowRequest 设置 NetFlow 请求结构体
// @Summary 设置 NetFlow
// @Description 设置网桥的 NetFlow 配置
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Param data body SetNetFlowRequest true "网桥、NetFlow 服务器、端口等"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/set-netflow [post]
type SetNetFlowRequest struct {
	Bridge   string `json:"bridge" binding:"required"`
	Target   string `json:"target" binding:"required"`
	EngineID int    `json:"engineID"`
}
func SetNetFlowHandler(c *gin.Context) {
	var req SetNetFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetNetFlow(req.Bridge, req.Target, req.EngineID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetSFlowRequest 设置 sFlow 请求结构体
// @Summary 设置 sFlow
// @Description 设置网桥的 sFlow 配置
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Param data body SetSFlowRequest true "网桥、sFlow 服务器、采样率等"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/set-sflow [post]
type SetSFlowRequest struct {
	Bridge     string   `json:"bridge" binding:"required"`
	Targets    []string `json:"targets" binding:"required"`
	Sampling   int      `json:"sampling"`
	Header     int      `json:"header"`
	Polling    int      `json:"polling"`
	Agent      string   `json:"agent"`
}
func SetSFlowHandler(c *gin.Context) {
	var req SetSFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetSFlow(req.Bridge, req.Targets, req.Sampling, req.Header, req.Polling, req.Agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetStpRequest 设置 STP 请求结构体
// @Summary 设置 STP
// @Description 设置网桥的 STP 配置
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Param data body SetStpRequest true "网桥、stp_enable"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/set-stp [post]
type SetStpRequest struct {
	Bridge   string `json:"bridge" binding:"required"`
	Enable   bool   `json:"enable" binding:"required"`
}
func SetStpHandler(c *gin.Context) {
	var req SetStpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetStp(req.Bridge, req.Enable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetQosRequest 设置 QoS 请求结构体
// @Summary 设置 QoS
// @Description 设置端口的 QoS 配置
// @Tags OVS-QoS
// @Accept json
// @Produce json
// @Param data body SetQosRequest true "端口、type、max-rate、queues等"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/port/set-qos [post]
type SetQosRequest struct {
	PortName string            `json:"portName" binding:"required"`
	Type     string            `json:"type" binding:"required"`
	MaxRate  string            `json:"maxRate"`
	Queues   map[string]string `json:"queues"`
}
func SetQosHandler(c *gin.Context) {
	var req SetQosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetQos(req.PortName, req.Type, req.MaxRate, req.Queues); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetRstpRequest 设置 RSTP 请求结构体
// @Summary 设置 RSTP
// @Description 设置网桥的 RSTP 配置
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Param data body SetRstpRequest true "网桥、rstp_enable"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/set-rstp [post]
type SetRstpRequest struct {
	Bridge   string `json:"bridge" binding:"required"`
	Enable   bool   `json:"enable" binding:"required"`
}
func SetRstpHandler(c *gin.Context) {
	var req SetRstpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetRstp(req.Bridge, req.Enable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetIpfixRequest 设置 IPFIX 请求结构体
// @Summary 设置 IPFIX
// @Description 设置网桥的 IPFIX 配置
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Param data body SetIpfixRequest true "网桥、IPFIX 服务器、采样率等"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/set-ipfix [post]
type SetIpfixRequest struct {
	Bridge   string   `json:"bridge" binding:"required"`
	Targets  []string `json:"targets" binding:"required"`
	Sampling int      `json:"sampling"`
	ObsDomainID int   `json:"obsDomainID"`
	ObsPointID  int   `json:"obsPointID"`
}
func SetIpfixHandler(c *gin.Context) {
	var req SetIpfixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetIpfix(req.Bridge, req.Targets, req.Sampling, req.ObsDomainID, req.ObsPointID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DumpFlowsRequest 查询流缓存请求结构体
// @Summary 查询流缓存
// @Description 查询网桥的流缓存（dump-flows）
// @Tags OVS-Bridge
// @Accept json
// @Produce json
// @Param data body DumpFlowsRequest true "网桥"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bridge/dump-flows [post]
type DumpFlowsRequest struct {
	Bridge string `json:"bridge" binding:"required"`
}
func DumpFlowsHandler(c *gin.Context) {
	var req DumpFlowsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := service.DumpFlows(req.Bridge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": output})
} 