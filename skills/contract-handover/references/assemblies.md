# 合同组：识别、完整读取与交接

合同组须完整读取父组与有效成员，并满足当前权限、状态和影响范围条件，不能套用普通合同角色规则。有后继、未知关联组、未知灰度准入或缺失回读路径时不写入。父组类型与当前通道权限不能可靠确认时停止，不能把其他组类型的成功结果作为本次可执行依据。

参数表约定见 [调用参数约定](transport.md)；业务类型码与 ID 选择见 [资源类型](resource-types.md)。

先按 [父组 ID 定位](assembly-lookup.md) 使用已核父组入口或成员的真实 assemblyId；缺少映射时只读补齐。已确认组成员的状态和角色从本节 childrenInfo 取得，不强制调用普通 contract/info。父组详情的 contractAssemblyId 必须与定位 ID 一致；原选成员须在完整成员列表中按 contractId 精确匹配，且其 assemblyId 指向该父组。ID、类型或关系矛盾时停止，不换父组继续。

## 父合同组详情与回查

`GET /open-web/web/clm/contractAssembly/info`

请求形态示例（姓名和 ID 仅为占位，执行前按本次实际值替换）：

```http
GET /open-web/web/clm/contractAssembly/info?contractAssemblyId=30001&batchContractChildrenPageIndex=0&batchContractChildrenPageSize=100&needChildrenInfo=false
```
| 参数 | 位置 | 类型 | 本流程要求 | 值、来源与限制 |
|---|---|---|---|---|
| `contractAssemblyId` | query | string（内部 ID） | 必须传 | 已确认的父合同组 ID，不能填搜索结果外层 groupId。 |
| `batchContractChildrenPageIndex` | query | integer | 必须传 | 父组详情固定 `0`，不把详情接口当成员翻页接口。 |
| `batchContractChildrenPageSize` | query | integer | 必须传 | 固定 `100`。 |
| `needChildrenInfo` | query | boolean | 必须传 | 固定 `false`，完整成员通过 childrenInfo 分页读取，不依赖内嵌 contracts。 |

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

请求形态示例（姓名和 ID 仅为占位，执行前按本次实际值替换）：

```http
GET /open-web/web/clm/contractAssembly/childrenInfo?contractAssemblyId=30001&batchContractChildrenPageIndex=0&batchContractChildrenPageSize=100&needApproveRecord=false
```
| 参数 | 位置 | 类型 | 本流程要求 | 值、来源与限制 |
|---|---|---|---|---|
| `contractAssemblyId` | query | string（内部 ID） | 必须传 | 已确认的父合同组 ID，不能填搜索结果外层 groupId。 |
| `batchContractChildrenPageIndex` | query | integer | 必须传 | 从 0 开始，每页加 1。 |
| `batchContractChildrenPageSize` | query | integer | 必须传 | 固定 `100`。 |
| `needApproveRecord` | query | boolean | 必须传 | 固定 `false`，不请求成员审批记录；仍读取当前角色。 |

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

- 类型统一按 [资源类型表](resource-types.md) 解释。父组类型 10/17 与成员类型 11/0 不可互换；核对真实详情与成员关系，不能用搜索外层 groupId 代替父组 ID。
- 成员从 pageIndex=0、pageSize=100 逐页读取到空页，按 contractId 去重，与父组前后 childrenCount 核对。父组 contracts 可能只有首批，不能直接当全部。数量不符或漂移时重新核实，不能宣布成员完整。
- 只选择组内一份合同但用户明确其他成员不得改变时，说明整组影响并停止；不能静默把父组及其他成员加入已授权范围。
- 写入以父组真实 ID、类型和 groupId 形成一个 resourceList 入口。组整体确认，不能拆成成员来规避范围限制。
- 父组角色分别计算：实际操作人等于旧归属人且接收人不同，改变归属人；实际操作人等于旧申请人且接收人不同，改变申请人。服务端允许的第三方管理员分支以旧归属人计算。两个标记都不成立是无新增变化。只更换申请人的请求不借合同组路径变相执行。
- 所有有效成员跟随父组的两个角色标记，不能按每个成员自身角色关系重新计算。父组有效交接后按接收人首部门更新，无部门保留。子合同部门由后端处理，不作为读取或结果判据，不用父组部门或员工部门补值。
- 本次父组授权／新增历史 + 完整成员关系 + 每份成员角色和需求人回读，可以证明成员结果；不要求每份成员独立授权记录。成员缺失、角色不符或缺回读的对象待核实。
- 流程归属、旧授权停用、通知发送与送达分别核验。没有通知证据不自动判角色交接失败，也不能声称副作用全通过。没有合法回查方法时明确保留未核实项，不能先假定它们已发生。

## 参数与分页边界

两个接口的 contractAssemblyId 都是正十进制父组 ID；batchContractChildrenPageIndex 为非负整数，batchContractChildrenPageSize 为正整数，本流程固定 100。父组 info 只传 pageIndex=0；childrenInfo 才逐页加 1，不使用合同搜索的分页令牌。

needChildrenInfo=false 表示不请求内嵌完整成员，不保证响应中 contracts 一定为空；成员完整性以 childrenInfo 全分页加父组 childrenCount 核对。needApproveRecord=false 表示不要求成员审批记录，不是关闭当前角色回读。

childrenCount 应为非负整数；缺失、成员缺 ID、重复与计数不一致、前后总数变化都不能判断成员完整。nextId<=0/null 仅用于无后继父组条件，不是全部关联链已查全的证明。taskInstanceId、noticeId、mobile 等其他入口参数不在本流程中传入，也不借它们替换权限上下文。
