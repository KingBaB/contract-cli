# mdm legal get --code Parameters

本页专用于 `contract-cli mdm legal get --code`。

- 接口：`GET /open-apis/mdm/v1/legal_entities`
- 身份：仅 `app`
- 请求体：官方接口不要求 JSON 请求体
- 官方 OpenAPI：[获取法人实体](https://docs.qfei.cn/373518112e0.md)
- 校验基线：`contract-cli` 当前分支；CLM `master@334c18fc2e` 存在对应实现时，以 Controller、DTO 和业务校验补充官方规格。

## 目录

- [CLI 参数映射](#cli-参数映射)
- [枚举与约束](#枚举与约束)
- [示例](#示例)
- [来源差异说明](#来源差异说明)

## CLI 参数映射

| 参数 | 请求位置 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- | --- |
| --user-id-type | $query.user_id_type | string | 可选，默认 `user_id` | 用户类型 |
| --code | $query.legalEntity | string | 必填（服务端） | 法人实体编码 |
| --page-size | $query.page_size | integer | 可选 | 分页大小 |
| --page-token | $query.page_token | string | 可选 | 分页标记，第一次请求不填，表示从头开始遍历；分页查询结果还有更多项时会同时返回新的 page_token，下次遍历可采用该 page_token 获取查询结果 |
| --profile | 本地上下文 | string | 可选 | 不传时使用当前 profile。 |
| --as | 本地上下文 | enum | 可选 | 仅支持 `app`；不传时使用 profile 默认身份。 |
| --user-id | $query.user_id | string | 可选 | 传入时透传；MDM create/update 除外。 |
| --output | CLI 输出 | enum | 可选 | `json`、`yaml` 或 `table`；默认 `json`。 |
| --raw | CLI 输出 | boolean | 可选 | 原样输出服务端响应 body。 |

## 枚举与约束

- 根据法人实体编码，来获取对应的法人实体信息。参数均采用驼峰式

## 示例

```bash
contract-cli mdm legal get --code <code> --profile contract --as app
```

## 来源差异说明

- 官方规格路径：`GET /open-apis/mdm/v1/legal_entities`
- 本页描述请求参数；响应 envelope 和输出格式遵循共享 Skill。
