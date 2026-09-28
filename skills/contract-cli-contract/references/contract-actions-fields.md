# Contract Actions And File Fields

这份文档轻量覆盖字段较少的动作类和文件类命令：`upload-file`、`download-file`、`submit`、`resubmit`、`delete`。

官方来源：

- [提交合同](https://docs.qfei.cn/367579451e0)
- [重新提交合同](https://docs.qfei.cn/367657850e0)
- [删除草稿合同](https://docs.qfei.cn/367578376e0)

## 1. `contract upload-file`

```bash
contract-cli contract upload-file --profile contract --as app --file ./合同正文.docx --file-type text
```

请求是 `multipart/form-data`：

| CLI 参数 | 表单字段 | 类型 | 必填性 | 业务含义 |
| --- | --- | --- | --- | --- |
| `--file` | `file` | file | 必填 | 本地真实文件。 |
| `--file-type` | `file_type` | `string` | 必填 | 文件类型。 |
| `--file-name` | `file_name` | `string` | 可选 | 上传给后端的文件名；默认本地文件名。 |

常用 `file_type`：

- `text`：合同文本
- `attachment`：其他附件
- `scan`：归档扫描件
- `cause`：合同附件
- `archiveAttachment`：归档附件
- `customPictureAttachment`：图片附件
- `customTableAttachment`：表格附件
- `customFileAttachment`：文件附件

响应重点读取：

- `data.file_id`

注意：

- `--file` 不是 JSON 请求体文件
- 当前 CLI 不支持官方 API 里可能出现的 `need_convert_to_pdf`
- 文件必须是普通文件，大小小于等于 `200MB`

## 2. `contract download-file`

```bash
contract-cli contract download-file <file-id> --profile contract --as app --output-file ./contract.pdf
```

字段：

| CLI 参数 | 请求位置 | 类型 | 必填性 | 业务含义 |
| --- | --- | --- | --- | --- |
| `<file-id>` | `$path.file_id` | `string` | 必填 | 平台文件 id。 |
| `--output-file` | 本地文件 | `string` | 推荐 | 保存路径。 |
| `--force` | 本地行为 | `boolean` | 可选 | 覆盖已存在输出文件。 |
| `--raw` | stdout | `boolean` | 可选 | 把二进制内容写到 stdout。 |

注意：

- 当前仅支持 `--as app`
- 不返回 JSON；成功标准是文件正确写入或 `--raw` 输出二进制
- Agent/CI/远程环境推荐显式传 `--output-file`

## 3. `contract submit`

```bash
contract-cli contract submit <contract-id> --profile contract --as app
```

字段：

| CLI 参数 | 请求位置 | 类型 | 必填性 | 业务含义 |
| --- | --- | --- | --- | --- |
| `<contract-id>` | `$path.contract_id` | `string` | 必填 | 合同 id，目前仅支持通过开发平台创建的草稿合同。 |
| `--input-file` / `--data` | body | JSON | 可选 | CLI 允许透传 body，但官方文档没有定义请求体字段。 |

响应读取：

- `data.contract_id`
- `data.process_instance_id`

## 4. `contract resubmit`

```bash
contract-cli contract resubmit <contract-id> --profile contract --as app
```

字段与响应和 `submit` 一致。

限制：

- 合同被撤回或拒绝后，且当前审批节点为提交节点时可重新提交
- 官方状态说明：`revoked` code `2`，`rejected` code `4`

响应读取：

- `data.contract_id`
- `data.process_instance_id`

## 5. `contract delete`

```bash
contract-cli contract delete <contract-id> --profile contract --as app
```

字段：

| CLI 参数 | 请求位置 | 类型 | 必填性 | 业务含义 |
| --- | --- | --- | --- | --- |
| `<contract-id>` | `$path.contract_id` | `string` | 必填 | 合同 id。 |

限制：

- 当前仅支持 `--as app`
- 仅支持删除同一应用下创建的草稿状态合同
- CLI 不额外要求 `--yes`

响应：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```
