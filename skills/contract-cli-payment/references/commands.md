# Payment Commands Reference

## 付款申请

```bash
contract-cli payment create --contract 7023646046559404327 --profile contract --as app --input-file payment.json
contract-cli payment update payment_123 --contract 7023646046559404327 --profile contract --as app --input-file payment-update.json
contract-cli payment get payment_123 --contract 7023646046559404327 --profile contract --as app
contract-cli payment list --contract 7023646046559404327 --profile contract --as app --page-size 10 --page-token next
```

接口路径：

- `create`：`POST /open-apis/contract/v1/contracts/{contract_id}/payments`
- `update`：`PATCH /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}`
- `get`：`GET /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}`
- `list`：`GET /open-apis/contract/v1/contracts/{contract_id}/payments`

请求体字段：看 [payment-fields.md](payment-fields.md)

## 付款计划

```bash
contract-cli payment plan notify --profile contract --as app --input-file notify.json
contract-cli payment plan search --profile contract --as app --input-file payment-plan-search.json
```

接口路径：

- `notify`：`POST /open-apis/contract/v1/payment/notify`
- `search`：`POST /open-apis/contract/v1/payments/search`

请求体字段：看 [payment-plan-fields.md](payment-plan-fields.md)

## 付款记录

```bash
contract-cli payment record create --contract 7023646046559404327 --payment payment_123 --profile contract --as app --input-file payment-record.json
contract-cli payment record update record_123 --contract 7023646046559404327 --payment payment_123 --profile contract --as app --input-file payment-record-update.json
contract-cli payment record get record_123 --contract 7023646046559404327 --payment payment_123 --profile contract --as app
contract-cli payment record list --plan payment_plan_uuid_123 --profile contract --as app
```

接口路径：

- `record create`：`POST /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records`
- `record update`：`PATCH /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records/{payment_record_id}`
- `record get`：`GET /open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records/{payment_record_id}`
- `record list`：`GET /open-apis/contract/v1/contracts/payments/{payment_plan_uuid}/payment_records`

请求体字段：看 [payment-record-fields.md](payment-record-fields.md)
