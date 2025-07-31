package service

import (
	"os/exec"
	"strings"
)

// CreateNetns 创建网络命名空间
func CreateNetns(name string) error {
	cmd := exec.Command("ip", "netns", "add", name)
	return cmd.Run()
}

// DeleteNetns 删除网络命名空间
func DeleteNetns(name string) error {
	cmd := exec.Command("ip", "netns", "del", name)
	return cmd.Run()
}

// ListNetns 列出所有网络命名空间
func ListNetns() ([]string, error) {
	cmd := exec.Command("ip", "netns", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			result = append(result, fields[0])
		}
	}
	return result, nil
} 