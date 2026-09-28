# mdm vendor list-all Parameters

本页专用于 `contract-cli mdm vendor list-all`。

- 接口：`GET /open-apis/mdm/v1/vendors/list_all`
- 身份：仅 `app`
- 请求体：官方接口不要求 JSON 请求体
- 官方 OpenAPI：[获取交易方全量数据](https://docs.qfei.cn/373500124e0.md)
- 校验基线：`contract-cli` 当前分支；CLM `master@334c18fc2e` 存在对应实现时，以 Controller、DTO 和业务校验补充官方规格。

## 目录

- [CLI 参数映射](#cli-参数映射)
- [枚举与约束](#枚举与约束)
- [示例](#示例)
- [来源差异说明](#来源差异说明)

## CLI 参数映射

| 参数 | 请求位置 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- | --- |
| --page-size | $query.page_size | integer | 可选 | 分页大小<br>示例值：10<br>数据校验规则：<br>最大值：20 |
| --page-token | $query.page_token | string | 可选 | 分页标记，第一次请求不填，表示从头开始遍历；分页查询结果还有更多项时会同时返回新的 page_token，下次遍历可采用该 page_token 获取查询结果<br>示例值："此处可不填" |
| --user-id-type | $query.user_id_type | string | 可选，默认 `user_id` | 用户 ID 类型，参考 用户身份体系 |
| --profile | 本地上下文 | string | 可选 | 不传时使用当前 profile。 |
| --as | 本地上下文 | enum | 可选 | 仅支持 `app`；不传时使用 profile 默认身份。 |
| --user-id | $query.user_id | string | 可选 | 传入时透传；MDM create/update 除外。 |
| --output | CLI 输出 | enum | 可选 | `json`、`yaml` 或 `table`；默认 `json`。 |
| --raw | CLI 输出 | boolean | 可选 | 原样输出服务端响应 body。 |

## 枚举与约束

- 交易方全量数据分页查询张。参数均采用驼峰式
- 官方规格使用 `pageSize` / `pageToken`；当前 CLI 实际发送 `page_size` / `page_token`。

## 示例

```bash
contract-cli mdm vendor list-all --profile contract --as app
```

## 来源差异说明

- 官方规格路径：`GET /open-apis/mdm/v1/vendors/list_all`
- 本页描述请求参数；响应 envelope 和输出格式遵循共享 Skill。
