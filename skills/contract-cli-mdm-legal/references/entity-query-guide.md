# Entity Query Guide

这份文档是 `contract-cli mdm legal ...` 的主入口，优先解决两个问题：

- 我现在是要找候选法人实体，还是已经拿到了法人实体 id
- 这个命令面的 `--name` 到底如何映射到底层查询参数

推荐阅读顺序：

1. 先看本文件，选查询场景
2. 再看 [entity-query-parameters.md](entity-query-parameters.md) 查精确参数映射
3. 最后看 [commands.md](commands.md) 直接抄命令示例

## 1. 命令面与硬约束

当前结构化命令分为查询和 app-only 写入：

```bash
contract-cli mdm legal list --profile contract --name "上海主体"
contract-cli mdm legal get 7023646046559404327 --profile contract
contract-cli mdm legal get --profile contract --as app --code L0001
contract-cli mdm legal create --profile contract --as app --user-id <operator-user-id> --input-file legal-create.json
contract-cli mdm legal update 7003410079584092448 --profile contract --as app --user-id <operator-user-id> --input-file legal-update.json
```

硬约束：

- 不暴露 `--operator`
- 默认输出就是开放平台原始 envelope；如果要脚本消费，建议加 `--output json`
- 如果要排障或确认原始响应，建议加 `--raw`
- `mdm legal list` 同时支持 `user` 和 `app`
- `mdm legal get` 也同时支持 `user` 和 `app`
- `mdm legal get --code` 当前仅支持 `app`
- `mdm legal create/update` 当前仅支持 `app`
- `mdm legal create/update` 必须传 `--user-id`，用于提供当前操作人上下文
- `create` 请求体不要包含后端生成的 `legalEntity` / `legal_entity` 编码；`update` 请求体必须包含后端返回的 `id` 和 `legalEntity` 编码
- `update` 编码字段必须使用 camelCase `legalEntity`，不要写成 `legal_entity`
- 创建/更新请求体除上述规则外仍直接透传，其他字段是否必填以 `mdm fields list --biz-line legal_entity` 和后端配置为准

## 2. 场景配方

### 2.1 按名称找我方法人主体

适用场景：

- 创建合同前只知道主体名称
- 需要拿候选列表，再从结果里挑法人实体 id

最小命令：

```bash
contract-cli mdm legal list --profile contract --name "上海主体"
```

常见追加参数：

- `--page-size 20`
- `--page-token <next-token>`
- `--as app --user-id-type employee_id`

补充说明：

- user 路由走 `/open-apis/contract/v1/mcp/legal_entities`
- app 路由走 `/open-apis/mdm/v1/legal_entities/list_all`
- 文档显示文本使用 `legal_entities/list_all`，但超链接目标误指到了 `vendors`
- 文档还写了“查询参数采用驼峰式”，CLI 当前仍保持 `legalEntity/page_size/page_token` 的既有透传映射

### 2.2 分页扫法人实体列表

适用场景：

- 不按名字过滤，直接分页遍历
- 或者已经有上一页返回的 `page_token`

最小命令：

```bash
contract-cli mdm legal list --profile contract --page-size 20
```

app 示例：

```bash
contract-cli mdm legal list --profile contract --as app --name "主体A" --page-size 20 --user-id-type employee_id
```

### 2.3 已知 id 直接查详情

适用场景：

- 已经从搜索结果、外部系统或历史合同里拿到了法人实体 id

最小命令：

```bash
contract-cli mdm legal get 7023646046559404327 --profile contract
```

app 示例：

```bash
contract-cli mdm legal get 7003410079584092448 --profile contract --as app --user-id-type employee_id
```

补充说明：

- user 路由走 `/open-apis/contract/v1/mcp/legal_entities/{legal_entity_id}`
- app 路由走 `/open-apis/mdm/v1/legal_entities/{legal_entity_id}`
- 按这次确认方案，app 除了 path 参数外，还会额外拼接同名 query `legal_entity_id`
- 文档把 `legal_entity_id` 写在查询参数表里，所以 CLI 采用“path + query 双带”的保守实现

### 2.4 按法人实体编码查询

适用场景：

- 已知法人实体编码，但还没有法人实体 id
- 需要调用生产文档里的 `GET /open-apis/mdm/v1/legal_entities`，而不是全量列表 `list_all`

最小命令：

```bash
contract-cli mdm legal get --profile contract --as app --code L0001
```

常见追加参数：

- `--page-size 20`
- `--page-token <next-token>`
- `--user-id-type employee_id`

补充说明：

- `--code` 模式仅支持 app 身份
- `--code` 映射到底层 query `legalEntity`
- 路由走 `/open-apis/mdm/v1/legal_entities`
- 不接受 `--input-file` / `--data`

### 2.5 创建或更新法人实体

适用场景：

- 需要把外部法人主体同步到合同主数据
- 已知法人实体 id，需要修改主体字段

最小命令：

```bash
contract-cli mdm legal create --profile contract --as app --user-id <operator-user-id> --input-file legal-create.json
contract-cli mdm legal update 7003410079584092448 --profile contract --as app --user-id <operator-user-id> --input-file legal-update.json
```

`legal-create.json` 示例不要带后端生成的 `legalEntity` / `legal_entity` 编码：

```json
{
  "legal_entity_text": "法大大我方11",
  "status": 1
}
```

`legal-update.json` 示例必须带后端返回的 `id` 和 `legalEntity` 编码：

```json
{
  "id": "7003410079584092448",
  "legalEntity": "L0001",
  "legal_entity_text": "法大大我方11",
  "status": 1
}
```

补充说明：

- `create` 走 `POST /open-apis/mdm/v1/legal_entities`
- `update` 走 `PUT /open-apis/mdm/v1/legal_entities/{legal_entity_id}`
- 写接口必须传 `--user-id`
- 字段配置是动态的，不要只凭负向样例判断必填；先查 `mdm fields list --biz-line legal_entity`

## 3. 什么时候不要走这里

- 想先确认法人实体字段定义：改看 [../../contract-cli-mdm-fields/SKILL.md](../../contract-cli-mdm-fields/SKILL.md)
- 想查合同我方主体选择逻辑：回到 [../../contract-cli-contract/SKILL.md](../../contract-cli-contract/SKILL.md)
