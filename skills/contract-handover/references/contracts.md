# 合同定位与普通合同详情

已有可信内部合同 ID 时直接读取详情；只有编号或名称时先做有条件搜索。

## 按编号或名称定位合同

`POST /open-web/web/clm/es/search/contractList`

| 参数 | 位置 | 必填 | 值或来源 |
|---|---|---|---|
| `pageIndex` | body | 是 | 固定 `0`。 |
| `pageSize` | body | 是 | 固定 `100`。 |
| `searchTabCode` | body | 是 | 全部合同为 0，我的合同为 1；沿用用户指定范围。 |
| `conditionUnits` | body | 否 | 名称字段 CONTRACT_NAME_FIELD，编号字段 CONTRACT_NUMBER_FIELD；searchValue 为原始字符串，unionType 固定 SHOULD；有关键词时不得省略。 |
| `filterUnits` | body | 否 | 仅传已核实字段与语义的筛选条件；没有相应字段配方时停止该筛选，不猜字段 code。 |
| `allConditionUnits` | body | 是 | 固定 `false`。 |
| `searchAfter` | body | 否 | 首轮省略；续页原样传前页数组，保留每项类型和整数精度。 |
| `sortType` | body | 是 | 固定 `"SUBMITTED_TIME"`。 |
| `order` | body | 是 | 固定 `"DESC"`。 |
| `withStatistics` | body | 是 | 固定 `false`。 |

关键响应字段（路径相对完整响应）：

| 含义 | 路径 | 类型 | 可为空 |
|---|---|---|---|
| 聚合组 | `$.data.content` | array | 是 |
| 合同入口 | `$.data.content[*].searchContractListItems[*]` | array | 是 |
| 合同 ID | `$.data.content[*].searchContractListItems[*].contractId` | string | 是 |
| 合同编号 | `$.data.content[*].searchContractListItems[*].contractNumber` | string | 是 |
| 合同名称 | `$.data.content[*].searchContractListItems[*].contractName` | string | 是 |
| 业务类型码 | `$.data.content[*].searchContractListItems[*].businessTypeCode` | integer | 是 |
| 所属父组 ID | `$.data.content[*].searchContractListItems[*].assemblyId` | string | 是 |
| 聚合组总数 | `$.data.filteredTotalElements` | integer | 是 |
| 续页游标 | `$.data.searchAfter` | array | 是 |
| 分页能力 | `$.data.paginationCapability` | object | 是 |
| 各 tab 统计 | `$.data.totalElements` | array | 是 |

编号搜索请求示例（替换示例编号，并按用户范围选择 searchTabCode）：

```json
{"pageIndex":0,"pageSize":100,"searchTabCode":1,"conditionUnits":[{"searchField":"CONTRACT_NUMBER_FIELD","searchValue":"用户指定编号","unionType":"SHOULD"}],"allConditionUnits":false,"sortType":"SUBMITTED_TIME","order":"DESC","withStatistics":false}
```

搜索是候选定位，不保证精确匹配；逐条比较 contractNumber 或名称，重名时让用户选择。零结果只能说明本搜索范围未命中，不能证明草稿或协商记录不存在。没有已验证替代查询时说明范围限制。

结果按组聚合，从 `searchContractListItems` 提取入口。`filteredTotalElements` 是聚合组数，不能当合同数；`totalElements` 是各 tab 统计。续页固定 pageIndex=0、排序不变，原样传 searchAfter，直到终止页并核对组数；重复游标、数量漂移或窗口限制而无有效游标时不能宣称查全。后续页 paginationCapability 为空，不抹掉已有游标。

用户给了条件就不能退化为全租户搜索。名称、编号和综合关键词不互换。当前配方只确认名称、编号字段；未知筛选条件不能自行构造。完整搜索全选的业务执行目前不开放。

## 普通合同详情与回查

`GET /open-web/web/clm/contract/info`

| 参数 | 位置 | 必填 | 值或来源 |
|---|---|---|---|
| `contractId` | query | 是 | 已定位的普通合同内部 ID；不能填合同编号、聚合组 ID 或父组 ID。 |
| `tianyancha` | query | 是 | 固定 `false`。 |
| `generateContractName` | query | 是 | 固定 `false`。 |

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
- 普通指定合同测试要求 assemblyId 明确为 null 或 0；缺失或类型关系不明时停止该写入。groupId 与 assemblyId 不能互换；前后链为空也不能证明没有同组关联记录。
- 普通交接的角色规则见 [交接与结果核对](handover.md)。状态码 1 为已取消，不交接；已知普通准入范围为 2～17、另含 23，但状态落入范围不独自构成交接资格。草稿、协商等依赖未核实配置的分支不猜测放行，仍由服务端最终判定。
- 查询历史时分别显示交接前人员、接收人、操作人和时间，不用当前归属人补齐历史；空记录只表示接口未返回记录。时间单位或时区未知时保留原值。
- 不改用会更新交易方状态的聚合详情接口，也不通过猜参数、换身份来绕过详情异常。
