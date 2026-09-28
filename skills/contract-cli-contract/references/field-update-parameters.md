# contract field update Parameters

本页专用于 `contract-cli contract field update`。

- 接口：`PUT /open-apis/contract/v1/attribute_definition`
- 身份：仅 `app`
- 请求体：JSON 必填，`--input-file` 与 `--data` 二选一
- 官方 OpenAPI：[更新合同字段信息](https://docs.qfei.cn/367650665e0.md)
- 校验基线：`contract-cli` 当前分支；CLM `master@334c18fc2e` 存在对应实现时，以 Controller、DTO 和业务校验补充官方规格。

## 目录

- [CLI 参数映射](#cli-参数映射)
- [请求体字段](#请求体字段)
- [枚举与约束](#枚举与约束)
- [示例](#示例)
- [来源差异说明](#来源差异说明)

## CLI 参数映射

| 参数 | 请求位置 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- | --- |
| --input-file | $body | JSON file | 二选一必填 | 从文件读取 JSON；与 `--data` 互斥。 |
| --data | $body | JSON string | 二选一必填 | 内联 JSON；与 `--input-file` 互斥。 |
| --profile | 本地上下文 | string | 可选 | 不传时使用当前 profile。 |
| --as | 本地上下文 | enum | 可选 | 仅支持 `app`；不传时使用 profile 默认身份。 |
| --user-id-type | $query.user_id_type | string | 可选 | 不传时 CLI 默认发送 `user_id`。 |
| --user-id | $query.user_id | string | 可选 | 传入时透传；MDM create/update 除外。 |
| --output | CLI 输出 | enum | 可选 | `json`、`yaml` 或 `table`；默认 `json`。 |
| --raw | CLI 输出 | boolean | 可选 | 原样输出服务端响应 body。 |

## 请求体字段

字段名、类型和服务端必填性来自官方 OpenAPI；“CLI 必填/禁止”是结构化命令的额外本地校验。父对象可选时，其内部必填字段标记为“父对象存在时必填”。

| JSON 路径 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- |
| module_name | string | 必填 | 字段组名称 |
| attribute_name | string | 必填 | 字段名称 |
| value_scopes | array<object> | 可选 | 选项列表<br>[{<br>  "label":"深圳",<br>  "value":"sz"<br>},{<br>  "label":"广州",<br>  "value":"gz"<br>}] |
| value_scopes[].label | string | 父对象存在时必填 | label |
| value_scopes[].value | string | 父对象存在时必填 | value |

## 枚举与约束

- 目前仅支持修改下拉列表的选项列表范围

## 示例

```bash
contract-cli contract field update --input-file request.json --profile contract --as app
```

官方请求体示例（动态字段接口仍须以当前租户配置为准）：

```json
{
  "module_name": "甘蒙",
  "attribute_name": "陆一诺",
  "value_scopes": [
    {
      "label": "voluptate fugiat irure",
      "value": "nisi ea reprehenderit"
    }
  ]
}
```

## 来源差异说明

- 官方规格路径：`PUT /open-apis/contract/v1/attribute_definition`
- 本页描述请求参数；响应 envelope 和输出格式遵循共享 Skill。
