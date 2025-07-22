package api

import (
	"net/http"
	"ovs-manager/service"
	"github.com/gin-gonic/gin"
)

// ScenarioStep 表示场景中的单个操作步骤。
// Action: 操作类型（如 add_bridge、add_port、set_port_vlan 等）
// Params: 该操作所需的参数，key-value 形式，具体内容取决于 action
// 例如：{Action: "add_bridge", Params: {"name": "br0"}}
type ScenarioStep struct {
	Action string                 `json:"action" binding:"required"` // 步骤类型
	Params map[string]interface{} `json:"params" binding:"required"` // 步骤参数
}

// ScenarioApplyRequest 场景引导请求体
// Scenario: 可选，内置场景模板名（如 "vxlan_vlan_isolation"）
// Steps: 自定义步骤数组，若传递则优先生效
// Params: 可选，模板参数覆盖（如 {"bridge": "br-demo", "tag": 200}），会自动合并到模板 steps 的 Params 字段
// 只传 scenario 时，后端会自动填充对应模板步骤
// 只传 steps 时，按 steps 顺序执行
// 两者都不传则报错
type ScenarioApplyRequest struct {
	Scenario string                 `json:"scenario"` // 场景模板名（可选）
	Steps    []ScenarioStep         `json:"steps"`    // 自定义步骤（可选）
	Params   map[string]interface{} `json:"params"`   // 覆盖模板参数（可选）
}

// ScenarioStepResult 表示单个步骤的执行结果
// Action: 步骤类型
// Success: 是否成功
// Error: 错误信息（如有）
// Output: 额外输出（如有）
type ScenarioStepResult struct {
	Action  string      `json:"action"`
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Output  interface{} `json:"output,omitempty"`
}

// ScenarioApplyResponse 场景引导接口的返回体
// Success: 所有步骤是否全部成功
// Results: 每一步的详细结果
type ScenarioApplyResponse struct {
	Success bool                 `json:"success"`
	Results []ScenarioStepResult `json:"results"`
}

// 内置场景模板，常用场景一键化
// key: 模板名，value: 步骤数组
// 可扩展更多模板，如 vxlan_vlan_isolation、patch_trunk 等
var scenarioTemplates = map[string][]ScenarioStep{
	"vxlan_vlan_isolation": {
		// 创建网桥
		{Action: "add_bridge", Params: map[string]interface{}{ "name": "br-int" }},
		// 添加 internal 端口
		{Action: "add_port", Params: map[string]interface{}{ "bridge": "br-int", "portName": "vnet0", "type": "internal" }},
		// 设置端口 VLAN tag
		{Action: "set_port_vlan", Params: map[string]interface{}{ "portName": "vnet0", "tag": 100 }},
		// 添加 bond
		{Action: "add_bond", Params: map[string]interface{}{ "bridge": "br-int", "bondName": "bond0", "slaves": []interface{}{ "eth0", "eth1" }, "bondMode": "balance-tcp" }},
	},
	"patch_trunk": {
		// 创建两个网桥
		{Action: "add_bridge", Params: map[string]interface{}{ "name": "br0" }},
		{Action: "add_bridge", Params: map[string]interface{}{ "name": "br1" }},
		// 添加 patch 端口并互为 peer
		{Action: "add_patch_port", Params: map[string]interface{}{ "bridge": "br0", "portName": "patch0", "peer": "patch1" }},
		{Action: "add_patch_port", Params: map[string]interface{}{ "bridge": "br1", "portName": "patch1", "peer": "patch0" }},
		// 设置 patch0 为 trunk 模式
		{Action: "set_port_vlan_mode", Params: map[string]interface{}{ "portName": "patch0", "vlanMode": "trunk" }},
	},
}

// mergeParams 合并模板步骤参数和用户传入的 params，用户参数优先
func mergeParams(stepParams, userParams map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(stepParams)+len(userParams))
	for k, v := range stepParams {
		merged[k] = v
	}
	for k, v := range userParams {
		merged[k] = v
	}
	return merged
}

// ScenarioApplyHandler 场景引导接口
// 支持三种用法：
// 1. 传 scenario，自动按模板执行一组步骤（可用 params 覆盖默认参数）
// 2. 传 steps，自定义步骤顺序和参数
// 3. 传 scenario+params，模板结构+自定义参数，兼顾易用和灵活
// 返回每一步的 success/error/output，便于前端引导和展示
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
		// 合并 params 到每个模板步骤
		steps = make([]ScenarioStep, len(tpl))
		for i, s := range tpl {
			steps[i] = ScenarioStep{
				Action: s.Action,
				Params: mergeParams(s.Params, req.Params),
			}
		}
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