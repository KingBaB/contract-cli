# Entity Commands Reference

```bash
contract-cli mdm legal list --profile contract --name "上海主体"
contract-cli mdm legal list --profile contract --name "上海主体" --page-size 20 --page-token next
contract-cli mdm legal list --profile contract --as app --name "主体A" --page-size 20 --user-id-type employee_id
contract-cli mdm legal get 7023646046559404327 --profile contract
contract-cli mdm legal get 7003410079584092448 --profile contract --as app --user-id-type employee_id
contract-cli mdm legal get --profile contract --as app --code L0001 --page-size 10
contract-cli mdm legal create --profile contract --as app --user-id <operator-user-id> --input-file legal-create.json
contract-cli mdm legal update 7003410079584092448 --profile contract --as app --user-id <operator-user-id> --input-file legal-update.json
```

说明：

- `mdm legal list` 会按身份路由：
  - `user` -> `contract/v1/mcp/legal_entities`
  - `app` -> `mdm/v1/legal_entities/list_all`
- `mdm legal get` 也会按身份路由：
  - `user` -> `contract/v1/mcp/legal_entities/{legal_entity_id}`
  - `app` -> `mdm/v1/legal_entities/{legal_entity_id}`，并额外带上 query `legal_entity_id`
- 常用于合同创建前选择我方法人主体
- `mdm legal get --code` 仅支持 app 身份，走 `GET /open-apis/mdm/v1/legal_entities`，`--code` 映射到 query `legalEntity`
- `create/update` 当前仅支持 app 身份
- `create/update` 必须传 `--user-id`
- `create` 请求体不要带后端生成的 `legalEntity` / `legal_entity` 编码；`update` 请求体必须带后端返回的 `id` 和 camelCase `legalEntity` 编码
- 创建/更新前建议先查 `contract-cli mdm fields list --profile contract --as app --biz-line legal_entity`
