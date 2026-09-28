# Template Instance Fields

这份文档用于 `contract-cli contract template instantiate`，即根据模板编号和字段值创建模板实例。

官方来源：[创建模版实例](https://docs.qfei.cn/367567259e0)

## 1. CLI 命令面与硬约束

```bash
contract-cli contract template instantiate --profile contract --as app --input-file template-instance.json --user-id-type employee_id
```

硬约束：

- `--as user` 走 `/open-apis/contract/v1/mcp/template_instances`
- `--as app` 走 `POST /open-apis/contract/v1/template_instances`
- `--input-file` 与 `--data` 必须传一个且互斥
- 官方 app 文档的 query 只有 `user_id_type`
- `create_user_id` 与 `create_employee_code` 二选一；官方说明优先取 `create_user_id`
- `template_number` 来自 `contract template list/get`
- 创建补充协议时，`template_field_list` 需要把所需字段全部填写，不会继承原合同模板自定义字段

## 2. 请求体字段

| 字段 | 类型 | 必填性 | 业务含义 |
| --- | --- | --- | --- |
| `source_id` | `string` | 可选 | 外部系统来源单据 id。 |
| `create_user_id` | `string` | 条件必填 | 创建用户 id。与 `create_employee_code` 二选一，优先取该字段。 |
| `create_employee_code` | `string` | 条件必填 | 创建员工工号。与 `create_user_id` 二选一。 |
| `template_number` | `string` | 必填 | 模板编号。 |
| `template_field_list` | `array<object>` | 可选 | 模板字段信息列表。 |

`template_field_list[]`：

| 字段 | 类型 | 必填性 | 业务含义 |
| --- | --- | --- | --- |
| `field_key` | `string` | 必填 | 模板字段编码，对应模板详情里的 `template_fields[].field_code`。 |
| `field_value` | `string` | 必填 | 字段值，必须是 JSON 字符串。 |
| `edit_disabled` | `boolean` | 可选 | 是否禁止用户修改字段值。 |

## 3. `field_value` 格式

`field_value` 是字符串，字符串内容本身是 JSON，统一包在 `content` 里。

文本：

```json
{"content":"名称"}
```

数值：

```json
{"content":10086}
```

日期，格式 `yyyy-MM-dd`：

```json
{"content":"2022-04-21"}
```

日期区间：

```json
{"content":["2022-04-21","2022-04-30"]}
```

单选：

```json
{"content":"option1"}
```

多选：

```json
{"content":["option1","option2"]}
```

附件：

```json
{"content":["file_123"]}
```

金额：

```json
{"content":{"amount":"14.00","currency":"CNY"}}
```

## 4. 最小示例

```json
{
  "source_id": "contract-cli-skill-field-smoke-20260526",
  "create_user_id": "ou_xxx",
  "template_number": "202204160001",
  "template_field_list": [
    {
      "field_key": "name",
      "field_value": "{\"content\":\"名称\"}"
    },
    {
      "field_key": "money",
      "field_value": "{\"content\":{\"amount\":\"14.00\",\"currency\":\"CNY\"}}",
      "edit_disabled": true
    }
  ]
}
```

## 5. 响应读取

读取路径：`data.template_instance`

| 字段 | 类型 | 业务含义 |
| --- | --- | --- |
| `source_id` | `string` | 外部系统单据 id。 |
| `template_number` | `string` | 模板编号。 |
| `template_id` | `string` | 模板 id。 |
| `template_instance_id` | `string` | 模板实例 id，可用于 `contract create` 的 `template_instance_id`。 |
