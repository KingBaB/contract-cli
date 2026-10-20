# 身份、权限与接收人

先读 [通道约定](transport.md)。只读当前身份，不接受用户粘贴的凭证。

## 当前操作者

`GET /open-web/web/clm/employee/info`

无业务 query 或 body 参数。

关键响应字段（路径相对完整响应）：

| 含义 | 路径 | 类型 | 可为空 |
|---|---|---|---|
| 员工 ID | `$.data.employeeId` | string | 否 |
| 租户 ID | `$.data.tenantId` | string | 否 |
| 名称 | `$.data.name` | string | 是 |
| 部门列表 | `$.data.departments` | array | 是 |

## 功能权限

`GET /open-web/web/clm/permission/function/list`

无业务 query 或 body 参数。

关键响应字段（路径相对完整响应）：

| 含义 | 路径 | 类型 | 可为空 |
|---|---|---|---|
| 权限列表 | `$.data` | array | 否 |
| 功能权限码 | `$.data[*].code` | array | 否 |

## 查找接收人

`GET /open-web/web/clm/employee/findEmployeesByKeywords`

| 参数 | 位置 | 必填 | 值或来源 |
|---|---|---|---|
| `keywords` | query | 是 | 用户输入的接收人姓名关键词。 |
| `includeInvalid` | query | 是 | 固定 `false`。 |

关键响应字段（路径相对完整响应）：

| 含义 | 路径 | 类型 | 可为空 |
|---|---|---|---|
| 员工列表 | `$.data` | array | 否 |
| 员工 ID | `$.data[*].employeeId` | array | 否 |
| 租户 ID | `$.data[*].tenantId` | array | 否 |
| 部门列表 | `$.data[*].departments` | array | 是 |
| 候选展示名 | `$.data[*].displayName` | array | 是 |
| 平台状态 | `$.data[*].status` | array | 是 |

## 接收人详情

`GET /open-web/web/clm/employee/getEmployeesByIds`

| 参数 | 位置 | 必填 | 值或来源 |
|---|---|---|---|
| `employeeIds` | query | 是 | 已选择的一个内部 employeeId，十进制字符串；不能填开放 user_id。 |
| `extraDepartmentInfo` | query | 是 | 固定 `true`。 |

关键响应字段（路径相对完整响应）：

| 含义 | 路径 | 类型 | 可为空 |
|---|---|---|---|
| 员工列表 | `$.data` | array | 否 |
| 员工 ID | `$.data[0].employeeId` | string | 是 |
| 租户 ID | `$.data[0].tenantId` | string | 是 |
| 平台状态 | `$.data[0].status` | integer | 是 |
| 是否失效 | `$.data[0].expired` | boolean | 是 |
| 是否锁定 | `$.data[0].locked` | boolean | 是 |
| 部门列表 | `$.data[0].departments` | array | 是 |
| 首部门 ID | `$.data[0].departments[0].departmentId` | string | 是 |

## 如何判定接收人和权限

- 功能权限码为 `manage_contract_handover_personal` 或 `manage_contract_handover_all`；两者都没有时不交接。有权限不代表某份合同一定允许交接，服务端还会校验角色和状态。
- 姓名查询可能重名，用部门等信息让用户选定一个人，不默认第一条。已有可靠内部 ID 时仍需确认本次正常候选和详情能对应同一人。
- `includeInvalid=false` 的候选与详情必须匹配同一 employeeId、tenantId；status 必须非空且一致。不能把某个固定数字（如 1）当作通用在职证明。
- 写前重读候选和详情。明确离职、删除、停用、`expired/locked=true`、租户不同或字段矛盾时停止。expired/locked 为空不是 false；不将缺失字段解释为已核实有效。
- 部门使用返回的第一项；无部门时保留原合同部门，不替接收人选择其他部门。无法消歧或必要有效性不能确认时不写入。
