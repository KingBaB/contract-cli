# Payment Plan Fields Reference

本页适用于：

- `contract-cli payment plan notify --input-file notify.json`
- `contract-cli payment plan search --input-file payment-plan-search.json`

字段口径按 CLM 后端 `PaymentLineUpsertDTO`、`PaymentPlanQueryDTO` 和 `PaymentPlanEsService` 校验整理；JSON 字段使用 snake_case。

## 同步付款记录

`payment plan notify` 路由：`POST /open-apis/contract/v1/payment/notify`

顶层字段：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `finance_number` | string | 是 | 财务侧付款单号，不能为空白 |
| `finance_amount` | number | 是 | 付款单金额 |
| `finance_currency` | string | 是 | 币种，使用 `CurrencyType.enName` |
| `finance_url` | string | 否 | 付款单链接 |
| `applicant_user_id` | integer | 二选一 | 申请人飞书 larkId |
| `applicant_employee_id` | string | 二选一 | 申请人员工 ID；和 `applicant_user_id` 至少传一个 |
| `payment_lines` | array | 是 | 付款行记录，不能为空 |

`payment_lines[]` 字段：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `source_id` | string | 是 | 付款行号，不能为空白 |
| `transaction_amount` | number | 是 | 付款行金额 |
| `transaction_status` | string | 是 | 交易状态枚举名称 |
| `transaction_time` | string | 是 | `yyyy-MM-dd HH:mm:ss` |
| `currency` | string | 是 | 币种 enName |
| `trading_party_id` | string | 二选一 | 交易方 ID |
| `trading_party_code` | string | 二选一 | 交易方编码；和 `trading_party_id` 至少传一个 |
| `trading_party_account` | object | 是 | 交易方账户 |
| `url` | string | 否 | 付款行外部链接；不传时后端按空字符串处理 |
| `relations` | array | 是 | 关联合同/付款计划，不能为空 |

`trading_party_account` 字段：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `account_number` | string | 是 | 账号，不能为空 |
| `account_type` | string | 否 | 账户类型 |

`relations[]` 字段：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `contract_id` | string | 是 | 合同 ID，不能为空 |
| `amount` | number | 是 | 本关联金额 |
| `payment_plan_uuid` | string | 否 | 付款计划 UUID |
| `contract_number` | string | 否 | 合同编号 |

`transaction_status` 支持：

- `IN_APPROVAL`
- `APPROVAL_COMPLETE`
- `APPROVAL_REVOKE`
- `APPROVAL_WITHDRAW`
- `PAYMENT_PROCESSING`
- `PAYMENT_SUCCESS`
- `PAYMENT_FAIL`
- `DISCARD`
- `IN_DRAFT`
- `DELETE`

示例：

```json
{
  "finance_number": "FIN-20260706-001",
  "finance_amount": 1000.5,
  "finance_currency": "CNY",
  "finance_url": "https://finance.example.com/payments/FIN-20260706-001",
  "applicant_user_id": 1152307041667121520,
  "payment_lines": [
    {
      "source_id": "line-001",
      "transaction_amount": 1000.5,
      "transaction_status": "PAYMENT_SUCCESS",
      "transaction_time": "2026-07-06 15:30:00",
      "currency": "CNY",
      "trading_party_code": "TP-001",
      "trading_party_account": {
        "account_number": "6222000000000000"
      },
      "relations": [
        {
          "contract_id": "7023646046559404327",
          "payment_plan_uuid": "payment-plan-uuid-001",
          "amount": 1000.5
        }
      ]
    }
  ]
}
```

## 搜索付款计划

`payment plan search` 路由：`POST /open-apis/contract/v1/payments/search`

顶层字段：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `page_size` | integer | 否 | 默认 `20`；后端校验范围 `1..200` |
| `page_index` | integer | 否 | 默认 `0`；从 0 开始 |
| `sort_field` | string | 否 | 必须是支持的 `field_code` |
| `order` | string | 否 | `asc` 或 `desc` |
| `user_id` | integer | 否 | 飞书 larkId，用于权限过滤 |
| `submitter_user_ids` | integer[] | 否 | 提交人飞书 larkId 列表 |
| `department_ids` | integer[] | 否 | 合同所有人部门 larkId 列表 |
| `payment_reminder_type_codes` | integer[] | 否 | `0` 需付款，`1` 无需付款 |
| `must_conditions` | object[] | 否 | 必须满足的条件 |
| `should_conditions` | object[] | 否 | 至少满足一个的条件 |

条件对象字段：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `field_code` | string | 是 | 搜索字段编码 |
| `operator_type_code` | integer | 是 | 操作符编码 |
| `field_values` | string[] | 见说明 | 除 `operator_type_code = 18` 外必须非空且元素非空 |

当前代码只支持这些查询操作符：

| code | 名称 | 适用字段类型 |
| --- | --- | --- |
| `0` | 等于 | keyword、text_keyword、number |
| `2` | 大于 | number/date |
| `3` | 小于 | number/date |
| `4` | 大于等于 | number/date |
| `5` | 小于等于 | number/date |
| `6` | 包含 | keyword、number |
| `8` | 属于 | keyword、number |
| `18` | 不为空 | keyword |
| `19` | 模糊匹配 | text_keyword |
| `20` | 区间 | number/date，`field_values` 必须正好 2 个 |

支持的 `field_code`：

| field_code | 输入说明 |
| --- | --- |
| `contractId` | 数字字符串 |
| `contractStatus` | 状态 code |
| `contractNumber` | 合同编号 |
| `contractName` | 合同名称 |
| `contractRemark` | 合同备注 |
| `contractLegalCode` | 我方主体编码 |
| `contractDepartmentId` | 部门 ID，通常由 `department_ids` 自动追加 |
| `contractSubmitUserId` | 提交人 ID，通常由 `submitter_user_ids` 自动追加 |
| `contractCategoryCode` | 合同分类编码 |
| `contractSignDate` | `yyyy-MM-dd` |
| `contractPropertyCode` | 合同属性 code |
| `contractDurationDate` | `yyyy-MM-dd`；区间用 `20`，大于/小于会分别映射到开始/结束日期 |
| `paymentTradingCode` | 交易方编码 |
| `paymentDateExpect` | `yyyy-MM-dd` |
| `currency` | 币种 enName，如 `CNY` |
| `paymentPlanRemark` | 付款计划备注 |
| `paymentUnpaidAmount` | 数字字符串 |
| `paymentAmount` | 数字字符串 |
| `paymentUuids` | 付款计划 UUID |
| `paymentPlanStatus` | 付款计划状态 code |

注意：

- 不要直接传 `contractDurationDateStart` 或 `contractDurationDateEnd`；当前代码禁止这两个字段作为输入。
- `should_conditions` 单独使用时当前后端校验路径有空指针风险；优先使用 `must_conditions`，需要 `should_conditions` 时同时提供合法的 `must_conditions`。

示例：

```json
{
  "page_size": 20,
  "page_index": 0,
  "order": "desc",
  "must_conditions": [
    {
      "field_code": "paymentUuids",
      "operator_type_code": 8,
      "field_values": ["payment-plan-uuid-001"]
    },
    {
      "field_code": "currency",
      "operator_type_code": 0,
      "field_values": ["CNY"]
    }
  ]
}
```
