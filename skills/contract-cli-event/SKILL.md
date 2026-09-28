---
name: contract-cli-event
version: 1.0.1
description: "contract-cli 事件出口 IP 查询技能：用 app 身份分页查询 `/open-apis/event/v1/outbound_ip`。当用户要使用 `contract-cli event outbound-ip list` 查询事件出口 IP 白名单时触发。"
---

# contract-cli Event

CRITICAL — 开始前 MUST 先读取 [../contract-cli-shared/SKILL.md](../contract-cli-shared/SKILL.md)。

## 适用命令

- `contract-cli event outbound-ip list`

## 关键规则

- 当前仅支持 `--as app`
- 走 `GET /open-apis/event/v1/outbound_ip`
- 支持 `--page-size`、`--page-token`
- `--page-size` 必须在 `10` 到 `50` 之间；不传时由后端默认分页
- 不接受 `--input-file` / `--data`
- 完整参数、类型和约束：读 [references/outbound-ip-list-parameters.md](references/outbound-ip-list-parameters.md)

## 示例

```bash
contract-cli event outbound-ip list --profile contract --as app --page-size 10
contract-cli event outbound-ip list --profile contract --as app --page-size 10 --page-token next
```
