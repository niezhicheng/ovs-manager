# OVS 管理平台 API 文档说明

本目录用于存放 OVS 管理平台的 OpenAPI 规范文档及接口分组说明，便于前端、测试、运维等团队成员查阅和对接。

## 目录结构
- `openapi.yaml`：完整 OpenAPI 3.0 规范文档（可导入 apiflox、Swagger UI、Postman 等）
- `README.md`：接口分组与路由说明（本文件）

## 路由分组与主要接口

### 1. 交换机（Bridge）相关
- `/api/ovs/bridge/list`         查询网桥列表
- `/api/ovs/bridge/add`          新增网桥
- `/api/ovs/bridge/delete`       删除网桥
- `/api/ovs/bridge/set-netflow`  设置 NetFlow
- `/api/ovs/bridge/set-sflow`    设置 sFlow
- `/api/ovs/bridge/set-stp`      设置 STP
- `/api/ovs/bridge/set-rstp`     设置 RSTP
- `/api/ovs/bridge/set-ipfix`    设置 IPFIX
- `/api/ovs/bridge/set-mcast-snooping` 组播监听
- `/api/ovs/bridge/set-datapath-type`  datapath 切换
- `/api/ovs/bridge/dump-flows`   查询流缓存

### 2. 端口（Port）相关
- `/api/ovs/port/list`           查询端口列表
- `/api/ovs/port/add`            新增端口
- `/api/ovs/port/delete`         删除端口
- `/api/ovs/port/set-vlan`       设置 VLAN tag
- `/api/ovs/port/set-vlan-mode`  设置 VLAN mode
- `/api/ovs/port/set-trunks`     设置 trunks
- `/api/ovs/port/remove-property`移除端口属性
- `/api/ovs/port/info`           查询端口属性
- `/api/ovs/patch/add`           添加 Patch Port
- `/api/ovs/tunnel/add`          添加 Tunnel Port（GRE/Geneve/VXLAN）
- `/api/ovs/port/set-bfd`        设置 BFD
- `/api/ovs/port/set-cfm`        设置 CFM
- `/api/ovs/port/set-qos`        设置 QoS
- `/api/ovs/port/set-hfsc-qos`   设置 HFSC QoS
- `/api/ovs/port/updown`         设置 up/down
- `/api/ovs/port/addr`           分配 IP

### 3. Bond 相关
- `/api/ovs/bond/add`            新增 Bond
- `/api/ovs/bond/set`            设置 Bond 属性
- `/api/ovs/bond/show`           查询 Bond 状态
- `/api/ovs/bond/delete`         删除 Bond

### 4. 端口镜像（Mirror）相关
- `/api/ovs/mirror/add`          新增端口镜像
- `/api/ovs/mirror/delete`       删除端口镜像
- `/api/ovs/mirror/list`         查询端口镜像

### 5. 流表（Flow）相关
- `/api/ovs/flow/list-v2`        查询流表规则
- `/api/ovs/flow/add-v2`         添加流表规则
- `/api/ovs/flow/delete-v2`      删除流表规则

### 6. 网络命名空间（Netns）相关
- `/api/netns/create`            新增命名空间
- `/api/netns/delete`            删除命名空间
- `/api/netns/list`              查询命名空间列表

### 7. VXLAN 相关
- `/api/ovs/vxlan/add`           新增 VXLAN 端口
- `/api/ovs/vxlan/delete`        删除 VXLAN 端口

### 8. 场景引导（Scenario）相关
- `/api/ovs/scenario/apply`      场景引导式一键操作（支持模板+参数覆盖、自定义步骤）
  - 支持 scenario+params 组合，详见 openapi.yaml

## 如何使用 openapi.yaml
1. 打开 apiflox、Swagger UI、Postman 等工具
2. 导入本目录下的 `openapi.yaml`
3. 可在线测试所有接口，查看参数和响应结构

## 场景引导接口说明
- 支持一键部署常见 OVS 网络场景（如 VXLAN 隔离、Patch+Trunk 等）
- 支持自定义步骤，或用 params 覆盖模板参数，兼顾易用和灵活
- 返回每一步详细结果，便于前端引导和自动化运维

---
如需补充接口示例、响应示例、或有其它文档需求，请联系开发者。 