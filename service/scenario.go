package service

import (
	"fmt"
)

// ExecuteScenarioStep 统一调度场景步骤
// 返回 error, output
func ExecuteScenarioStep(action string, params map[string]interface{}) (error, interface{}) {
	switch action {
	case "add_bridge":
		name, _ := params["name"].(string)
		return AddBridge(name), nil
	case "delete_bridge":
		name, _ := params["name"].(string)
		return DeleteBridge(name), nil
	case "add_port":
		bridge, _ := params["bridge"].(string)
		portName, _ := params["portName"].(string)
		portType, _ := params["type"].(string)
		return AddPort(bridge, portName, portType), nil
	case "delete_port":
		bridge, _ := params["bridge"].(string)
		portName, _ := params["portName"].(string)
		return DeletePort(bridge, portName), nil
	case "set_port_vlan":
		portName, _ := params["portName"].(string)
		tag, _ := toInt(params["tag"])
		return SetPortVlanTag(portName, tag), nil
	case "set_port_vlan_mode":
		portName, _ := params["portName"].(string)
		vlanMode, _ := params["vlanMode"].(string)
		return SetPortVlanMode(portName, vlanMode), nil
	case "set_port_trunks":
		portName, _ := params["portName"].(string)
		trunks, _ := toIntSlice(params["trunks"])
		return SetPortTrunks(portName, trunks), nil
	case "add_patch_port":
		bridge, _ := params["bridge"].(string)
		portName, _ := params["portName"].(string)
		peer, _ := params["peer"].(string)
		return AddPatchPort(bridge, portName, peer), nil
	case "add_bond":
		bridge, _ := params["bridge"].(string)
		bondName, _ := params["bondName"].(string)
		slaves, _ := toStringSlice(params["slaves"])
		bondMode, _ := params["bondMode"].(string)
		lacp, _ := params["lacp"].(string)
		otherOptions, _ := toStringMap(params["otherOptions"])
		return AddBond(bridge, bondName, slaves, bondMode, lacp, otherOptions), nil
	case "set_bfd":
		portName, _ := params["portName"].(string)
		bfd, _ := toStringMap(params["bfd"])
		return SetBfd(portName, bfd), nil
	case "set_cfm":
		portName, _ := params["portName"].(string)
		cfm, _ := toStringMap(params["cfm"])
		return SetCfm(portName, cfm), nil
	case "set_qos":
		portName, _ := params["portName"].(string)
		typeStr, _ := params["type"].(string)
		maxRate, _ := params["maxRate"].(string)
		queues, _ := toStringMap(params["queues"])
		return SetQos(portName, typeStr, maxRate, queues), nil
	case "set_hfsc_qos":
		portName, _ := params["portName"].(string)
		maxRate, _ := params["maxRate"].(string)
		queues, _ := toStringMap(params["queues"])
		return SetHfscQos(portName, maxRate, queues), nil
	case "add_tunnel_port":
		bridge, _ := params["bridge"].(string)
		portName, _ := params["portName"].(string)
		typeStr, _ := params["type"].(string)
		options, _ := toStringMap(params["options"])
		return AddTunnelPort(bridge, portName, typeStr, options), nil
	case "set_netflow":
		bridge, _ := params["bridge"].(string)
		target, _ := params["target"].(string)
		engineID, _ := toInt(params["engineID"])
		return SetNetFlow(bridge, target, engineID), nil
	case "set_sflow":
		bridge, _ := params["bridge"].(string)
		targets, _ := toStringSlice(params["targets"])
		sampling, _ := toInt(params["sampling"])
		header, _ := toInt(params["header"])
		polling, _ := toInt(params["polling"])
		agent, _ := params["agent"].(string)
		return SetSFlow(bridge, targets, sampling, header, polling, agent), nil
	case "set_stp":
		bridge, _ := params["bridge"].(string)
		enable, _ := params["enable"].(bool)
		return SetStp(bridge, enable), nil
	case "set_rstp":
		bridge, _ := params["bridge"].(string)
		enable, _ := params["enable"].(bool)
		return SetRstp(bridge, enable), nil
	case "set_ipfix":
		bridge, _ := params["bridge"].(string)
		targets, _ := toStringSlice(params["targets"])
		sampling, _ := toInt(params["sampling"])
		obsDomainID, _ := toInt(params["obsDomainID"])
		obsPointID, _ := toInt(params["obsPointID"])
		return SetIpfix(bridge, targets, sampling, obsDomainID, obsPointID), nil
	case "set_mcast_snooping":
		bridge, _ := params["bridge"].(string)
		enable, _ := params["enable"].(bool)
		return SetMcastSnooping(bridge, enable), nil
	case "set_datapath_type":
		bridge, _ := params["bridge"].(string)
		datapathType, _ := params["datapathType"].(string)
		return SetDatapathType(bridge, datapathType), nil
	case "add_mirror":
		bridge, _ := params["bridge"].(string)
		name, _ := params["name"].(string)
		srcPorts, _ := toStringSlice(params["selectSrcPorts"])
		dstPorts, _ := toStringSlice(params["selectDstPorts"])
		var selectVlan *int
		if v, ok := params["selectVlan"]; ok {
			vint, _ := toInt(v)
			selectVlan = &vint
		}
		outputPort, _ := params["outputPort"].(string)
		var outputVlan *int
		if v, ok := params["outputVlan"]; ok {
			vint, _ := toInt(v)
			outputVlan = &vint
		}
		selectAll, _ := params["selectAll"].(bool)
		return AddMirror(bridge, name, srcPorts, dstPorts, selectVlan, outputPort, outputVlan, selectAll), nil
	case "delete_mirror":
		bridge, _ := params["bridge"].(string)
		name, _ := params["name"].(string)
		return DeleteMirror(bridge, name), nil
	case "add_flow":
		bridge, _ := params["bridge"].(string)
		flow, _ := params["flow"].(string)
		return AddFlowV2(bridge, flow), nil
	case "delete_flow":
		bridge, _ := params["bridge"].(string)
		match, _ := params["match"].(string)
		return DeleteFlowV2(bridge, match), nil
	case "create_netns":
		name, _ := params["name"].(string)
		return CreateNetns(name), nil
	case "delete_netns":
		name, _ := params["name"].(string)
		return DeleteNetns(name), nil
	default:
		return fmt.Errorf("unsupported action: %s", action), nil
	}
}

// 工具函数
func toInt(v interface{}) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case float64:
		return int(val), true
	case float32:
		return int(val), true
	case string:
		var i int
		_, err := fmt.Sscanf(val, "%d", &i)
		return i, err == nil
	}
	return 0, false
}

func toIntSlice(v interface{}) ([]int, bool) {
	arr, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	res := make([]int, len(arr))
	for i, x := range arr {
		res[i], _ = toInt(x)
	}
	return res, true
}

func toStringSlice(v interface{}) ([]string, bool) {
	arr, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	res := make([]string, len(arr))
	for i, x := range arr {
		res[i], _ = x.(string)
	}
	return res, true
}

func toStringMap(v interface{}) (map[string]string, bool) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, false
	}
	res := make(map[string]string, len(m))
	for k, x := range m {
		if s, ok := x.(string); ok {
			res[k] = s
		} else {
			res[k] = fmt.Sprintf("%v", x)
		}
	}
	return res, true
} 