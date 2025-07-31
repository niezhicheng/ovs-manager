package api

import (
	"net/http"
	"ovs-manager/service"
	"github.com/gin-gonic/gin"
)

// AddBondRequest 新增 Bond 请求结构体
// @Summary 新增 Bond 端口
// @Description 新增 Bond 端口并可设置 bond_mode、lacp、其它参数
// @Tags OVS-Bond
// @Accept json
// @Produce json
// @Param data body AddBondRequest true "Bond 参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bond/add [post]
type AddBondRequest struct {
	Bridge      string            `json:"bridge" binding:"required"`
	BondName    string            `json:"bondName" binding:"required"`
	Slaves      []string          `json:"slaves" binding:"required"`
	BondMode    string            `json:"bondMode"`
	Lacp        string            `json:"lacp"`
	OtherOptions map[string]string `json:"otherOptions"`
}
func AddBondHandler(c *gin.Context) {
	var req AddBondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddBond(req.Bridge, req.BondName, req.Slaves, req.BondMode, req.Lacp, req.OtherOptions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetBondRequest 设置 Bond 属性请求结构体
// @Summary 设置 Bond 属性
// @Description 设置 Bond 端口的 bond_mode、lacp、其它参数
// @Tags OVS-Bond
// @Accept json
// @Produce json
// @Param data body SetBondRequest true "Bond 参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bond/set [post]
type SetBondRequest struct {
	BondName    string            `json:"bondName" binding:"required"`
	BondMode    string            `json:"bondMode"`
	Lacp        string            `json:"lacp"`
	OtherOptions map[string]string `json:"otherOptions"`
}
func SetBondHandler(c *gin.Context) {
	var req SetBondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.SetBond(req.BondName, req.BondMode, req.Lacp, req.OtherOptions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// ShowBondRequest 查询 Bond 状态请求结构体
// @Summary 查询 Bond 状态
// @Description 查询 Bond 端口的详细状态
// @Tags OVS-Bond
// @Accept json
// @Produce json
// @Param data body ShowBondRequest true "Bond 名称"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bond/show [post]
type ShowBondRequest struct {
	BondName string `json:"bondName" binding:"required"`
}
func ShowBondHandler(c *gin.Context) {
	var req ShowBondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bondShow, lacpShow, portInfo, err := service.ShowBond(req.BondName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bondShow": bondShow, "lacpShow": lacpShow, "portInfo": portInfo})
}

// ListBondsRequest 查询所有 Bond 端口请求结构体
// @Summary 查询所有 Bond 端口
// @Description 查询所有 Bond 端口及其成员和模式
// @Tags OVS-Bond
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bond/show [post]
func ListBondsHandler(c *gin.Context) {
	bonds, err := service.ListBonds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bonds": bonds})
}

// DeleteBondRequest 删除 Bond 请求结构体
// @Summary 删除 Bond 端口
// @Description 删除 Bond 端口
// @Tags OVS-Bond
// @Accept json
// @Produce json
// @Param data body DeleteBondRequest true "Bond 参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/bond/delete [post]
type DeleteBondRequest struct {
	Bridge   string `json:"bridge" binding:"required"`
	BondName string `json:"bondName" binding:"required"`
}
func DeleteBondHandler(c *gin.Context) {
	var req DeleteBondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeleteBond(req.Bridge, req.BondName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
} 