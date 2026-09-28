---
name: contract-cli-payment
version: 1.0.0
description: "contract-cli 付款命令技能：支持 app 身份下创建/更新/查看/查询付款申请、同步/搜索付款计划、创建/更新/查看/按付款计划查询付款记录。当用户要使用 `contract-cli payment ...` 操作付款能力时触发。"
---

# contract-cli Payment

CRITICAL — 开始前 MUST 先读取 [../contract-cli-shared/SKILL.md](../contract-cli-shared/SKILL.md)。

## 适用命令

- `contract-cli payment create`
- `contract-cli payment update <payment-id>`
- `contract-cli payment get <payment-id>`
- `contract-cli payment list`
- `contract-cli payment plan notify`
- `contract-cli payment plan search`
- `contract-cli payment record create`
- `contract-cli payment record update <payment-record-id>`
- `contract-cli payment record get <payment-record-id>`
- `contract-cli payment record list`

## 快速决策

- 创建付款申请：`payment create --contract <contract-id> --input-file payment.json`
- 更新付款信息：`payment update <payment-id> --contract <contract-id> --input-file payment-update.json`
- 查看付款信息：`payment get <payment-id> --contract <contract-id>`
- 查询付款申请列表：`payment list --contract <contract-id>`
- 同步付款记录：`payment plan notify --input-file notify.json`
- 搜索付款计划：`payment plan search --input-file payment-plan-search.json`
- 创建付款记录：`payment record create --contract <contract-id> --payment <payment-id> --input-file payment-record.json`
- 更新付款记录：`payment record update <payment-record-id> --contract <contract-id> --payment <payment-id> --input-file payment-record-update.json`
- 查询付款记录详情：`payment record get <payment-record-id> --contract <contract-id> --payment <payment-id>`
- 按付款计划查询付款记录：`payment record list --plan <payment-plan-uuid>`

## 字段文档导航

- 付款申请创建/更新请求体：读 [references/payment-fields.md](references/payment-fields.md)
- 付款计划同步/搜索请求体：读 [references/payment-plan-fields.md](references/payment-plan-fields.md)
- 付款记录创建/更新请求体：读 [references/payment-record-fields.md](references/payment-record-fields.md)
- 付款命令和路由总览：读 [references/commands.md](references/commands.md)

## 关键规则

- 当前 `payment *` 全部仅支持 `--as app`。
- 命令参数采用“主操作对象 ID 用位置参数，父资源/上下文 ID 用 flag”的方式。
- `contract_id` 作为父资源时使用 `--contract <contract-id>`。
- `payment_id` 作为父资源时使用 `--payment <payment-id>`。
- `payment_plan_uuid` 使用 `--plan <payment-plan-uuid>`。
- POST/PATCH 命令必须传 `--input-file` 或 `--data`，且二者互斥。
- GET/list 命令不接受 `--input-file` / `--data`。
- `payment list` 支持 `--page-size` 和 `--page-token`。
- `--user-id-type` / `--user-id` 是开放平台通用 query 参数：
  - `--user-id-type` 不传时默认拼接 `user_id_type=user_id`
  - 显式传 `--user-id-type <type>` 时覆盖默认值
  - `--user-id` 传了就原样拼到底层接口，不传就不带

## 路由

- `payment create`：`POST /open-apis/contract/v1/contracts/{contract_id}/payments`
- `payment update`：`PATCH /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}`
- `payment get`：`GET /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}`
- `payment list`：`GET /open-apis/contract/v1/contracts/{contract_id}/payments`
- `payment plan notify`：`POST /open-apis/contract/v1/payment/notify`
- `payment plan search`：`POST /open-apis/contract/v1/payments/search`
- `payment record create`：`POST /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records`
- `payment record update`：`PATCH /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records/{payment_record_id}`
- `payment record get`：`GET /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records/{payment_record_id}`
- `payment record list`：`GET /open-apis/contract/v1/contracts/payments/{payment_plan_uuid}/payment_records`

## 实现来源

- [internal/cli/payment_command.go](../../internal/cli/payment_command.go)
- [internal/openplatform/payment/service.go](../../internal/openplatform/payment/service.go)
- [references/commands.md](references/commands.md)
- [references/payment-fields.md](references/payment-fields.md)
- [references/payment-plan-fields.md](references/payment-plan-fields.md)
- [references/payment-record-fields.md](references/payment-record-fields.md)

## 操作建议

- 先确认 profile 已完成 app 登录：`contract-cli auth login --profile contract --as app`
- 复杂请求体优先用 `--input-file`
- 需要脚本消费时加 `--output json`
- 需要对照后端原始 envelope 时加 `--raw`
- 查询类命令缺少父资源 ID 时优先补 `--contract`、`--payment` 或 `--plan`

## 不要这样做

- 不要对 `payment *` 传 `--as user`
- 不要把付款命令写到 `contract` 子命令下，例如不要写 `contract payment create`
- 不要把 JSON 请求体放进 `--file`，JSON 请求体始终用 `--input-file`
- 不要给 GET/list 命令传 `--input-file` 或 `--data`
