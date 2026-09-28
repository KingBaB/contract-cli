# contract authorization grant Parameters

本页专用于 `contract-cli contract authorization grant`。

- 接口：`POST /open-apis/contract/v1/authorizations`
- 身份：仅 `app`
- 请求体：JSON 必填，`--input-file` 与 `--data` 二选一
- 官方 OpenAPI：[授予合同权限](https://docs.qfei.cn/367596624e0.md)
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
| --user-id-type | $query.user_id_type | string | 可选，默认 `user_id` | 用户 ID 类型 详情参考：用户身份体系<br>示例值："user_id"<br>可选值有：<br>- user_id：标识一个用户在某个租户内的身份。同一个用户在租户 A 和租户 B 内的 User ID 是不同的。在同一个租户内，一个用户的 User ID 在所有应用（包括商店应用）中都保持一致。User ID 主要用于在不同的应用间打通用户数据。<br>- union_id：标识一个用户在某个应用开发商下的身份。同一用户在同一开发商下的应用中的 Union ID 是相同的，在不同开发商下的应用中的 Union ID 是不同的。通过 Union ID，应用开发商可以把同个用户在多个应用中的身份关联起来。 |
| --input-file | $body | JSON file | 二选一必填 | 从文件读取 JSON；与 `--data` 互斥。 |
| --data | $body | JSON string | 二选一必填 | 内联 JSON；与 `--input-file` 互斥。 |
| --profile | 本地上下文 | string | 可选 | 不传时使用当前 profile。 |
| --as | 本地上下文 | enum | 可选 | 仅支持 `app`；不传时使用 profile 默认身份。 |
| --user-id | $query.user_id | string | 可选 | 传入时透传；MDM create/update 除外。 |
| --output | CLI 输出 | enum | 可选 | `json`、`yaml` 或 `table`；默认 `json`。 |
| --raw | CLI 输出 | boolean | 可选 | 原样输出服务端响应 body。 |

## 请求体字段

字段名、类型和服务端必填性来自官方 OpenAPI；“CLI 必填/禁止”是结构化命令的额外本地校验。父对象可选时，其内部必填字段标记为“父对象存在时必填”。

| JSON 路径 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- |
| business_id | string | 条件必填 | 合同id |
| business_type_code | integer | 可选 | 默认为0，0代表合同申请场景 |
| authorized_user_id | string | 必填 | 授权者id，参考 用户身份体系 |
| start_time | string | 必填 | 授权开始时间，毫秒级时间戳 |
| end_time | string | 必填 | 授权结束时间，毫秒级时间戳 |
| permanent | boolean | 可选 | 是否永久授权，默认为false |
| remark | string | 可选 | 备注 |
| source_system | string | 可选 | 来源系统<br>示例值："ZhishuOpenPlatform" |

## 枚举与约束

- CLM `master` 允许 `business_id` 与 `business_list` 二选一；授权用户和起止时间必须提供，业务 ID 必须是数字字符串。
- `business_id`（string，条件必填）：合同id
- `business_type_code`（integer，可选）：默认为0，0代表合同申请场景
- `start_time`（string，必填）：授权开始时间，毫秒级时间戳
- `end_time`（string，必填）：授权结束时间，毫秒级时间戳
- `permanent`（boolean，可选）：是否永久授权，默认为false

## 示例

```bash
contract-cli contract authorization grant --input-file request.json --profile contract --as app
```

官方请求体示例（动态字段接口仍须以当前租户配置为准）：

```json
{
  "business_id": "6965467645105668385",
  "business_type_code": 0,
  "authorized_user_id": "ae721f86",
  "start_time": "1626850544000",
  "end_time": "1626850999000",
  "permanent": true,
  "remark": "fjdkfj",
  "source_system": "ZhishuOpenPlatform"
}
```

## 来源差异说明

- 官方规格路径：`POST /open-apis/contract/v1/authorizations`
- 本页描述请求参数；响应 envelope 和输出格式遵循共享 Skill。
