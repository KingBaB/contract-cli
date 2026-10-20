# Print File Fields

这份文档用于 `contract-cli contract print-file`，即生成合同打印文件。

官方来源：[生成合同打印文件](https://docs.qfei.cn/367687887e0)

## 1. CLI 命令面与硬约束

```bash
contract-cli contract print-file --profile contract --as app --input-file print-file.json
```

硬约束：

- 当前仅支持 `--as app`
- 走 `POST /open-apis/contract/v1/files`
- `--input-file` 与 `--data` 必须传一个且互斥
- 返回的是平台文件 id，不是文件二进制；下载文件需要再用 `contract download-file <file-id>`

## 2. 请求体字段

| 字段 | 类型 | 必填性 | 业务含义 |
| --- | --- | --- | --- |
| `contract_id` | `string` | 通常必填 | 生成文件的合同 id。 |
| `operate_type` | `integer` | 必填 | 操作类型；当前官方文档只列出 `0`。 |

`operate_type`：

| 值 | 含义 |
| --- | --- |
| `0` | 合同信息文件 |

## 3. 最小示例

```json
{
  "contract_id": "7160880468122992647",
  "operate_type": 0
}
```

命令：

```bash
contract-cli contract print-file --profile contract --as app --data '{"contract_id":"7160880468122992647","operate_type":0}'
```

## 4. 响应读取

读取路径：`data.file_id`

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "file_id": "7107133499601535020"
  }
}
```

生成后下载：

```bash
contract-cli contract download-file 7107133499601535020 --profile contract --as app --output-file ./contract.pdf
```
