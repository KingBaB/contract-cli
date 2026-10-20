---
name: contract-cli-mdm-legal
version: 1.0.1
description: "contract-cli 法人实体主数据技能：列出法人实体候选列表、按 ID 获取详情、按编码查询，或用 app 身份创建、更新法人实体。当用户要使用 `contract-cli mdm legal ...` 操作合同域法人实体数据时触发。"
---

# contract-cli MDM Legal

CRITICAL — 开始前 MUST 先读取 [../contract-cli-shared/SKILL.md](../contract-cli-shared/SKILL.md)。

## 适用命令

- `contract-cli mdm legal list`
- `contract-cli mdm legal get <legal-entity-id>`
- `contract-cli mdm legal get --code <code>`
- `contract-cli mdm legal create`
- `contract-cli mdm legal update <legal-entity-id>`

## 快速决策

- 已知法人实体 ID：直接用 `mdm legal get`
- 已知法人实体编码：用 `mdm legal get --as app --code <code>`
- 只知道名称、需要候选列表：用 `mdm legal list`
- 想创建或更新法人实体：用 `mdm legal create|update --as app --user-id <operator-user-id> --input-file ...`
- 用户是想查字段配置：切到 [../contract-cli-mdm-fields/SKILL.md](../contract-cli-mdm-fields/SKILL.md)

## 关键规则

- `mdm legal list` 支持 `--name`、`--page-size`、`--page-token`
- `mdm legal list --as user` 走 `/open-apis/contract/v1/mcp/legal_entities`
- `mdm legal list --as app` 走 `/open-apis/mdm/v1/legal_entities/list_all`
- `mdm legal get --as user` 走 `/open-apis/contract/v1/mcp/legal_entities/{legal_entity_id}`
- `mdm legal get --as app` 走 `/open-apis/mdm/v1/legal_entities/{legal_entity_id}`，并额外透传 query `legal_entity_id`
- `mdm legal get --code` 当前仅支持 `--as app`
- `mdm legal create/update` 当前仅支持 `--as app`
- `mdm legal create` 走 `POST /open-apis/mdm/v1/legal_entities`
- `mdm legal update` 走 `PUT /open-apis/mdm/v1/legal_entities/{legal_entity_id}`
- `mdm legal get --code` 走 `GET /open-apis/mdm/v1/legal_entities`，`--code` 映射到 query `legalEntity`
- 不暴露 `--operator`
- `mdm legal create/update` 必须传 `--user-id`；`--user-id-type` 不传时默认 `user_id`
- `mdm legal create` 请求体不要包含后端生成的 `legalEntity` / `legal_entity` 编码
- `mdm legal update` 请求体必须包含后端返回的 `id` 和 `legalEntity` 编码；字段名使用 camelCase `legalEntity`，不要用 `legal_entity`
- 创建/更新请求体除上述规则外仍直接透传 JSON，其他字段是否必填受后台动态字段配置影响
- 推荐阅读顺序是：
  - 先读 [references/entity-query-guide.md](references/entity-query-guide.md) 选查询场景
  - 再读 [references/entity-query-parameters.md](references/entity-query-parameters.md) 查请求参数映射
  - 最后读 [references/commands.md](references/commands.md) 抄命令示例
- P3 命令按需读取完整参数文档：
  - 创建：[references/legal-create-parameters.md](references/legal-create-parameters.md)
  - 更新：[references/legal-update-parameters.md](references/legal-update-parameters.md)
  - 按编码查询：[references/legal-get-by-code-parameters.md](references/legal-get-by-code-parameters.md)

## 实现来源

- [internal/cli/vendor_command.go](../../internal/cli/vendor_command.go)
- [internal/openplatform/entity/service.go](../../internal/openplatform/entity/service.go)
- [references/entity-query-guide.md](references/entity-query-guide.md)
- [references/entity-query-parameters.md](references/entity-query-parameters.md)
- [references/commands.md](references/commands.md)
- [references/legal-create-parameters.md](references/legal-create-parameters.md)
- [references/legal-update-parameters.md](references/legal-update-parameters.md)
- [references/legal-get-by-code-parameters.md](references/legal-get-by-code-parameters.md)

## 操作建议

- 只知道主体名称时，先 `mdm legal list` 拿候选，再 `mdm legal get` 查详情
- 只知道主体编码时，直接 `mdm legal get --as app --code <code>`
- 创建合同前如果只是要选我方主体，优先记住法人实体 id
- 创建或更新法人实体前，先用 `mdm fields list --biz-line legal_entity` 确认动态字段
- 想确认写接口字段定义时，切到 `mdm fields list`

## 不要这样做

- 不要把 `mdm legal list/get` 当成字段配置查询
- 不要在 create 请求体里传后端生成的 `legalEntity` / `legal_entity` 编码
- 不要在 update 请求体里把编码字段写成 `legal_entity`
- 不要在 create/update 请求体里写入 token、密钥或无关个人敏感信息
