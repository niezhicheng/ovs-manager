// OVS 管理平台 API
// 结构参考 gin-vue-admin，仅做命令行管理
// Author: 你自己的名字
//
// 启动命令：go run main.go
//
// 健康检查接口：GET /ping
package main

import (
	"ovs-manager/router"
)

func main() {
	r := router.InitRouter()
	r.Run(":8080")
}
