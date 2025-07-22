package service

import (
	"fmt"
	"os/exec"
	"strings"
)

// ListPorts 列出指定 bridge 的所有端口
func ListPorts(bridge string) ([]string, error) {
	cmd := exec.Command("ovs-vsctl", "list-ports", bridge)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := []string{}
	for _, line := range string(output) {
		if line != '\n' {
			lines = append(lines, string(line))
		}
	}
	return lines, nil
}

// AddPort 向指定 bridge 添加端口，可指定类型
func AddPort(bridge, port, portType string) error {
	if portType == "internal" {
		cmd := exec.Command("ovs-vsctl", "add-port", bridge, port, "--", "set", "Interface", port, "type=internal")
		return cmd.Run()
	}
	if portType != "" {
		cmd := exec.Command("ovs-vsctl", "add-port", bridge, port, "--", "set", "Interface", port, "type="+portType)
		return cmd.Run()
	}
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, port)
	return cmd.Run()
}

// DeletePort 从指定 bridge 删除端口
func DeletePort(bridge, port string) error {
	cmd := exec.Command("ovs-vsctl", "del-port", bridge, port)
	return cmd.Run()
}

// BindPortToNetns 将端口绑定到指定命名空间
func BindPortToNetns(portName, netns string) error {
	cmd := exec.Command("ip", "link", "set", portName, "netns", netns)
	return cmd.Run()
}

// UnbindPortFromNetns 将端口解绑到主命名空间
func UnbindPortFromNetns(portName string) error {
	cmd := exec.Command("ip", "link", "set", portName, "netns", "1")
	return cmd.Run()
}

// SetPortUpDown 设置端口 up/down
func SetPortUpDown(netns, portName string, up bool) error {
	state := "down"
	if up {
		state = "up"
	}
	cmd := exec.Command("ip", "netns", "exec", netns, "ip", "link", "set", portName, state)
	return cmd.Run()
}

// SetPortAddr 给端口分配 IP 地址
func SetPortAddr(netns, portName, ip string) error {
	cmd := exec.Command("ip", "netns", "exec", netns, "ip", "addr", "add", ip, "dev", portName)
	return cmd.Run()
}

// SetPortVlanTag 设置端口 VLAN tag
func SetPortVlanTag(portName string, tag int) error {
	cmd := exec.Command("ovs-vsctl", "set", "port", portName, "tag="+string(rune(tag)))
	return cmd.Run()
}

// SetPortVlanMode 设置端口 VLAN mode
func SetPortVlanMode(portName, vlanMode string) error {
	cmd := exec.Command("ovs-vsctl", "set", "port", portName, fmt.Sprintf("vlan_mode=%s", vlanMode))
	return cmd.Run()
}

// SetPortTrunks 设置端口 trunks
func SetPortTrunks(portName string, trunks []int) error {
	trunksStr := make([]string, len(trunks))
	for i, t := range trunks {
		trunksStr[i] = fmt.Sprintf("%d", t)
	}
	cmd := exec.Command("ovs-vsctl", "set", "port", portName, fmt.Sprintf("trunks=%s", strings.Join(trunksStr, ",")))
	return cmd.Run()
}

// RemovePortProperty 移除端口属性
func RemovePortProperty(portName, property string, value interface{}) error {
	var valStr string
	switch v := value.(type) {
	case int:
		valStr = fmt.Sprintf("%d", v)
	case string:
		valStr = v
	case []int:
		strs := make([]string, len(v))
		for i, t := range v {
			strs[i] = fmt.Sprintf("%d", t)
		}
		valStr = strings.Join(strs, ",")
	case []string:
		valStr = strings.Join(v, ",")
	default:
		return fmt.Errorf("unsupported value type")
	}
	cmd := exec.Command("ovs-vsctl", "remove", "port", portName, property, valStr)
	return cmd.Run()
}

// AddPatchPort 添加 Patch Port
func AddPatchPort(bridge, portName, peer string) error {
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, portName, "--", "set", "interface", portName, "type=patch", fmt.Sprintf("options:peer=%s", peer))
	return cmd.Run()
}

// AddTunnelPort 添加 GRE/Geneve Tunnel Port
func AddTunnelPort(bridge, portName, typ string, options map[string]string) error {
	args := []string{"add-port", bridge, portName, "--", "set", "interface", portName, fmt.Sprintf("type=%s", typ)}
	for k, v := range options {
		args = append(args, fmt.Sprintf("options:%s=%s", k, v))
	}
	cmd := exec.Command("ovs-vsctl", args...)
	return cmd.Run()
}

// PortInfo 查询端口/interface 详细属性
func PortInfo(portName string) (string, error) {
	cmd := exec.Command("ovs-vsctl", "list", "interface", portName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// SetBfd 设置 BFD
func SetBfd(portName string, bfd map[string]string) error {
	args := []string{"set", "interface", portName}
	for k, v := range bfd {
		args = append(args, fmt.Sprintf("bfd:%s=%s", k, v))
	}
	cmd := exec.Command("ovs-vsctl", args...)
	return cmd.Run()
}

// SetCfm 设置 CFM (802.1ag)
func SetCfm(portName string, cfm map[string]string) error {
	args := []string{"set", "interface", portName}
	for k, v := range cfm {
		args = append(args, fmt.Sprintf("cfm:%s=%s", k, v))
	}
	cmd := exec.Command("ovs-vsctl", args...)
	return cmd.Run()
}

// SetMcastSnooping 设置组播监听
func SetMcastSnooping(bridge string, enable bool) error {
	val := "false"
	if enable {
		val = "true"
	}
	cmd := exec.Command("ovs-vsctl", "set", "Bridge", bridge, fmt.Sprintf("mcast_snooping_enable=%s", val))
	return cmd.Run()
}

// SetHfscQos 设置 HFSC QoS
func SetHfscQos(portName, maxRate string, queues map[string]string) error {
	args := []string{"set", "port", portName, "qos=@newqos"}
	qosArgs := []string{"--", "--id=@newqos", "create", "qos", "type=hfsc"}
	if maxRate != "" {
		qosArgs = append(qosArgs, fmt.Sprintf("other-config:max-rate=%s", maxRate))
	}
	if len(queues) > 0 {
		queueStrs := make([]string, 0, len(queues))
		for k, v := range queues {
			queueStrs = append(queueStrs, fmt.Sprintf("%s=%s", k, v))
		}
		qosArgs = append(qosArgs, fmt.Sprintf("queues=%s", strings.Join(queueStrs, ",")))
	}
	cmd := exec.Command("ovs-vsctl", append(args, qosArgs...)...)
	return cmd.Run()
}

// SetDatapathType 设置网桥 datapath_type
func SetDatapathType(bridge, datapathType string) error {
	cmd := exec.Command("ovs-vsctl", "set", "Bridge", bridge, fmt.Sprintf("datapath_type=%s", datapathType))
	return cmd.Run()
} 