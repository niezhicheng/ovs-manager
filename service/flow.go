package service

import (
	"os/exec"
)

// ListFlowsV2 查询指定 bridge 的所有流表（支持自定义表达式）
func ListFlowsV2(bridge string) (string, error) {
	cmd := exec.Command("ovs-ofctl", "dump-flows", bridge)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// AddFlowV2 添加流表规则
func AddFlowV2(bridge, flow string) error {
	cmd := exec.Command("ovs-ofctl", "add-flow", bridge, flow)
	return cmd.Run()
}

// DeleteFlowV2 删除流表规则（支持全删和条件删）
func DeleteFlowV2(bridge, match string) error {
	if match == "" {
		cmd := exec.Command("ovs-ofctl", "del-flows", bridge)
		return cmd.Run()
	}
	cmd := exec.Command("ovs-ofctl", "del-flows", bridge, match)
	return cmd.Run()
} 