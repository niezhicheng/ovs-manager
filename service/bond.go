package service

import (
	"fmt"
	"os/exec"
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
