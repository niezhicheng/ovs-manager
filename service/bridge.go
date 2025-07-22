package service

import (
	"fmt"
	"os/exec"
	"strings"
)

// ListBridges 调用 ovs-vsctl 列出所有 bridge
func ListBridges() ([]string, error) {
	cmd := exec.Command("ovs-vsctl", "list-br")
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

// AddBridge 新增 bridge
func AddBridge(name string) error {
	cmd := exec.Command("ovs-vsctl", "add-br", name)
	return cmd.Run()
}

// DeleteBridge 删除 bridge
func DeleteBridge(name string) error {
	cmd := exec.Command("ovs-vsctl", "del-br", name)
	return cmd.Run()
}

// SetNetFlow 设置 NetFlow
func SetNetFlow(bridge, target string, engineID int) error {
	args := []string{"set", "Bridge", bridge, fmt.Sprintf("netflow=@nf")}
	nfArgs := []string{"--", "--id=@nf", "create", "NetFlow", fmt.Sprintf("targets=[\"%s\"]", target)}
	if engineID != 0 {
		nfArgs = append(nfArgs, fmt.Sprintf("engine_id=%d", engineID))
	}
	cmd := exec.Command("ovs-vsctl", append(args, nfArgs...)...)
	return cmd.Run()
}

// SetSFlow 设置 sFlow
func SetSFlow(bridge string, targets []string, sampling, header, polling int, agent string) error {
	targetStrs := make([]string, len(targets))
	for i, t := range targets {
		targetStrs[i] = fmt.Sprintf("\"%s\"", t)
	}
	args := []string{"set", "Bridge", bridge, fmt.Sprintf("sflow=@sf")}
	sfArgs := []string{"--", "--id=@sf", "create", "sFlow", fmt.Sprintf("targets=[%s]", strings.Join(targetStrs, ","))}
	if sampling != 0 {
		sfArgs = append(sfArgs, fmt.Sprintf("sampling=%d", sampling))
	}
	if header != 0 {
		sfArgs = append(sfArgs, fmt.Sprintf("header=%d", header))
	}
	if polling != 0 {
		sfArgs = append(sfArgs, fmt.Sprintf("polling=%d", polling))
	}
	if agent != "" {
		sfArgs = append(sfArgs, fmt.Sprintf("agent=%s", agent))
	}
	cmd := exec.Command("ovs-vsctl", append(args, sfArgs...)...)
	return cmd.Run()
}

// SetStp 设置 STP
func SetStp(bridge string, enable bool) error {
	val := "false"
	if enable {
		val = "true"
	}
	cmd := exec.Command("ovs-vsctl", "set", "Bridge", bridge, fmt.Sprintf("stp_enable=%s", val))
	return cmd.Run()
}

// SetQos 设置 QoS
func SetQos(portName, qosType, maxRate string, queues map[string]string) error {
	args := []string{"set", "port", portName, fmt.Sprintf("qos=@newqos")}
	qosArgs := []string{"--", "--id=@newqos", "create", "qos", fmt.Sprintf("type=%s", qosType)}
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

// SetRstp 设置 RSTP
func SetRstp(bridge string, enable bool) error {
	val := "false"
	if enable {
		val = "true"
	}
	cmd := exec.Command("ovs-vsctl", "set", "Bridge", bridge, fmt.Sprintf("rstp_enable=%s", val))
	return cmd.Run()
}

// SetIpfix 设置 IPFIX
func SetIpfix(bridge string, targets []string, sampling, obsDomainID, obsPointID int) error {
	targetStrs := make([]string, len(targets))
	for i, t := range targets {
		targetStrs[i] = fmt.Sprintf("\"%s\"", t)
	}
	args := []string{"set", "Bridge", bridge, fmt.Sprintf("ipfix=@ipf")}
	ipfixArgs := []string{"--", "--id=@ipf", "create", "IPFIX", fmt.Sprintf("targets=[%s]", strings.Join(targetStrs, ","))}
	if sampling != 0 {
		ipfixArgs = append(ipfixArgs, fmt.Sprintf("sampling=%d", sampling))
	}
	if obsDomainID != 0 {
		ipfixArgs = append(ipfixArgs, fmt.Sprintf("obs_domain_id=%d", obsDomainID))
	}
	if obsPointID != 0 {
		ipfixArgs = append(ipfixArgs, fmt.Sprintf("obs_point_id=%d", obsPointID))
	}
	cmd := exec.Command("ovs-vsctl", append(args, ipfixArgs...)...)
	return cmd.Run()
}

// DumpFlows 查询流缓存
func DumpFlows(bridge string) (string, error) {
	cmd := exec.Command("ovs-ofctl", "dump-flows", bridge)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
} 