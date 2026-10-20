# 交接请求与结果核对

仅在 [主流程](../SKILL.md) 的环境、范围和确认条件满足后调用。

## 执行交接

`POST /open-web/web/clm/permission/handover`

| 参数 | 位置 | 必填 | 值或来源 |
|---|---|---|---|
| `manual` | body | 是 | 固定 `true`。 |
| `reason` | body | 是 | 固定 `""`。 |
| `targetEmployeeId` | body | 是 | 已确认接收人的内部 employeeId，写前重新读取。 |
| `targetTenantId` | body | 是 | 接收人详情 tenantId，必须与当前操作者 tenantId 相同。 |
| `resourceList` | body | 是 | 1～100 个去重入口；字段详见下文，不提交动态搜索条件。 |

关键响应字段（路径相对完整响应）：

| 含义 | 路径 | 类型 | 可为空 |
|---|---|---|---|
| 本次返回授权 | `$.data.authorizations` | array | 是 |
| 授权 ID | `$.data.authorizations[*].authorizationId` | array | 是 |
| 业务对象 ID | `$.data.authorizations[*].businessId` | array | 是 |
| 业务类型码 | `$.data.authorizations[*].businessTypeCode` | array | 是 |
| 接收人 ID | `$.data.authorizations[*].authorizedEmployeeId` | array | 是 |
| 接收租户 ID | `$.data.authorizations[*].authorizedTenantId` | array | 是 |
| 操作人 ID | `$.data.authorizations[*].approvedEmployeeId` | array | 是 |
| 创建时间 | `$.data.authorizations[*].createTime` | array | 是 |

## resourceList 与请求示例

| 字段 | 类型 | 来源 |
|---|---|---|
| businessId | 十进制字符串 | 当前详情中的合同 ID，合同组场景为父组 ID |
| businessTypeCode | 整数 | 当前入口详情的真实类型码，不把所有对象填成 0 |
| parentId | 十进制字符串 | 当前入口详情的 groupId，不是 assemblyId，也不是搜索聚合组 ID |

以下只展示 JSON 结构，所有示例 ID 必须替换为本次实际读取且已确认的数据；大整数保持字符串：

```json
{"manual":true,"reason":"","targetEmployeeId":"9223372036854775804","targetTenantId":"9223372036854775803","resourceList":[{"businessId":"9223372036854775806","businessTypeCode":0,"parentId":"9223372036854775805"}]}
```

不发送 `handoverByAdmin`、operatorId、onlyOwner、onlySubmitter 等身份或角色选择参数；也不把验收标记、搜索条件写进 body。100 是本 Skill 的单批入口策略，不是服务端容量承诺。组不拆成成员逐个写，重叠影响不跨批。

## 普通合同的预期变化

| 交接前 | 交接后 |
|---|---|
| 当前归属人与申请人相同 | 两者均变更为接收人 |
| 当前归属人与申请人不同 | 归属人变更，申请人保持原值 |

此规则适用于普通归属人交接路径，不是权限判定。管理员第三方交接分支由服务端按当前权限与角色决定，不提交管理员标记。有全部权限不代表一定进入第三方分支；仅申请人操作的未验证角色分支不得套用归属人规则执行。仅要求修改申请人或需求人时不调用交接。

普通合同成功时部门按接收人第一部门更新，无部门则保留；需求人不应被交接修改。这些字段有现成结果可另行说明，不把其缺证当作指定普通核心测试的角色交接失败。关联记录、草稿、流程、旧授权及通知由后端既有规则处理，不能承诺全部审批待办迁移、原人丧失所有查看权或所有旧提醒停止。

## 逐份判定

保存写前角色和历史记录。调用一次后，以本次授权中的 businessId、businessTypeCode、authorizedEmployeeId、authorizedTenantId，与写后新增交接记录及角色前后变化核对。回读接口见 [合同详情](contracts.md) 或 [合同组](assemblies.md)。历史记录中的 preAuthorizedEmployee、authorizedEmployee、approvedEmployee 分别是交接前人员、接收人、操作人，不能互换。

| 输出状态 | 判定 |
|---|---|
| 成功 | 本次授权／新增交接记录能与对象和接收人对应，角色回读符合预期 |
| 无新增变化 | 写前已经是目标状态，未证明本次新增交接；不计为本次成功 |
| 未执行 | 确认前排除、取消，或后续批次没有发送 |
| 跳过 | 有可靠证据证明服务端未执行，说明已知原因 |
| 失败 | 接口或可靠回查明确该对象失败 |
| 待核实 | 已发送但超时、缺响应、缺记录、字段矛盾或回查不可用 |

只有整体 authorizations，没有完整逐项失败／跳过清单。HTTP 200、success=true、空授权或数量不足不能判全部成功；403/404 或列表消失也不能独自判成功或失败。已有目标角色或旧历史不能当本次成功证据。

请求不是整批原子事务，超时或失败可能已经部分完成。出现不确定结果时停止后续批次，只读回查；不自动重试、不切换 CLI/MCP 重放、不交回原人当回滚。即使后来核实成功也不自动续发被停止批次。用户明确续办时，重新预检剩余范围；未决对象及与其有共同影响的对象不得再次提交。

指定普通核心测试只报告指定入口的结果；关联影响未核实应单独列出，不声称全量副作用已通过。
