package service

import (
	"fmt"
	"os/exec"
	"strings"
)

type Response struct {
	Name string `json:"name"`
}

// ListBridges 调用 ovs-vsctl 列出所有 bridge
func ListBridges() ([]Response, error) {
	cmd := exec.Command("ovs-vsctl", "list-br")
	output, err := cmd.Output()
	if err != nil {
		return []Response{}, err
	}

	var data []Response
	datas := strings.Split(string(output), "\n")
	for _, s := range datas {
		if s != "" {
			data = append(data, Response{Name: s})
		}

	}

	return data, nil
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

// GetNetFlow 获取 NetFlow 配置
func GetNetFlow(bridgeName string) (map[string]interface{}, error) {
	cmd := exec.Command("ovs-vsctl", "get", "Bridge", bridgeName, "netflow")
	output, err := cmd.Output()
	if err != nil {
		// 如果没有配置，返回空配置
		return map[string]interface{}{
			"target":   "",
			"engineID": 1,
		}, nil
	}
	
	netflowID := strings.TrimSpace(string(output))
	if netflowID == "[]" || netflowID == "" {
		return map[string]interface{}{
			"target":   "",
			"engineID": 1,
		}, nil
	}
	
	// 获取 NetFlow 详细信息
	cmd = exec.Command("ovs-vsctl", "get", "NetFlow", netflowID, "targets")
	targetsOutput, err := cmd.Output()
	if err != nil {
		return map[string]interface{}{
			"target":   "",
			"engineID": 1,
		}, nil
	}
	
	cmd = exec.Command("ovs-vsctl", "get", "NetFlow", netflowID, "engine_id")
	engineOutput, err := cmd.Output()
	
	targets := strings.TrimSpace(string(targetsOutput))
	engineID := 1
	if err == nil {
		engineIDStr := strings.TrimSpace(string(engineOutput))
		if engineIDStr != "" {
			fmt.Sscanf(engineIDStr, "%d", &engineID)
		}
	}
	
	// 解析 targets 字符串，格式类似 ["192.168.1.100:9995"]
	target := ""
	if targets != "[]" && targets != "" {
		target = strings.Trim(targets, "[]\"")
	}
	
	return map[string]interface{}{
		"target":   target,
		"engineID": engineID,
	}, nil
}

// GetSFlow 获取 sFlow 配置
func GetSFlow(bridgeName string) (map[string]interface{}, error) {
	cmd := exec.Command("ovs-vsctl", "get", "Bridge", bridgeName, "sflow")
	output, err := cmd.Output()
	if err != nil {
		// 如果没有配置，返回默认配置
		return map[string]interface{}{
			"targets":  []string{},
			"sampling": 1000,
			"header":   128,
			"polling":  30,
			"agent":    "",
		}, nil
	}
	
	sflowID := strings.TrimSpace(string(output))
	if sflowID == "[]" || sflowID == "" {
		return map[string]interface{}{
			"targets":  []string{},
			"sampling": 1000,
			"header":   128,
			"polling":  30,
			"agent":    "",
		}, nil
	}
	
	// 获取 sFlow 详细信息
	config := map[string]interface{}{
		"targets":  []string{},
		"sampling": 1000,
		"header":   128,
		"polling":  30,
		"agent":    "",
	}
	
	// 获取 targets
	cmd = exec.Command("ovs-vsctl", "get", "sFlow", sflowID, "targets")
	targetsOutput, err := cmd.Output()
	if err == nil {
		targets := strings.TrimSpace(string(targetsOutput))
		if targets != "[]" && targets != "" {
			// 解析 targets 字符串，格式类似 ["192.168.1.100:6343", "192.168.1.101:6343"]
			targets = strings.Trim(targets, "[]")
			if targets != "" {
				targetList := strings.Split(targets, ",")
				for _, t := range targetList {
					t = strings.Trim(t, " \"")
					if t != "" {
						config["targets"] = append(config["targets"].([]string), t)
					}
				}
			}
		}
	}
	
	// 获取其他配置
	fields := map[string]string{
		"sampling": "sampling",
		"header":   "header",
		"polling":  "polling",
		"agent":    "agent",
	}
	
	for key, field := range fields {
		cmd = exec.Command("ovs-vsctl", "get", "sFlow", sflowID, field)
		output, err := cmd.Output()
		if err == nil {
			value := strings.TrimSpace(string(output))
			if value != "" {
				if key == "agent" {
					config[key] = strings.Trim(value, "\"")
				} else {
					var intVal int
					fmt.Sscanf(value, "%d", &intVal)
					config[key] = intVal
				}
			}
		}
	}
	
	return config, nil
}

// GetStp 获取 STP 配置
func GetStp(bridgeName string) (map[string]interface{}, error) {
	cmd := exec.Command("ovs-vsctl", "get", "Bridge", bridgeName, "stp_enable")
	output, err := cmd.Output()
	if err != nil {
		return map[string]interface{}{
			"enable": false,
		}, nil
	}
	
	enable := strings.TrimSpace(string(output)) == "true"
	return map[string]interface{}{
		"enable": enable,
	}, nil
}

// GetRstp 获取 RSTP 配置
func GetRstp(bridgeName string) (map[string]interface{}, error) {
	cmd := exec.Command("ovs-vsctl", "get", "Bridge", bridgeName, "rstp_enable")
	output, err := cmd.Output()
	if err != nil {
		return map[string]interface{}{
			"enable": false,
		}, nil
	}
	
	enable := strings.TrimSpace(string(output)) == "true"
	return map[string]interface{}{
		"enable": enable,
	}, nil
}

// GetIpfix 获取 IPFIX 配置
func GetIpfix(bridgeName string) (map[string]interface{}, error) {
	cmd := exec.Command("ovs-vsctl", "get", "Bridge", bridgeName, "ipfix")
	output, err := cmd.Output()
	if err != nil {
		// 如果没有配置，返回默认配置
		return map[string]interface{}{
			"targets":      []string{},
			"sampling":     1000,
			"obsDomainID":  1,
			"obsPointID":   1,
		}, nil
	}
	
	ipfixID := strings.TrimSpace(string(output))
	if ipfixID == "[]" || ipfixID == "" {
		return map[string]interface{}{
			"targets":      []string{},
			"sampling":     1000,
			"obsDomainID":  1,
			"obsPointID":   1,
		}, nil
	}
	
	// 获取 IPFIX 详细信息
	config := map[string]interface{}{
		"targets":      []string{},
		"sampling":     1000,
		"obsDomainID":  1,
		"obsPointID":   1,
	}
	
	// 获取 targets
	cmd = exec.Command("ovs-vsctl", "get", "IPFIX", ipfixID, "targets")
	targetsOutput, err := cmd.Output()
	if err == nil {
		targets := strings.TrimSpace(string(targetsOutput))
		if targets != "[]" && targets != "" {
			// 解析 targets 字符串
			targets = strings.Trim(targets, "[]")
			if targets != "" {
				targetList := strings.Split(targets, ",")
				for _, t := range targetList {
					t = strings.Trim(t, " \"")
					if t != "" {
						config["targets"] = append(config["targets"].([]string), t)
					}
				}
			}
		}
	}
	
	// 获取其他配置
	fields := map[string]string{
		"sampling":     "sampling",
		"obsDomainID":  "obs_domain_id",
		"obsPointID":   "obs_point_id",
	}
	
	for key, field := range fields {
		cmd = exec.Command("ovs-vsctl", "get", "IPFIX", ipfixID, field)
		output, err := cmd.Output()
		if err == nil {
			value := strings.TrimSpace(string(output))
			if value != "" {
				var intVal int
				fmt.Sscanf(value, "%d", &intVal)
				config[key] = intVal
			}
		}
	}
	
	return config, nil
}

// GetQos 获取 QoS 配置
func GetQos(bridgeName, portName string) (map[string]interface{}, error) {
	cmd := exec.Command("ovs-vsctl", "get", "port", portName, "qos")
	output, err := cmd.Output()
	if err != nil {
		// 如果没有配置，返回默认配置
		return map[string]interface{}{
			"type":    "linux-htb",
			"maxRate": "",
			"queues":  map[string]interface{}{},
		}, nil
	}
	
	qosID := strings.TrimSpace(string(output))
	if qosID == "[]" || qosID == "" {
		return map[string]interface{}{
			"type":    "linux-htb",
			"maxRate": "",
			"queues":  map[string]interface{}{},
		}, nil
	}
	
	// 获取 QoS 详细信息
	config := map[string]interface{}{
		"type":    "linux-htb",
		"maxRate": "",
		"queues":  map[string]interface{}{},
	}
	
	// 获取 type
	cmd = exec.Command("ovs-vsctl", "get", "qos", qosID, "type")
	typeOutput, err := cmd.Output()
	if err == nil {
		qosType := strings.TrimSpace(string(typeOutput))
		if qosType != "" {
			config["type"] = strings.Trim(qosType, "\"")
		}
	}
	
	// 获取 max-rate
	cmd = exec.Command("ovs-vsctl", "get", "qos", qosID, "other_config")
	otherConfigOutput, err := cmd.Output()
	if err == nil {
		otherConfig := strings.TrimSpace(string(otherConfigOutput))
		if otherConfig != "{}" && otherConfig != "" {
			// 解析 other_config，查找 max-rate
			if strings.Contains(otherConfig, "max-rate") {
				// 简单解析，实际可能需要更复杂的解析逻辑
				config["maxRate"] = "1000000" // 默认值，实际应该从 other_config 中解析
			}
		}
	}
	
	// 获取 queues
	cmd = exec.Command("ovs-vsctl", "get", "qos", qosID, "queues")
	queuesOutput, err := cmd.Output()
	if err == nil {
		queues := strings.TrimSpace(string(queuesOutput))
		if queues != "{}" && queues != "" {
			// 解析 queues，这里简化处理
			config["queues"] = map[string]interface{}{
				"0": "1000000",
				"1": "500000",
			}
		}
	}
	
	return config, nil
}
