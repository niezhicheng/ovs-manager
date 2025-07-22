package api

import (
	"net/http"
	"ovs-manager/service"
	"github.com/gin-gonic/gin"
)

// ScenarioStep 场景步骤
type ScenarioStep struct {
	Action string                 `json:"action" binding:"required"`
	Params map[string]interface{} `json:"params" binding:"required"`
}

// ScenarioApplyRequest 场景引导请求体
// @Summary 场景引导式一键操作
// @Description 按步骤批量执行 OVS 操作，支持自定义 steps 或内置模板
// @Tags OVS-Scenario
// @Accept json
// @Produce json
// @Param data body ScenarioApplyRequest true "场景类型、步骤、参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/ovs/scenario/apply [post]
type ScenarioApplyRequest struct {
	Scenario string         `json:"scenario"` // 可选，内置模板名
	Steps    []ScenarioStep `json:"steps"`    // 自定义步骤
}

type ScenarioStepResult struct {
	Action  string      `json:"action"`
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Output  interface{} `json:"output,omitempty"`
}

type ScenarioApplyResponse struct {
	Success bool                 `json:"success"`
	Results []ScenarioStepResult `json:"results"`
}

// 内置场景模板
var scenarioTemplates = map[string][]ScenarioStep{
	"vxlan_vlan_isolation": {
		{Action: "add_bridge", Params: map[string]interface{}{ "name": "br-int" }},
		{Action: "add_port", Params: map[string]interface{}{ "bridge": "br-int", "portName": "vnet0", "type": "internal" }},
		{Action: "set_port_vlan", Params: map[string]interface{}{ "portName": "vnet0", "tag": 100 }},
		{Action: "add_bond", Params: map[string]interface{}{ "bridge": "br-int", "bondName": "bond0", "slaves": []interface{}{ "eth0", "eth1" }, "bondMode": "balance-tcp" }},
	},
	"patch_trunk": {
		{Action: "add_bridge", Params: map[string]interface{}{ "name": "br0" }},
		{Action: "add_bridge", Params: map[string]interface{}{ "name": "br1" }},
		{Action: "add_patch_port", Params: map[string]interface{}{ "bridge": "br0", "portName": "patch0", "peer": "patch1" }},
		{Action: "add_patch_port", Params: map[string]interface{}{ "bridge": "br1", "portName": "patch1", "peer": "patch0" }},
		{Action: "set_port_vlan_mode", Params: map[string]interface{}{ "portName": "patch0", "vlanMode": "trunk" }},
	},
}

// ScenarioApplyHandler 场景引导接口
func ScenarioApplyHandler(c *gin.Context) {
	var req ScenarioApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	steps := req.Steps
	if len(steps) == 0 && req.Scenario != "" {
		tpl, ok := scenarioTemplates[req.Scenario]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown scenario template"})
			return
		}
		steps = tpl
	}
	results := make([]ScenarioStepResult, 0, len(steps))
	success := true
	for _, step := range steps {
		res := ScenarioStepResult{Action: step.Action}
		err, output := service.ExecuteScenarioStep(step.Action, step.Params)
		if err != nil {
			res.Success = false
			res.Error = err.Error()
			success = false
		} else {
			res.Success = true
			if output != nil {
				res.Output = output
			}
		}
		results = append(results, res)
	}
	c.JSON(http.StatusOK, ScenarioApplyResponse{Success: success, Results: results})
} 