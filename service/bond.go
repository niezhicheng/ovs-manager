package service

import (
	"fmt"
	"os/exec"
	"strings"
)

// AddBond 新增 Bond 端口并设置属性
func AddBond(bridge, bondName string, slaves []string, bondMode, lacp string, otherOptions map[string]string) error {
	args := []string{"add-bond", bridge, bondName}
	args = append(args, slaves...)
	cmd := exec.Command("ovs-vsctl", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	setArgs := []string{"set", "port", bondName}
	if bondMode != "" {
		setArgs = append(setArgs, fmt.Sprintf("bond_mode=%s", bondMode))
	}
	if lacp != "" {
		setArgs = append(setArgs, fmt.Sprintf("lacp=%s", lacp))
	}
	for k, v := range otherOptions {
		setArgs = append(setArgs, fmt.Sprintf("%s=%s", k, v))
	}
	if len(setArgs) > 3 {
		cmd2 := exec.Command("ovs-vsctl", setArgs...)
		if err := cmd2.Run(); err != nil {
			return err
		}
	}
	return nil
}

// SetBond 设置 Bond 端口属性
func SetBond(bondName, bondMode, lacp string, otherOptions map[string]string) error {
	setArgs := []string{"set", "port", bondName}
	if bondMode != "" {
		setArgs = append(setArgs, fmt.Sprintf("bond_mode=%s", bondMode))
	}
	if lacp != "" {
		setArgs = append(setArgs, fmt.Sprintf("lacp=%s", lacp))
	}
	for k, v := range otherOptions {
		setArgs = append(setArgs, fmt.Sprintf("%s=%s", k, v))
	}
	if len(setArgs) > 3 {
		cmd := exec.Command("ovs-vsctl", setArgs...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// ShowBond 查询 Bond 详细状态
func ShowBond(bondName string) (string, string, string, error) {
	bondShow, err := exec.Command("ovs-appctl", "bond/show", bondName).CombinedOutput()
	if err != nil {
		return "", "", "", err
	}
	lacpShow, err := exec.Command("ovs-appctl", "lacp/show", bondName).CombinedOutput()
	if err != nil {
		return string(bondShow), "", "", err
	}
	portInfo, err := exec.Command("ovs-vsctl", "list", "port", bondName).CombinedOutput()
	if err != nil {
		return string(bondShow), string(lacpShow), "", err
	}
	return string(bondShow), string(lacpShow), string(portInfo), nil
}

// DeleteBond 删除 Bond 端口
func DeleteBond(bridge, bondName string) error {
	cmd := exec.Command("ovs-vsctl", "del-port", bridge, bondName)
	return cmd.Run()
}

// ListBonds 返回所有 Bond 端口及其成员和模式
func ListBonds() ([]map[string]interface{}, error) {
	// 获取所有 port
	out, err := exec.Command("ovs-vsctl", "list", "port").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var bonds []map[string]interface{}
	var bond map[string]interface{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if bond != nil && bond["name"] != nil && bond["bond_mode"] != nil && bond["interfaces"] != nil {
				bonds = append(bonds, bond)
			}
			bond = nil
			continue
		}
		if strings.HasPrefix(line, "_uuid") {
			bond = make(map[string]interface{})
		}
		if bond == nil {
			continue
		}
		if strings.HasPrefix(line, "name") {
			bond["name"] = strings.TrimSpace(strings.TrimPrefix(line, "name               : "))
		}
		if strings.HasPrefix(line, "bond_mode") {
			bond["mode"] = strings.TrimSpace(strings.TrimPrefix(line, "bond_mode          : "))
		}
		if strings.HasPrefix(line, "interfaces") {
			members := strings.TrimSpace(strings.TrimPrefix(line, "interfaces         : "))
			members = strings.Trim(members, "[]")
			if members != "" {
				bond["members"] = strings.Split(members, ", ")
			} else {
				bond["members"] = []string{}
			}
		}
	}
	// 获取每个 Bond 所属的网桥
	for _, bond := range bonds {
		bondName := bond["name"].(string)
		bridge, err := getBondBridge(bondName)
		if err == nil {
			bond["bridge"] = bridge
		} else {
			bond["bridge"] = ""
		}
	}
	return bonds, nil
}

// getBondBridge 获取 Bond 所属的网桥
func getBondBridge(bondName string) (string, error) {
	// 获取所有网桥
	out, err := exec.Command("ovs-vsctl", "list", "bridge").Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	var currentBridge string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name") {
			currentBridge = strings.TrimSpace(strings.TrimPrefix(line, "name               : "))
		}
		if strings.HasPrefix(line, "ports") {
			ports := strings.TrimSpace(strings.TrimPrefix(line, "ports              : "))
			ports = strings.Trim(ports, "[]")
			if ports != "" {
				portList := strings.Split(ports, ", ")
				for _, port := range portList {
					if port == bondName {
						return currentBridge, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("bond %s not found in any bridge", bondName)
}
