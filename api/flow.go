package api

import (
	"net/http"
	"ovs-manager/service"
	"github.com/gin-gonic/gin"
)

// ListFlowsV2Handler 查询流表规则接口
// @Summary 查询流表规则
// @Description 查询指定 bridge 的所有流表规则
// @Tags OVS-Flow
// @Accept json
// @Produce json
// @Param data body ListFlowsV2Request true "网桥名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/flow/list-v2 [post]
type ListFlowsV2Request struct {
	Bridge string `json:"bridge" binding:"required"`
}
func ListFlowsV2Handler(c *gin.Context) {
	var req ListFlowsV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := service.ListFlowsV2(req.Bridge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": output})
}

// AddFlowV2Request 添加流表规则请求结构体
// @Summary 添加流表规则
// @Description 添加流表规则，flow 字符串为完整表达式
// @Tags OVS-Flow
// @Accept json
// @Produce json
// @Param data body AddFlowV2Request true "网桥名称和流表表达式"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/flow/add-v2 [post]
type AddFlowV2Request struct {
	Bridge string `json:"bridge" binding:"required"`
	Flow string `json:"flow" binding:"required"`
}
func AddFlowV2Handler(c *gin.Context) {
	var req AddFlowV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddFlowV2(req.Bridge, req.Flow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeleteFlowV2Request 删除流表规则请求结构体
// @Summary 删除流表规则
// @Description 删除流表规则，支持全删和条件删
// @Tags OVS-Flow
// @Accept json
// @Produce json
// @Param data body DeleteFlowV2Request true "网桥名称和匹配条件"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/flow/delete-v2 [post]
type DeleteFlowV2Request struct {
	Bridge string `json:"bridge" binding:"required"`
	Match string `json:"match"`
}
func DeleteFlowV2Handler(c *gin.Context) {
	var req DeleteFlowV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeleteFlowV2(req.Bridge, req.Match); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
} 