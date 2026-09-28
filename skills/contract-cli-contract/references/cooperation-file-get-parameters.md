# contract cooperation file get Parameters

本页专用于 `contract-cli contract cooperation file get`。

- 接口：`GET /open-apis/contract/v1/contracts/{contract_id}/cooperation/file_info`
- 身份：仅 `app`
- 请求体：官方接口不要求 JSON 请求体
- 官方 OpenAPI：[获取最新协商文件信息](https://docs.qfei.cn/465721208e0.md)
- 校验基线：`contract-cli` 当前分支；CLM `master@334c18fc2e` 存在对应实现时，以 Controller、DTO 和业务校验补充官方规格。

## 目录

- [CLI 参数映射](#cli-参数映射)
- [枚举与约束](#枚举与约束)
- [示例](#示例)
- [来源差异说明](#来源差异说明)

## CLI 参数映射

| 参数 | 请求位置 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- | --- |
| <contract-id> | $path.contract_id | string | 必填（服务端） | 合同id |
| --profile | 本地上下文 | string | 可选 | 不传时使用当前 profile。 |
| --as | 本地上下文 | enum | 可选 | 仅支持 `app`；不传时使用 profile 默认身份。 |
| --user-id-type | $query.user_id_type | string | 可选 | 不传时 CLI 默认发送 `user_id`。 |
| --user-id | $query.user_id | string | 可选 | 传入时透传；MDM create/update 除外。 |
| --output | CLI 输出 | enum | 可选 | `json`、`yaml` 或 `table`；默认 `json`。 |
| --raw | CLI 输出 | boolean | 可选 | 原样输出服务端响应 body。 |

## 枚举与约束

- CLM `master` 按数字 ID 解析路径参数，必须传数字字符串。

## 示例

```bash
contract-cli contract cooperation file get <contract-id> --profile contract --as app
```

## 来源差异说明

- 官方规格路径：`GET /open-apis/contract/v1/contracts/{contract_id}/cooperation/file_info`
- 本页描述请求参数；响应 envelope 和输出格式遵循共享 Skill。
