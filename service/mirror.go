package service

import (
	"fmt"
	"os/exec"
)

// AddMirror 新增端口镜像
func AddMirror(bridge, name string, selectSrcPorts, selectDstPorts []string, selectVlan *int, outputPort string, outputVlan *int, selectAll bool) error {
	args := []string{"--", "--id=@m", "create", "Mirror", fmt.Sprintf("name=%s", name)}
	if selectAll {
		args = append(args, "select_all=true")
	}
	for _, p := range selectSrcPorts {
		args = append(args, fmt.Sprintf("select-src-port=%s", p))
	}
	for _, p := range selectDstPorts {
		args = append(args, fmt.Sprintf("select-dst-port=%s", p))
	}
	if selectVlan != nil {
		args = append(args, fmt.Sprintf("select_vlan=%d", *selectVlan))
	}
	if outputPort != "" {
		args = append(args, fmt.Sprintf("output-port=%s", outputPort))
	}
	if outputVlan != nil {
		args = append(args, fmt.Sprintf("output_vlan=%d", *outputVlan))
	}
	args = append([]string{"--", "set", "Bridge", bridge, "mirrors=@m"}, args...)
	cmd := exec.Command("ovs-vsctl", args...)
	return cmd.Run()
}

// DeleteMirror 删除端口镜像
func DeleteMirror(bridge, name string) error {
	cmd := exec.Command("ovs-vsctl", "--", "clear", "Bridge", bridge, "mirrors")
	return cmd.Run()
}

// ListMirrors 查询端口镜像
func ListMirrors(bridge string) (string, error) {
	// 查询网桥的镜像配置
	cmd := exec.Command("ovs-vsctl", "get", "Bridge", bridge, "mirrors")
	output, err := cmd.Output()
	if err != nil {
		// 如果网桥不存在或没有镜像配置，返回空字符串
		if err.Error() == "exit status 1" {
			return "", nil
		}
		return "", err
	}
	
	// 如果没有镜像配置，返回空字符串
	outputStr := string(output)
	if outputStr == "[]" || outputStr == "" || outputStr == "\n" || outputStr == "[]\n" {
		return "", nil
	}
	
	// 如果有镜像配置，查询镜像的详细信息
	// 使用find命令查找所有镜像
	cmd = exec.Command("ovs-vsctl", "find", "Mirror")
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
} 