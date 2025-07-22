package api

import (
	"net/http"
	"ovs-manager/service"
	"github.com/gin-gonic/gin"
)

// AddMirrorHandler 新增端口镜像接口
// @Summary 新增端口镜像
// @Description 新增端口镜像
// @Tags OVS-Mirror
// @Accept json
// @Produce json
// @Param data body AddMirrorRequest true "端口镜像参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/mirror/add [post]
type AddMirrorRequest struct {
	Bridge string `json:"bridge" binding:"required"`
	Name string `json:"name" binding:"required"`
	SelectSrcPorts []string `json:"selectSrcPorts"`
	SelectDstPorts []string `json:"selectDstPorts"`
	SelectVlan *int `json:"selectVlan"`
	OutputPort string `json:"outputPort"`
	OutputVlan *int `json:"outputVlan"`
	SelectAll bool `json:"selectAll"`
}
func AddMirrorHandler(c *gin.Context) {
	var req AddMirrorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddMirror(req.Bridge, req.Name, req.SelectSrcPorts, req.SelectDstPorts, req.SelectVlan, req.OutputPort, req.OutputVlan, req.SelectAll); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeleteMirrorRequest 删除端口镜像请求结构体
// @Summary 删除端口镜像
// @Description 删除端口镜像
// @Tags OVS-Mirror
// @Accept json
// @Produce json
// @Param data body DeleteMirrorRequest true "端口镜像参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/mirror/delete [post]
type DeleteMirrorRequest struct {
	Bridge string `json:"bridge" binding:"required"`
	Name string `json:"name" binding:"required"`
}
func DeleteMirrorHandler(c *gin.Context) {
	var req DeleteMirrorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeleteMirror(req.Bridge, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// ListMirrorsRequest 查询端口镜像请求结构体
// @Summary 查询端口镜像
// @Description 查询端口镜像
// @Tags OVS-Mirror
// @Accept json
// @Produce json
// @Param data body ListMirrorsRequest true "端口镜像参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/mirror/list [post]
type ListMirrorsRequest struct {
	Bridge string `json:"bridge" binding:"required"`
}
func ListMirrorsHandler(c *gin.Context) {
	var req ListMirrorsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := service.ListMirrors(req.Bridge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": output})
} 