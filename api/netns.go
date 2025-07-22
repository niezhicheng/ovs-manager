package api

import (
	"net/http"
	"ovs-manager/service"
	"github.com/gin-gonic/gin"
)

// CreateNetnsRequest 创建命名空间请求结构体
// @Summary 创建网络命名空间
// @Description 创建一个新的网络命名空间
// @Tags Netns
// @Accept json
// @Produce json
// @Param data body CreateNetnsRequest true "命名空间名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/netns/create [post]
type CreateNetnsRequest struct {
	Name string `json:"name" binding:"required"`
}
func CreateNetnsHandler(c *gin.Context) {
	var req CreateNetnsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.CreateNetns(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeleteNetnsRequest 删除命名空间请求结构体
// @Summary 删除网络命名空间
// @Description 删除一个网络命名空间
// @Tags Netns
// @Accept json
// @Produce json
// @Param data body DeleteNetnsRequest true "命名空间名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/netns/delete [post]
type DeleteNetnsRequest struct {
	Name string `json:"name" binding:"required"`
}
func DeleteNetnsHandler(c *gin.Context) {
	var req DeleteNetnsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeleteNetns(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// ListNetnsHandler 命名空间列表接口
// @Summary 获取所有网络命名空间
// @Description 获取所有网络命名空间列表
// @Tags Netns
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/netns/list [post]
func ListNetnsHandler(c *gin.Context) {
	netns, err := service.ListNetns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"netns": netns})
} 