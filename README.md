# OVS 管理平台后端

## 项目简介
本项目是基于 Go + Gin 的 Open vSwitch (OVS) 管理平台后端，支持 OVS 交换机、端口、Bond、流表、镜像、QoS、隧道、BFD、CFM、组播、STP/RSTP、NetFlow/sFlow/IPFIX、网络命名空间等全部常用与高级功能，接口风格参考 gin-vue-admin，便于前后端分离和自动化运维。

## 主要功能
- 交换机（Bridge）管理：增删查、NetFlow/sFlow/IPFIX、STP/RSTP、组播、datapath 切换等
- 端口（Port）管理：增删查、VLAN、trunk、patch、tunnel、up/down、IP、BFD、CFM、QoS、HFSC
- Bond 管理：增删查改、LACP、主备、负载均衡
- 端口镜像（Mirror）：增删查
- 流表（Flow）：增删查、OpenFlow 多表、流缓存
- 网络命名空间（Netns）：增删查
- 支持 GRE/VXLAN/Geneve 隧道
- 支持所有主流 OVS 高级特性

## API 文档
- OpenAPI 3.0 规范文档见 [`openapi.yaml`](./openapi.yaml)
- 可直接导入 [apiflox](https://apiflox.com/) 或 Swagger UI、Postman 等工具
- 主要接口均为 POST，参数通过 JSON body 传递，返回统一为 JSON

## 如何运行
1. 安装依赖
   ```bash
   go mod tidy
   ```
2. 启动后端服务
   ```bash
   go run main.go
   ```
3. 确保已安装并启动 openvswitch 服务
   ```bash
   sudo systemctl start openvswitch
   ```

## 如何导入 API 到 apiflox
1. 打开 [apiflox](https://apiflox.com/) 或本地 Swagger 工具
2. 选择导入 OpenAPI/Swagger 文件
3. 上传本项目根目录下的 `openapi.yaml`
4. 即可在线测试和查看所有接口

## 常见问题
- **端口/命名空间/网桥等操作需 root 权限**，建议以 sudo 运行后端或配置 sudoers
- **部分高级功能依赖 OVS 2.9+ 版本**，请确保 openvswitch 版本兼容
- **如需批量操作、定制返回结构、或前端对接示例，请联系开发者**

## 贡献与反馈
如有功能建议、Bug 反馈、或希望支持更多 OVS 高级特性，欢迎提 Issue 或 PR！

---

> 本项目仅为 OVS 管理后端，前端可结合 gin-vue-admin、element-plus、apiflox 等工具快速开发。
