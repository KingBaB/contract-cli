# mdm file download Parameters

本页专用于 `contract-cli mdm file download`。

- 接口：`GET /open-apis/mdm/v1/file/download/{file_id}`
- 身份：仅 `app`
- 请求体：官方接口不要求 JSON 请求体
- 官方 OpenAPI：[下载主数据附件](https://docs.qfei.cn/373518485e0.md)
- 校验基线：`contract-cli` 当前分支；CLM `master@334c18fc2e` 存在对应实现时，以 Controller、DTO 和业务校验补充官方规格。

## 目录

- [CLI 参数映射](#cli-参数映射)
- [枚举与约束](#枚举与约束)
- [示例](#示例)
- [来源差异说明](#来源差异说明)

## CLI 参数映射

| 参数 | 请求位置 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- | --- |
| <file-id> | $path.file_id | string | 必填（服务端） | 文件id |
| --output-file | 本地文件 | string | 可选 | 保存到指定路径；不传时拉起保存文件弹窗。 |
| --force | 本地写入 | boolean | 可选 | 指定 `--output-file` 时允许覆盖已有文件。 |
| --profile | 本地上下文 | string | 可选 | 不传时使用当前 profile。 |
| --as | 本地上下文 | enum | 可选 | 仅支持 `app`；不传时使用 profile 默认身份。 |
| --user-id-type | $query.user_id_type | string | 可选 | 不传时 CLI 默认发送 `user_id`。 |
| --user-id | $query.user_id | string | 可选 | 传入时透传；MDM create/update 除外。 |
| --output | CLI 输出 | enum | 可选 | `json`、`yaml` 或 `table`；默认 `json`。 |
| --raw | CLI 输出 | boolean | 可选 | 把二进制响应写到 stdout，不打印额外提示。 |

## 枚举与约束

- 使用此方式可以下载小于等于100MB的文件

## 示例

```bash
contract-cli mdm file download <file-id> --profile contract --as app
```

## 来源差异说明

- 官方规格路径：`GET /open-apis/mdm/v1/file/download/{file_id}`
- 本页描述请求参数；响应 envelope 和输出格式遵循共享 Skill。
