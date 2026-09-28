# Category Fields

这份文档用于 `contract-cli contract category list`，即查询合同类型目录。

官方来源：[查询合同类型目录](https://docs.qfei.cn/366393507e0)

## 1. CLI 命令面

```bash
# 后续使用 user 身份创建合同
contract-cli contract category list --profile contract --as user --lang zh-CN

# 后续使用 app 身份创建合同
contract-cli contract category list --profile contract --as app --lang zh-CN
```

硬约束：

- `--as user` 走 `/open-apis/contract/v1/mcp/contract_categorys`
- `--as app` 走 `/open-apis/contract/v1/contract_categorys`
- CLI 把 `--lang` 作为 query 参数透传
- 官方文档的 body 里也出现 `lang`，但当前 CLI 命令面只支持 `--lang`
- 为创建合同选择分类时，查询命令必须使用与后续 `contract create` 相同的 `--profile` 和 `--as`

## 2. 参数

| CLI 参数 | 请求位置 | 类型 | 必填性 | 业务含义 |
| --- | --- | --- | --- | --- |
| `--lang` | `$query.lang` | `string` | 可选 | 语言。常见 `zh-CN`、`en-US`、`ja-JP`。 |

## 3. 响应字段

读取路径：`data.contract_category_resource_vo.category_resources[]`

一级分类字段：

| 字段 | 类型 | 业务含义 |
| --- | --- | --- |
| `name` | `string` | 类型名称。 |
| `number` | `string` | 类型编码，唯一标识。 |
| `abbreviation` | `string` | 类型名称缩写，创建合同时常用。 |
| `id` | `string` | 类型 id。 |
| `description` | `string` | 描述。 |
| `children` | `array<object>` | 二级类型集合。 |

二级分类 `children[]`：

| 字段 | 类型 | 业务含义 |
| --- | --- | --- |
| `name` | `string` | 二级类型名称。 |
| `number` | `string` | 二级合同类型 number；模板列表 `--category-number` 常用这个值。 |
| `abbreviation` | `string` | 二级类型缩写；创建合同 `contract_category_abbreviation` 常用这个值。 |
| `id` | `string` | 二级类型 id。 |
| `description` | `string` | 描述。 |

## 4. 使用建议

- `contract_category_abbreviation` 不是固定枚举；它来自当前租户、当前身份可见的分类树
- 创建合同时选择具体末级分类，并取该节点的 `abbreviation`；常见二级分类路径是 `data.contract_category_resource_vo.category_resources[].children[].abbreviation`
- 查询模板列表需要分类编号时，优先取二级分类的 `number`
- 如果后端返回多级树，先选 `children[]` 中具体业务分类，不要只用一级分类
- `name` 是展示名称、`number` 是分类编号、`id` 是分类 id，三者都不能直接作为创建请求的 `contract_category_abbreviation`
- `contract_category_id` 不能替代创建请求必填的 `contract_category_abbreviation`
- 创建返回 `111717 合同分类简称参数非法` 时，检查是否传了查询结果中的真实末级 `abbreviation`，以及查询和创建是否使用了相同 profile/身份
