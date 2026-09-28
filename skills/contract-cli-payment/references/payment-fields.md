# Payment Fields Reference

本页适用于：

- `contract-cli payment create --contract <contract-id> --input-file payment.json`
- `contract-cli payment update <payment-id> --contract <contract-id> --input-file payment-update.json`

字段口径按 CLM 后端 `PaymentCreateDTO`、`PaymentUpdateDTO` 和 `PaymentService` 校验整理；JSON 字段使用 snake_case。

## 通用枚举

- `payment_status_code`：`0` 暂存，`1` 审批中，`2` 付款取消，`3` 付款成功，`4` 付款失败。当前后端枚举没有 `9`。
- `finance_system_code`：`0` 金蝶，`1` 飞书费控，`2` 理想汽车，`3` 合同中心对公。
- `currency_code`：使用 `CurrencyType.enName`，例如 `CNY`、`USD`、`JPY`、`EUR`、`GBP`、`SGD`、`THB`、`HKD`、`AUD`。
- 金额字段是 JSON number / decimal，不要写成带币种的字符串。
- 时间字段 `apply_time`、`complete_time` 是数字字符串，按毫秒时间戳处理。
- `apply_user_id` 是飞书 larkId 的数字字符串；后端只在 `NumberUtil.isLong` 成立时查找并写入申请人。

## 创建付款申请

`payment create` 路由：`POST /open-apis/contract/v1/contracts/{contract_id}/payments`

必填字段：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `finance_system_code` | integer | 财务系统枚举 |
| `payment_status_code` | integer | 付款状态枚举 |
| `finance_number` | string | 财务侧付款单号 |
| `currency_code` | string | 币种 enName |

可选字段：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `source_id` | string | 外部来源 ID |
| `finance_url` | string | 财务侧链接 |
| `apply_amount` | number | 申请金额 |
| `succeed_amount` | number | 成功金额 |
| `fail_amount` | number | 失败金额 |
| `has_invoice` | boolean | 是否有发票 |
| `apply_user_id` | string | 飞书 larkId 的数字字符串 |
| `apply_time` | string | 毫秒时间戳字符串；不传默认当前时间 |
| `complete_time` | string | 毫秒时间戳字符串 |
| `extra_info` | string | 扩展信息 |

示例：

```json
{
  "finance_system_code": 0,
  "payment_status_code": 1,
  "finance_number": "FIN-20260706-001",
  "currency_code": "CNY",
  "apply_amount": 1000.5,
  "succeed_amount": 0,
  "fail_amount": 0,
  "has_invoice": true,
  "apply_user_id": "1152307041667121520",
  "apply_time": "1783315200000",
  "source_id": "payment-source-001"
}
```

## 更新付款申请

`payment update` 路由：`PATCH /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}`

后端更新入口当前要求以下字段非空：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `finance_system_code` | integer | 财务系统枚举 |
| `payment_status_code` | integer | 付款状态枚举 |
| `finance_number` | string | 财务侧付款单号 |
| `finance_url` | string | 财务侧链接 |

其余可选字段与创建付款申请一致。可选字段只在传入且格式可用时更新；例如 `currency_code` 无法匹配 `CurrencyType.enName` 时不会写入，`apply_user_id` 不是数字字符串时不会更新申请人。

示例：

```json
{
  "finance_system_code": 0,
  "payment_status_code": 3,
  "finance_number": "FIN-20260706-001",
  "finance_url": "https://finance.example.com/payments/FIN-20260706-001",
  "currency_code": "CNY",
  "apply_amount": 1000.5,
  "succeed_amount": 1000.5,
  "fail_amount": 0,
  "complete_time": "1783318800000"
}
```
