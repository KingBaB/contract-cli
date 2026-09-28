---
name: contract-cli-mdm-file
version: 1.0.1
description: "contract-cli 主数据文件下载技能：用 app 身份下载 `/open-apis/mdm/v1/file/download/{file_id}` 主数据附件。当用户要使用 `contract-cli mdm file download` 下载主数据文件时触发。"
---

# contract-cli MDM File

CRITICAL — 开始前 MUST 先读取 [../contract-cli-shared/SKILL.md](../contract-cli-shared/SKILL.md)。

## 适用命令

- `contract-cli mdm file download <file-id>`

## 关键规则

- 当前仅支持 `--as app`
- 走 `GET /open-apis/mdm/v1/file/download/{file_id}`
- 默认拉起保存文件弹窗；Agent/CI/远程环境优先传 `--output-file`
- 已有目标文件时，必须显式传 `--force` 才允许覆盖
- 管道或脚本消费二进制时用 `--raw`
- 不接受 `--input-file` / `--data`
- 完整参数、类型和约束：读 [references/file-download-parameters.md](references/file-download-parameters.md)

## 示例

```bash
contract-cli mdm file download <file-id> --profile contract --as app --output-file ./attachment.bin
contract-cli mdm file download <file-id> --profile contract --as app --output-file ./attachment.bin --force
contract-cli mdm file download <file-id> --profile contract --as app --raw > attachment.bin
```
