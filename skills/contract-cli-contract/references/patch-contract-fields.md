# Patch Contract Fields

这份文档用于 `contract-cli contract patch`。当前官方“更新合同”文档只覆盖文件、扫描件和归档附件相关更新，不要把它理解成任意更新合同基础字段。

官方来源：[更新合同](https://docs.qfei.cn/367579742e0)

## 1. CLI 命令面与硬约束

```bash
contract-cli contract patch <contract-id> --profile contract --as app --input-file contract-patch.json
```

硬约束：

- 当前仅支持 `--as app`
- 走 `PATCH /open-apis/contract/v1/contracts/{contract_id}`
- `--input-file` 与 `--data` 必须传一个且互斥
- 官方字段均是文件 id 或文件 id 映射，文件 id 需要先通过 `contract upload-file` 获取
- 官方 schema 把多个字段标为 required，但字段描述又按合同状态限定；实际调用时以目标合同状态和后端校验为准

## 2. 字段表

| 字段 | 类型 | 业务含义 | 状态限制 |
| --- | --- | --- | --- |
| `ocr_file_id` | `string` | OCR 扫描件文件 id。 | 仅支持合同状态为盖章中，`contract_status_code=6`。 |
| `archive_attachment_map` | `object` | 合同附件扫描件文件 id 映射。key 是合同附件文件 id，value 是扫描件文件 id。 | 扫描件类型为 `archiveAttachment`，仅支持归档中，`contract_status_code=8`。 |
| `scan_file_id` | `string` | 归档文件 id。 | 仅支持归档中，`contract_status_code=8`。 |
| `archive_attachment_file_ids` | `array<string>` | 合同附件扫描件文件 id 列表。 | 仅支持已归档，`contract_status_code=9`。 |

## 3. 示例

盖章中合同补 OCR 扫描件：

```json
{
  "ocr_file_id": "file_ocr_xxx"
}
```

归档中合同补归档主文件：

```json
{
  "scan_file_id": "file_scan_xxx"
}
```

归档中合同补附件扫描件映射：

```json
{
  "archive_attachment_map": {
    "attachment_file_id": "archive_attachment_file_id"
  }
}
```

已归档合同补附件扫描件列表：

```json
{
  "archive_attachment_file_ids": [
    "archive_attachment_file_id"
  ]
}
```

## 4. 不要这样做

- 不要用 `{"title":"demo"}`、`{"contract_name":"demo"}` 这类基础字段 payload 当作已确认能力
- 不要传本地文件路径；这里的值都应该是平台 `file_id`
- 不要对非专用测试数据执行 patch 验收
