package service

import (
	"fmt"
	"os/exec"
	"strings"
)

// PortInfo 端口信息结构体
type PortInfoResponse struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Up    bool   `json:"up"`
	Alias string `json:"alias"`
}

// PatchPortInfo 用于全局 patch 端口列表
type PatchPortInfo struct {
	Bridge string `json:"bridge"`
	Name   string `json:"name"`
	Peer   string `json:"peer"`
}

// ListPorts 列出指定 bridge 的所有端口，包含类型信息和状态
func ListPorts(bridge string) ([]PortInfoResponse, error) {
	// 获取端口列表
	cmd := exec.Command("ovs-vsctl", "list-ports", bridge)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ports []PortInfoResponse
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line != "" {
			// 获取端口类型
			portType := getPortType(line)
			// 获取端口状态
			portUp := getPortStatus(line)
			portAlias := getPortAlias(line)
			ports = append(ports, PortInfoResponse{
				Name:  line,
				Type:  portType,
				Up:    portUp,
				Alias: portAlias,
			})
		}
	}
	return ports, nil
}

// getPortType 获取端口类型
func getPortType(portName string) string {
	// 查询端口类型
	cmd := exec.Command("ovs-vsctl", "get", "Interface", portName, "type")
	output, err := cmd.Output()
	if err != nil {
		// 如果获取类型失败，默认为 normal
		return "normal"
	}

	portType := strings.TrimSpace(string(output))
	if portType == "" {
		return "normal"
	}
	return portType
}

// getPortStatus 获取端口状态（up/down）
func getPortStatus(portName string) bool {
	// 使用 ip link show 命令检查端口状态
	cmd := exec.Command("ip", "link", "show", portName)
	output, err := cmd.Output()
	if err != nil {
		// 如果命令失败，默认为down状态
		return false
	}

	// 检查输出中是否包含 "UP" 状态
	outputStr := string(output)
	return strings.Contains(outputStr, "UP")
}

// getPortAlias 获取端口别名（external-ids:ovs-port-name）
func getPortAlias(portName string) string {
	cmd := exec.Command("ovs-vsctl", "get", "Interface", portName, "external-ids:ovs-port-name")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	alias := strings.TrimSpace(string(output))
	if alias == "[]" {
		return ""
	}
	alias = strings.Trim(alias, "\"")
	return alias
}

// AddPort 向指定 bridge 添加端口，可指定类型
// 这个函数现在作为通用入口，根据类型调用专门的函数
func AddPort(bridge, port, portType, nicName string) error {
	switch portType {
	case "internal":
		return AddInternalPort(bridge, port)
	case "patch":
		return AddPatchPortWithoutPeer(bridge, port)
	case "vxlan":
		return AddVxlanPort(bridge, port)
	case "gre":
		return AddGrePort(bridge, port)
	case "tap":
		return AddTapPort(bridge, port)
	case "tun":
		return AddTunPort(bridge, port)
	case "bond":
		return AddBondPort(bridge, port)
	case "normal", "":
		return AddNormalPort(bridge, port, nicName)
	default:
		return AddCustomTypePort(bridge, port, portType)
	}
}

// AddNormalPort 添加普通端口（默认类型）
func AddNormalPort(bridge, port, nicName string) error {
	if nicName == "" {
		cmd := exec.Command("ovs-vsctl", "add-port", bridge, port)
		return cmd.Run()
	}
	// 使用网卡名称添加到网桥，但设置别名
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, nicName, "--", "set", "Interface", nicName, "external-ids:ovs-port-name="+port)
	return cmd.Run()
}

// AddInternalPort 添加内部端口
func AddInternalPort(bridge, port string) error {
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, port, "--", "set", "Interface", port, "type=internal")
	return cmd.Run()
}

// AddGrePort 添加GRE隧道端口
func AddGrePort(bridge, port string) error {
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, port, "--", "set", "Interface", port, "type=gre")
		return cmd.Run()
	}

// AddCustomTypePort 添加自定义类型端口
func AddCustomTypePort(bridge, port, portType string) error {
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, port, "--", "set", "Interface", port, "type="+portType)
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
func SetPortUpDown(portName string, up bool) error {
	state := "down"
	if up {
		state = "up"
	}
	cmd := exec.Command("ip", "link", "set", portName, state)
	return cmd.Run()
}

// SetPortAddr 给端口分配 IP 地址
func SetPortAddr(portName, ip string) error {
	cmd := exec.Command("ip", "addr", "add", ip, "dev", portName)
	return cmd.Run()
}

// GetPortAddrs 获取端口的IP地址列表
func GetPortAddrs(portName string) ([]string, error) {
	cmd := exec.Command("ip", "addr", "show", portName)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ips []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "inet ") {
			// 提取IP地址，格式如 "inet 192.168.1.10/24"
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				ips = append(ips, parts[1])
			}
		}
	}
	return ips, nil
}

// DeletePortAddr 删除端口的指定IP地址
func DeletePortAddr(portName, ip string) error {
	cmd := exec.Command("ip", "addr", "del", ip, "dev", portName)
	return cmd.Run()
}

// SetPortVlanTag 设置端口 VLAN tag
func SetPortVlanTag(portName string, tag int) error {
	cmd := exec.Command("ovs-vsctl", "set", "port", portName, fmt.Sprintf("tag=%d", tag))
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

// AddPatchPort 添加 patch 端口
func AddPatchPort(bridge, portName, peer string) error {
	if peer == "" {
		// 创建不设置对端的patch端口
		cmd := exec.Command("ovs-vsctl", "add-port", bridge, portName, "--", "set", "Interface", portName, "type=patch")
		return cmd.Run()
	}
	// 创建设置对端的patch端口
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, portName, "--", "set", "Interface", portName, "type=patch", "options:peer="+peer)
	return cmd.Run()
}

// AddPatchPortWithoutPeer 添加不设置对端的 patch 端口
func AddPatchPortWithoutPeer(bridge, portName string) error {
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, portName, "--", "set", "Interface", portName, "type=patch")
	return cmd.Run()
}

// SetPatchPortPeer 为patch端口设置对端
func SetPatchPortPeer(portName, peer string) error {
	cmd := exec.Command("ovs-vsctl", "set", "Interface", portName, "options:peer="+peer)
	return cmd.Run()
}

// AddPatchPortPair 一键成对创建 patch 端口
func AddPatchPortPair(bridgeA, portA, bridgeB, portB string) error {
	if err := AddPatchPort(bridgeA, portA, portB); err != nil {
		return err
	}
	if err := AddPatchPort(bridgeB, portB, portA); err != nil {
		return err
	}
	return nil
}

// AddVxlanPort 添加VXLAN隧道端口（基础版本，不设置参数）
func AddVxlanPort(bridge, port string) error {
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, port, "--", "set", "Interface", port, "type=vxlan")
	return cmd.Run()
}

// AddBondPort 添加Bond端口（基础版本，不设置成员）
func AddBondPort(bridge, port string) error {
	// 注意：bond端口通常需要成员，这里创建一个空的bond端口
	cmd := exec.Command("ovs-vsctl", "add-bond", bridge, port)
	return cmd.Run()
}

// AddBondPortWithMembers 添加 bond 端口（带成员和模式）
func AddBondPortWithMembers(bridge, portName string, members []string, mode string) error {
	args := []string{"add-bond", bridge, portName}
	args = append(args, members...)
	args = append(args, "bond_mode="+mode)
	cmd := exec.Command("ovs-vsctl", args...)
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

// AddTapPort 添加 tap 端口
func AddTapPort(bridge, portName string) error {
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, portName, "--", "set", "Interface", portName, "type=tap")
	return cmd.Run()
}

// AddTunPort 添加 tun 端口
func AddTunPort(bridge, portName string) error {
	cmd := exec.Command("ovs-vsctl", "add-port", bridge, portName, "--", "set", "Interface", portName, "type=tun")
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

// SetPortTypePeer 设置端口类型和 peer
func SetPortTypePeer(bridge, portName, typ, peer string) error {
	cmd := exec.Command("ovs-vsctl", "set", "Interface", portName, "type="+typ, "options:peer="+peer)
	return cmd.Run()
}

// SetPortAlias 设置端口别名（external-ids:ovs-port-name）
func SetPortAlias(portName, alias string) error {
	cmd := exec.Command("ovs-vsctl", "set", "Interface", portName, "external-ids:ovs-port-name="+alias)
	return cmd.Run()
}

// ListAllPatchPorts 返回所有 bridge 下的 patch 端口
func ListAllPatchPorts() ([]PatchPortInfo, error) {
	cmd := exec.Command("ovs-vsctl", "list-br")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	bridges := strings.Fields(string(output))
	var result []PatchPortInfo
	for _, bridge := range bridges {
		portsCmd := exec.Command("ovs-vsctl", "list-ports", bridge)
		portsOut, err := portsCmd.Output()
		if err != nil {
			continue
		}
		ports := strings.Fields(string(portsOut))
		for _, port := range ports {
			typeCmd := exec.Command("ovs-vsctl", "get", "interface", port, "type")
			typeOut, err := typeCmd.Output()
			if err != nil {
				continue
			}
			if strings.TrimSpace(string(typeOut)) == "patch" {
				peerCmd := exec.Command("ovs-vsctl", "get", "interface", port, "options:peer")
				peerOut, _ := peerCmd.Output()
				peer := strings.Trim(strings.TrimSpace(string(peerOut)), "\"")
				result = append(result, PatchPortInfo{Bridge: bridge, Name: port, Peer: peer})
			}
		}
	}
	return result, nil
}

// SetPortRoute 设置端口静态路由
func SetPortRoute(portName, destination, gateway string) error {
	cmd := exec.Command("ip", "route", "add", destination, "via", gateway, "dev", portName)
	return cmd.Run()
}

// DeletePortRoute 删除端口静态路由
func DeletePortRoute(portName, destination, gateway string) error {
	cmd := exec.Command("ip", "route", "del", destination, "via", gateway, "dev", portName)
	return cmd.Run()
}

// GetPortRoutes 获取端口路由列表
func GetPortRoutes(portName string) ([]string, error) {
	cmd := exec.Command("ip", "route", "show", "dev", portName)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var routes []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			routes = append(routes, line)
		}
	}
	return routes, nil
}
