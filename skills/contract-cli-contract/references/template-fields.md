# Template Fields

这份文档用于 `contract-cli contract template list` 和 `contract-cli contract template get`，重点帮助你从模板详情里拿到实例化所需字段。

官方来源：

- [查看模板列表](https://docs.qfei.cn/367142097e0)
- [查看模版详情](https://docs.qfei.cn/367571068e0)

## 1. CLI 命令面

```bash
contract-cli contract template list --profile contract --as app --category-number CAT-1 --page-size 20 --user-id ou_xxx --user-id-type employee_id
contract-cli contract template get <template-id> --profile contract --as app --user-id ou_xxx --user-id-type employee_id
```

硬约束：

- `template list --as user` 走 `/open-apis/contract/v1/mcp/templates`
- `template list --as app` 走 `/open-apis/contract/v1/templates`
- `template get --as user` 走 `/open-apis/contract/v1/mcp/templates/{template_id}`
- `template get --as app` 走 `/open-apis/contract/v1/templates/{template_id}`
- 官方 app 文档要求 `user_id`、`user_id_type`，列表还要求 `category_number`；CLI 只透传，不做本地必填校验

## 2. `template list` 查询参数

| CLI 参数 | 请求位置 | 类型 | 业务含义 |
| --- | --- | --- | --- |
| `--category-number` | `$query.category_number` | `string` | 合同模板二级分类编码。可先用 `contract category list` 获取二级分类 `number`。 |
| `--page-size` | `$query.page_size` | `integer` | 分页大小，官方最大值 `100`。 |
| `--page-token` | `$query.page_token` | `string` | 分页令牌。 |
| `--user-id` | `$query.user_id` | `string` | app 场景通常需要传调用人。 |
| `--user-id-type` | `$query.user_id_type` | `string` | 用户 ID 类型。 |

## 3. `template list` 响应字段

读取路径：`data.template_brief_infos[]`

| 字段 | 类型 | 业务含义 |
| --- | --- | --- |
| `template_id` | `string` | 模板 id，用于 `template get <template-id>`。 |
| `template_number` | `string` | 模板编号，用于 `template instantiate`。 |
| `template_name` | `string` | 模板名称。 |
| `description` | `string` | 模板描述。 |
| `edit_enable` | `boolean` | 模板是否可编辑。 |
| `publish_employee_name` | `string` | 发布人姓名。 |
| `publish_time` | `string` | 发布时间。 |

分页字段：

- `data.has_more`
- `data.page_token`

## 4. `template get` 响应字段

读取路径：`data.template`

| 字段 | 类型 | 业务含义 |
| --- | --- | --- |
| `template_id` | `string` | 模板 id。 |
| `template_name` | `string` | 模板名称。 |
| `template_number` | `string` | 模板编号，创建模板实例时传 `template_number`。 |
| `create_employee_name` | `string` | 创建人姓名。 |
| `published` | `boolean` | 是否已发布。 |
| `publish_employee_name` | `string` | 发布人姓名。 |
| `description` | `string` | 模板描述。 |
| `contract_category_name` | `string` | 合同二级分类名称。 |
| `publish_time` | `string` | 发布时间。 |
| `template_fields` | `array<object>` | 模板字段信息列表。 |

`template_fields[]`：

| 字段 | 类型 | 业务含义 |
| --- | --- | --- |
| `field_id` | `string` | 字段 id。 |
| `field_code` | `string` | 字段 code；实例化时对应 `template_field_list[].field_key`。 |
| `field_name` | `string` | 字段名称。 |
| `field_type` | `integer` | 字段类型。 |
| `allow_null` | `boolean` | 是否可为空。 |
| `value_scopes` | `array<object>` | 可选值范围。 |
| `value_scopes[].label` | `string` | 可选值外显名称。 |
| `value_scopes[].value` | `string` | 可选值存储值。 |
| `default_value` | `array<string>` | 默认取值。 |

## 5. 衔接模板实例

流程：

1. `contract category list` 找二级分类 `number`
2. `contract template list --category-number <number>` 找模板
3. `contract template get <template-id>` 读取 `template_number` 和 `template_fields`
4. 按 [template-instance-fields.md](template-instance-fields.md) 组装 `template_field_list`
