# 补齐已选合同的父组 ID

合同搜索仍由 [合同搜索 Skill](../../contract-cli-contract-search/SKILL.md) 完成。本步骤只补充已选对象的组关系，不重选合同、不重新解释原筛选条件，也不扩大用户选择的范围。

## 先使用已有映射

先核对与交接相同的环境、租户和 user 授权，并保留候选内部 ID、真实业务类型、编号／名称及字段来源。字段名相似或编号相同不证明内部 ID 等价。

- 已确认为父组的 type=10/17 候选：使用其真实内部入口 ID 读取父组详情，不先调普通 contract/info。父组行的 assemblyId 可能为负数占位值，不能覆盖其真实入口 ID。
- 已确认成员的候选：实际响应中已有语义明确的正十进制 assemblyId／assembly_id 时，直接读取该父组及完整成员，不要求成员先经过普通 contract/info。
- 已返回 v2.1 Presentation 时，先核对 presentation.schema.id=clm.contract-search-presentation、version=1；通过选中行 row_ref 读取 data.presentation.interaction.entities 中同一实体。核对其 contract_id 或 child_contract_id 与候选内部 ID、business_type_code 一致，再读取 assembly_id。role=PARENT 只表示展示父行，不能据此把 type=11 改成父组类型；无此扩展时不假定它存在，不自行切换响应版本。

type=11 成员缺少父组 ID，或已有明确组上下文却无法从候选取得组关系时，使用下述只读补充查询。不能把普通详情报错本身当成组成员证明。普通合同没有组线索时仍按 [普通详情](contracts.md#普通合同详情与回查) 处理，不为所有合同额外发起本查询。

## 只读补充查询

`POST /open-web/web/clm/es/search/contractList`

通过现有通用 Open Web 通道调用，不绕过网关。参数约定见 [通道说明](transport.md)。只在已有搜索不能提供所需父组映射时使用；搜索服务的业务拒绝、身份／环境不一致不能通过本步骤绕过。

| 参数 | 位置 | 类型 | 本流程要求 | 值、来源与限制 |
|---|---|---|---|---|
| pageIndex | body | integer | 必须传 | 固定 0；使用游标续查时仍为 0 |
| pageSize | body | integer | 必须传 | 固定 100 |
| searchTabCode | body | integer | 必须传 | 沿用本轮已确认的全部合同 0／我的合同 1；上下文不明时先明确范围，不自行扩大 |
| conditionUnits | body | array&lt;object&gt; | 必须传 | 仅一个非空编号或名称条件，子字段见下表；优先已选候选的完整编号，无编号时使用该候选完整名称 |
| allConditionUnits | body | boolean | 必须传 | 固定 false，不切到综合关键词模式 |
| sortType | body | string | 必须传 | 固定 SUBMITTED_TIME |
| order | body | string | 必须传 | 固定 DESC |
| withStatistics | body | boolean | 必须传 | 固定 false |
| searchAfter | body | array（不透明游标） | 续页必须传 | 首轮省略；原样保留前页数组及各项类型，不替换成搜索 Skill 的 page_token |

| conditionUnits 子字段 | 类型 | 本流程要求 | 值、来源与限制 |
|---|---|---|---|
| searchField | string | 必须传 | 编号用 CONTRACT_NUMBER_FIELD；仅无编号时名称用 CONTRACT_NAME_FIELD |
| searchValue | string | 必须传 | 已选候选的完整编号或名称，不使用聚合组 ID，不填空字符串 |
| unionType | string | 必须传 | 固定 SHOULD；不增加其他条件或 filterUnits |

以下编号仅展示请求形态，执行前替换为已选候选值，并沿用已确认的 searchTabCode；没有可用编号／名称时说明缺少定位条件，不做无条件枚举：

```json
{"pageIndex":0,"pageSize":100,"searchTabCode":0,"conditionUnits":[{"searchField":"CONTRACT_NUMBER_FIELD","searchValue":"CG-EXAMPLE-1017","unionType":"SHOULD"}],"allConditionUnits":false,"sortType":"SUBMITTED_TIME","order":"DESC","withStatistics":false}
```

仅在 HTTP／通道成功且下游 code=200、success=true 后读取响应；POST 在这里是查询，不是交接写入。

| 响应字段 | 类型 | 用途 |
|---|---|---|
| $.data.content[*].searchContractListItems[*].contractId | string（内部 ID） | 与已选候选内部 ID 精确绑定 |
| 同一条目的 contractNumber／contractName | string | 核对选中对象的编号／名称；模糊命中不形成新选择 |
| 同一条目的 businessTypeCode | integer | 核对实际成员／父组类型，不用其他条目的类型拼接 |
| 同一条目的 assemblyId | string，可为空或占位值 | 成员的真实正 ID 才能作为父组 ID；父组自身的占位值不使用 |
| $.data.content[*].groupId | string，可为空 | 搜索聚合组标识，不作为 contractAssemblyId 或写入 parentId |
| $.data.searchAfter | array，可为空 | 尚未找到选中对象时按原游标续查 |
| $.data.filteredTotalElements | integer，可为空 | 聚合组计数，不能当作成员总数 |

## 绑定后再读取父组

1. 在 searchContractListItems 中按已选内部 ID 精确匹配，并核对已有类型、编号／名称。无关条目不纳入本次清单；不能只凭相同编号／名称、第一条结果或同一外层 groupId 选父组。
2. 对匹配的 type=11 成员，取该条目正十进制 assemblyId；type=0 且实际有正 assemblyId 也按组成员处理。成员返回的 0、null、负值、缺失或冲突不能推导父组 ID，不能改传 businessTypeCode 尝试交接。
3. 对已核实的 type=10/17 父组条目，取该条目的真实 contractId 作为 contractAssemblyId，不取该行的占位 assemblyId。同一内部 ID 出现矛盾映射时停止；搜索结果不直接生成 resourceList。
4. 按 [合同组](assemblies.md) 读取父组详情，要求 data.contractAssemblyId 与定位 ID 一致，类型确为 10/17。再完整读取 childrenInfo，确认原选成员 contractId 及其 assemblyId 属于该父组；不要求该成员调用普通详情。父组或成员不匹配时停止，不能把另一父组纳入用户选择。
5. 匹配关系确认后才说明整组影响，按原规则取得用户对父组及全部成员范围的确认。搜索元数据不能代替父组／成员的当前角色、状态、权限和写后回查。type=17 仍须按自身当前权限与条件核实，不能沿用其他类型的成功结论。

本查询只为选中对象补齐映射，找到后即可通过父组和成员接口核实，不声称查全原搜索结果。未命中时保留条件、页签和排序，按真实 searchAfter 续查；空末页、无有效续页游标、重复游标或窗口限制仍未找到则报告未定位。禁止改成无条件搜索、请求用户手抄内部 ID，或用搜索聚合组 ID 猜测父组。接口不可路由、权限拒绝或关键字段缺失时报告具体缺口，不发起交接写入。
