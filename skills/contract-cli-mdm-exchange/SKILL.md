---
name: contract-cli-mdm-exchange
version: 1.0.1
description: "contract-cli 固定汇率技能：用 app 身份查询或更新 `/open-apis/mdm/v1/fixed_exchange_rate`。当用户要使用 `contract-cli mdm fixed-exchange-rate get|update` 查询、维护固定汇率时触发。"
---

# contract-cli MDM Exchange

CRITICAL — 开始前 MUST 先读取 [../contract-cli-shared/SKILL.md](../contract-cli-shared/SKILL.md)。

## 适用命令

- `contract-cli mdm fixed-exchange-rate get`
- `contract-cli mdm fixed-exchange-rate update`

## 快速决策

- 只查固定汇率：用 `get --source-currency --target-currency --effective-date`
- 新增或更新固定汇率：用 `update --input-file <json>`

## 关键规则

- 两个命令当前都仅支持 `--as app`
- `get` 走 `GET /open-apis/mdm/v1/fixed_exchange_rate`
- `update` 走 `PUT /open-apis/mdm/v1/fixed_exchange_rate`
- `get` 的 CLI flag `--effective-date` 会映射到底层 query `date`；`update` 请求体仍使用字段 `effective_date`
- `get` 不接受 `--input-file` / `--data`
- `update` 必须通过 `--input-file` 或 `--data` 传 JSON

## 参数文档

- 查询固定汇率：[references/fixed-exchange-rate-get-parameters.md](references/fixed-exchange-rate-get-parameters.md)
- 更新固定汇率：[references/fixed-exchange-rate-update-parameters.md](references/fixed-exchange-rate-update-parameters.md)

## 示例

```bash
contract-cli mdm fixed-exchange-rate get --profile contract --as app --source-currency CNY --target-currency USD --effective-date 2026-06-01
contract-cli mdm fixed-exchange-rate update --profile contract --as app --input-file fixed-exchange-rate.json
```

`fixed-exchange-rate.json` 最小形态：

```json
{
  "source_currency": "CNY",
  "target_currency": "USD",
  "exchange_rate": "7.1",
  "effective_date": "2026-06-01",
  "status": 1
}
```
