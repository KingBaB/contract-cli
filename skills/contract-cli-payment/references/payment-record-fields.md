# Payment Record Fields Reference

本页适用于：

- `contract-cli payment record create --contract <contract-id> --payment <payment-id> --input-file payment-record.json`
- `contract-cli payment record update <payment-record-id> --contract <contract-id> --payment <payment-id> --input-file payment-record-update.json`

字段口径按 CLM 后端 `PaymentRecordCreateDTO`、`PaymentRecordUpdateDTO`、`PaymentRecordService` 和 swagger 生成结果整理；JSON 字段使用 snake_case。

## 关键规则

- `contract_id`、`payment_id`、`payment_record_id` 在路径中都必须是数字字符串。
- 金额字段是 JSON number / decimal。
- `transaction_time` 使用 13 位毫秒时间戳字符串；创建时后端明确校验长度为 `13`。
- `currency_code` 使用 `CurrencyType.enName`。创建时不传或无效会默认 `CNY`；更新时无效值会被忽略。
- `operator_user_id` 是飞书 larkId，JSON integer；后端会用它查员工。
- 当前代码和 swagger 暴露的部门字段是 `department_lark_id`，JSON integer。飞书文档如写 `department_id`，与当前代码不一致。
- 创建时 `extra_info` 最长 `10000` 字符；更新时只在非空字符串时写入。

## 创建付款记录

`payment record create` 路由：`POST /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records`

字段：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `source_id` | string | 否 | 付款明细原始 ID |
| `related_id` | string | 否 | 关联 ID |
| `transaction_amount` | number | 否 | 交易金额；业务上建议传 |
| `succeed_amount` | number | 否 | 付款成功金额；业务上建议传 |
| `fail_amount` | number | 否 | 付款失败金额；业务上建议传 |
| `currency_code` | string | 否 | 币种 enName |
| `operator_user_id` | integer | 否 | 操作人飞书 larkId |
| `department_lark_id` | integer | 否 | 部门飞书 larkId |
| `transaction_time` | string | 否 | 13 位毫秒时间戳字符串 |
| `extra_info` | string | 否 | 扩展信息，最长 `10000` 字符 |

示例：

```json
{
  "source_id": "record-source-001",
  "related_id": "finance-line-001",
  "transaction_amount": 1000.5,
  "succeed_amount": 1000.5,
  "fail_amount": 0,
  "currency_code": "CNY",
  "operator_user_id": 1152307041667121520,
  "department_lark_id": 7263646046559404327,
  "transaction_time": "1783318800000",
  "extra_info": "{\"source_budget_record_id\":\"BR-001\"}"
}
```

## 更新付款记录

`payment record update` 路由：`PATCH /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records/{payment_record_id}`

请求体字段与创建付款记录一致，全部按“传了才更新”的方式处理：

- 字符串字段：`source_id`、`related_id`、`currency_code`、`transaction_time`、`extra_info` 只有非空字符串才写入。
- 金额字段：`transaction_amount`、`succeed_amount`、`fail_amount` 只要非 null 就写入。
- ID 字段：`operator_user_id`、`department_lark_id` 只要非 null 就查对应员工/部门；查不到会报错。

示例：

```json
{
  "transaction_amount": 1200,
  "succeed_amount": 1200,
  "fail_amount": 0,
  "currency_code": "CNY",
  "transaction_time": "1783322400000",
  "extra_info": "{\"source_budget_record_id\":\"BR-002\"}"
}
```
