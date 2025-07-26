package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ovs-manager/service"
)

// AddVxlanPortRequest 新增 VXLAN 端口请求结构体
// @Summary 新增 VXLAN 端口
// @Description 新增 VXLAN 端口，可指定 remote_ip、vni、key、local_ip
// @Tags OVS-VXLAN
// @Accept json
// @Produce json
// @Param data body AddVxlanPortRequest true "VXLAN 端口参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/vxlan/add [post]
type AddVxlanPortCustomRequest struct {
	Bridge   string `json:"bridge" binding:"required"`
	PortName string `json:"portName" binding:"required"`
	RemoteIP string `json:"remoteIP" binding:"required"`
	VNI      int    `json:"vni" binding:"required"`
	Key      string `json:"key"`
	LocalIP  string `json:"localIP"`
}

func AddVxlanPortCustomHandler(c *gin.Context) {
	var req AddVxlanPortCustomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddVxlanPortCustom(req.Bridge, req.PortName, req.RemoteIP, req.VNI, req.Key, req.LocalIP); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeleteVxlanPortRequest 删除 VXLAN 端口请求结构体
// @Summary 删除 VXLAN 端口
// @Description 删除 VXLAN 端口
// @Tags OVS-VXLAN
// @Accept json
// @Produce json
// @Param data body DeleteVxlanPortRequest true "VXLAN 端口参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/vxlan/delete [post]
type DeleteVxlanPortRequest struct {
	Bridge   string `json:"bridge" binding:"required"`
	PortName string `json:"portName" binding:"required"`
}

func DeleteVxlanPortHandler(c *gin.Context) {
	var req DeleteVxlanPortRequest
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
