package service

import (
	"fmt"
	"os/exec"
)

// AddVxlanPort 添加 VXLAN 端口
func AddVxlanPort(bridge, portName, remoteIP string, vni int, key, localIP string) error {
	args := []string{"add-port", bridge, portName, "--", "set", "interface", portName, "type=vxlan"}
	if remoteIP != "" {
		args = append(args, fmt.Sprintf("options:remote_ip=%s", remoteIP))
	}
	if vni != 0 {
		args = append(args, fmt.Sprintf("options:key=%d", vni))
	}
	if key != "" {
		args = append(args, fmt.Sprintf("options:key=%s", key))
	}
	if localIP != "" {
		args = append(args, fmt.Sprintf("options:local_ip=%s", localIP))
	}
	cmd := exec.Command("ovs-vsctl", args...)
	return cmd.Run()
} 