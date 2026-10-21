# 合同定位、搜索衔接与详情

## 复用已有搜索能力

用户刚完成查询并选择了合同，随后要求“把这几份交给某人”时，复用这些候选和原查询语义，不重新搜一遍。只有候选缺失、来源不明、环境/身份不一致或信息已变化时才重新定位。

需要新查询时读取已有 [合同搜索 Skill](../../contract-cli-contract-search/SKILL.md)，由它解析名称、编号、人员/部门、自定义字段、关键词与分页。固定使用与交接相同的 user 身份及环境，不因结果为空、权限错误或搜索能力不足而切到 app。该 Skill 不可用或其要求与当前环境不兼容时，说明缺少的能力并停止定位。仅在已选组对象缺少父组映射时，按 [父组 ID 补充查询](assembly-lookup.md) 补齐关系，不用它替代原合同搜索。

已选集合以用户确认的合同为准。搜索按组返回/展开的成员不一定都满足原筛选条件，必须核对每条实际对象。分页、总数和查询是否完整沿用搜索 Skill 的判断；“列出了几份”不等于查全。用户要求全部结果时先完成搜索与范围核对，再让用户确认固定清单，不向交接接口提交动态查询条件。

## 候选交接约定

下面是 Agent 在本轮任务中保留的业务上下文，不是需要发送到任何接口的新 JSON 字段，也不要求用户填写：

| 信息 | 来源与要求 |
|---|---|
| 环境与配置 | 已确认的环境、实际服务来源、配置目录和 profile；同名 profile 或同一合同编号不能证明环境相同 |
| 租户与身份 | 当前已授权 user 上下文及可核实的租户；不能把 app 搜索结果直接交给 user 写入，也不借用其他账号 |
| 原查询与选择 | 保留用户条件、排除项、候选来源、分页是否完整以及用户选中的具体对象 |
| 标识与展示 | 原样保留候选 contract_id、contract_number、contract_name、business_type_code 及其来源；缺失字段不推算、不伪造映射 |
| 已有组关系 | 保留实际返回的 assemblyId／assembly_id 及其字段来源；若有 Presentation，同时保留选中行 row_ref 与同一实体绑定，不把外层 group_id 当父组 ID |
| ID 语义 | 记录该搜索通道是否明确返回内部合同实体 ID；未知时保持不确定，不因为是数字字符串就判为内部 ID |

跨 Skill 传递这些必要信息即可，不传 Token、Cookie 或完整人员数据。上下文缺失且不能从本轮已授权配置核实时，停止交接，不靠重新询问同一姓名或编号掩盖身份差异。

## ID 映射与详情回读

1. 先按候选真实业务类型及组关系分流：已核父组 type=10/17 读取父组详情；已确认成员 type=11 或实际有正 assemblyId 的合同读取对应父组和完整成员。缺少父组 ID 时按 [补充查询](assembly-lookup.md) 定位；不把普通 contract/info 作为已确认组成员的前置条件。
2. 来源返回开放 ID、编码或不明语义的 contract_id 时，只把它作为候选标识。必须通过已有且已确认的映射能力取得内部 ID；本 Skill 不提供猜测、解码或数字转换映射。无可靠映射时说明无法继续，不让用户手抄内部 ID，不回退到另一环境搜索。
3. 普通合同来源已明确为内部合同 ID，且环境、租户和 user 上下文一致时，才能作为 contract/info 的 contractId 查询值。回读后核对 data.contractId、可用编号／名称、真实 businessTypeCode、groupId、assemblyId、当前角色与状态。snake_case 与 camelCase 只是字段名差异，不是 ID 等价证明；不一致或必需字段缺失时暂停。普通详情异常本身不能证明是组成员；已核组成员不为取得角色而反复调用普通详情。
4. 写入 businessId、businessTypeCode、parentId 只从本次交接详情或父组详情构造。不能把搜索 group_id 直接重命名为 parentId，也不能把 business_type_code 当作后端已批准交接的证据。

已有来源明确的内部 ID 可直接进入详情回读，无需先调用搜索。普通合同的 contractId 与父组 contractAssemblyId 不互换；识别出真实组成员后按 [合同组](assemblies.md) 确认父组与完整成员。

参数表约定见 [调用参数约定](transport.md)。

## 普通合同详情与回查

`GET /open-web/web/clm/contract/info`

请求形态示例（姓名和 ID 仅为占位，执行前按本次实际值替换）：

```http
GET /open-web/web/clm/contract/info?contractId=20001&tianyancha=false&generateContractName=false
```
| 参数 | 位置 | 类型 | 本流程要求 | 值、来源与限制 |
|---|---|---|---|---|
| `contractId` | query | string（内部 ID） | 必须传 | 已定位的普通合同内部 ID；不能填合同编号、聚合组 ID 或父组 ID。 |
| `tianyancha` | query | boolean | 必须传 | 固定 `false`，不启用天眼查扩展。 |
| `generateContractName` | query | boolean | 必须传 | 固定 `false`，不触发合同名称生成。 |

关键响应字段（路径相对完整响应）：

| 含义 | 路径 | 类型 | 可为空 |
|---|---|---|---|
| 合同 ID | `$.data.contractId` | string | 否 |
| 名称 | `$.data.contractName` | string | 是 |
| 业务类型码 | `$.data.businessTypeCode` | integer | 否 |
| 状态码 | `$.data.contractStatusCode` | integer | 否 |
| 状态名称 | `$.data.contractStatusName` | string | 是 |
| 关联 groupId | `$.data.groupId` | string | 是 |
| 父组 ID | `$.data.assemblyId` | string | 是 |
| 前序 ID | `$.data.previousId` | string | 是 |
| 后继 ID | `$.data.nextId` | string | 是 |
| 当前归属人 ID | `$.data.ownerEmployee.employeeId` | string | 是 |
| 归属人展示名 | `$.data.ownerEmployee.displayName` | string | 是 |
| 当前申请人 ID | `$.data.submitterEmployee.employeeId` | string | 是 |
| 申请人展示名 | `$.data.submitterEmployee.displayName` | string | 是 |
| 创建人 ID | `$.data.createEmployee.employeeId` | string | 是 |
| 合同部门 ID | `$.data.department.departmentId` | string | 是 |
| 需求人 ID 列表 | `$.data.demandEmployees[*].employeeId` | array | 是 |
| 交接历史 | `$.data.handover.records` | array | 是 |
| 交接前人员 | `$.data.handover.records[*].preAuthorizedEmployee` | object | 是 |
| 交接后人员 | `$.data.handover.records[*].authorizedEmployee` | object | 是 |
| 操作人 | `$.data.handover.records[*].approvedEmployee` | object | 是 |
| 历史记录时间 | `$.data.handover.records[*].createTime` | string | 是 |
| 授权来源码 | `$.data.handover.records[*].authorizationSourceTypeCode` | integer | 是 |

## 如何使用详情

- 写前和写后使用同一接口与参数，记录当前归属人、申请人、类型、状态、groupId、assemblyId 和交接历史。创建人不等于当前归属人。
- `assemblyId` 为真实父组 ID 时读取 [合同组](assemblies.md)，不直接把该成员当独立入口。type=0 也可能属于组，不能只凭类型判断。
- 普通独立入口要求 assemblyId 明确为 null 或 0；缺失或类型关系不明时停止该写入。groupId 与 assemblyId 不能互换；前后链为空也不能证明没有同组关联记录。
- 普通交接的角色规则见 [交接与结果核对](handover.md)。状态码 1 为已取消，不交接；已知普通准入范围为 2～17、另含 23，但状态落入范围不独自构成交接资格。草稿、协商等依赖未核实配置的分支不猜测放行，仍由服务端最终判定。
- 查询历史时分别显示交接前人员、接收人、操作人和时间，不用当前归属人补齐历史；空记录只表示接口未返回记录。时间单位或时区未知时保留原值。
- 不改用会更新交易方状态的聚合详情接口，也不通过猜参数、换身份来绕过详情异常。

## ID 来源与缺失处理

contractId、groupId、assemblyId、previousId、nextId 分别是合同实体、关联链、父合同组、前序、后继标识，不可互换。普通入口的 businessId 取 contractId，parentId 取该详情 groupId；父组入口按 [合同组](assemblies.md) 取值。搜索外层聚合 groupId 只能用于理解搜索结果，不用于构造 parentId。

详情中的关键 ID、businessTypeCode 或当前角色缺失时，停止该对象写入，不用 0、空字符串或名称补齐。assemblyId=0/null 只表示未指向已知父组，不证明不存在同 groupId 的关联资源。业务类型码完整对照与提交边界见 [资源类型](resource-types.md)。

交接历史的 authorizationSourceTypeCode 按接口原值保留，本流程没有完整来源枚举表，不用任意来源记录证明本次交接。姓名为空时可展示已核实的人员 ID，不以创建人或操作人补成当前归属人。
