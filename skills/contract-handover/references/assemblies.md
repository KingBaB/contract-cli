# 合同组：识别、完整读取与受控验证

合同组属于扩展场景，不能因普通合同通过就自动开放。当前只保留 test 下具体清单、完整成员可读取且无后继父组的受控验证；有后继、未知关联组、未知灰度准入或缺失回读路径时不写入。type17 的网关／权限表现未完成运行验证，不按 type10 的样本放行。日常请求可先只读识别并说明范围。

## 父合同组详情与回查

`GET /open-web/web/clm/contractAssembly/info`

| 参数 | 位置 | 必填 | 值或来源 |
|---|---|---|---|
| `contractAssemblyId` | query | 是 | 已确认的父合同组 ID，不能填搜索结果外层 groupId。 |
| `batchContractChildrenPageIndex` | query | 是 | 固定 `0`。 |
| `batchContractChildrenPageSize` | query | 是 | 固定 `100`。 |
| `needChildrenInfo` | query | 是 | 固定 `false`。 |

关键响应字段（路径相对完整响应）：

| 含义 | 路径 | 类型 | 可为空 |
|---|---|---|---|
| 父组 ID | `$.data.contractAssemblyId` | string | 否 |
| 名称 | `$.data.contractAssemblyName` | string | 是 |
| 业务类型码 | `$.data.businessTypeCode` | integer | 否 |
| 状态码 | `$.data.contractStatusCode` | integer | 否 |
| 关联 groupId | `$.data.groupId` | string | 是 |
| 前序 ID | `$.data.previousId` | string | 是 |
| 后继 ID | `$.data.nextId` | string | 是 |
| 当前归属人 ID | `$.data.ownerEmployee.employeeId` | string | 是 |
| 当前申请人 ID | `$.data.submitterEmployee.employeeId` | string | 是 |
| 合同部门 ID | `$.data.department.departmentId` | string | 是 |
| 有效成员总数 | `$.data.childrenCount` | integer | 是 |
| 成员列表 | `$.data.contracts` | array | 是 |
| 交接历史 | `$.data.handover.records` | array | 是 |

## 合同组成员与回查

`GET /open-web/web/clm/contractAssembly/childrenInfo`

| 参数 | 位置 | 必填 | 值或来源 |
|---|---|---|---|
| `contractAssemblyId` | query | 是 | 已确认的父合同组 ID，不能填搜索结果外层 groupId。 |
| `batchContractChildrenPageIndex` | query | 是 | 从 0 开始，每页加 1。 |
| `batchContractChildrenPageSize` | query | 是 | 固定 `100`。 |
| `needApproveRecord` | query | 是 | 固定 `false`。 |

关键响应字段（路径相对完整响应）：

| 含义 | 路径 | 类型 | 可为空 |
|---|---|---|---|
| 成员列表 | `$.data.contracts` | array | 否 |
| 成员合同 ID | `$.data.contracts[*].contractId` | array | 否 |
| 租户 ID | `$.data.contracts[*].tenantId` | array | 是 |
| 父组 ID | `$.data.contracts[*].assemblyId` | array | 是 |
| 关联 groupId | `$.data.contracts[*].groupId` | array | 是 |
| 业务类型码 | `$.data.contracts[*].businessTypeCode` | array | 是 |
| 合同编号 | `$.data.contracts[*].contractNumber` | array | 是 |
| 合同名称 | `$.data.contracts[*].contractName` | array | 是 |
| 状态码 | `$.data.contracts[*].contractStatusCode` | array | 是 |
| 前序 ID | `$.data.contracts[*].previousId` | array | 是 |
| 后继 ID | `$.data.contracts[*].nextId` | array | 是 |
| 归属人 ID | `$.data.contracts[*].ownerEmployee.employeeId` | array | 是 |
| 申请人 ID | `$.data.contracts[*].submitterEmployee.employeeId` | array | 是 |
| 创建人 ID | `$.data.contracts[*].createEmployee.employeeId` | array | 是 |
| 需求人信息 | `$.data.contracts[*].demandEmployees` | array | 是 |
| 需求人 ID 列表 | `$.data.contracts[*].demandEmployees[*].employeeId` | array | 是 |

## 组范围与结果规则

- 类型示例：10 为合同组申请，17 为模板批量合同，11 为组下主合同；这些不是全部类型白名单。核对真实详情与成员关系，不能用搜索外层 groupId 代替父组 ID。
- 成员从 pageIndex=0、pageSize=100 逐页读取到空页，按 contractId 去重，与父组前后 childrenCount 核对。父组 contracts 可能只有首批，不能直接当全部。数量不符或漂移时重新核实，不能宣布成员完整。
- 只选择组内一份合同但用户明确其他成员不得改变时，说明整组影响并停止；不能静默把父组及其他成员加入已授权范围。
- 写入以父组真实 ID、类型和 groupId 形成一个 resourceList 入口。组整体确认，不能拆成成员来规避范围限制。
- 父组角色分别计算：实际操作人等于旧归属人且接收人不同，改变归属人；实际操作人等于旧申请人且接收人不同，改变申请人。服务端允许的第三方管理员分支以旧归属人计算。两个标记都不成立是无新增变化。只更换申请人的请求不借合同组路径变相执行。
- 所有有效成员跟随父组的两个角色标记，不能按每个成员自身角色关系重新计算。父组有效交接后按接收人首部门更新，无部门保留。子合同部门由后端处理，不作为读取或结果判据，不用父组部门或员工部门补值。
- 本次父组授权／新增历史 + 完整成员关系 + 每份成员角色和需求人回读，可以证明成员结果；不要求每份成员独立授权记录。成员缺失、角色不符或缺回读的对象待核实。
- 流程归属、旧授权停用、通知发送与送达分别核验。没有通知证据不自动判角色交接失败，也不能声称副作用全通过。没有合法回查方法时明确保留未核实项，不能先假定它们已发生。
