# CLI P2 开发提示词

参照以下信息开发 P2 规划的 CLI 命令。

## 要求

1. 当发现接口和飞书文档描述不匹配，或者有疑问的时候，不要猜测，直接停止并询问。
2. skill 要同步编写，参考原来代码仓的编写方式，格式保持一致。
3. 遵守仓库 AGENTS.md：TDD 先行，先写失败测试，再做最小实现。
4. 每次代码生成或修改后，必须追加记录到 `docs/ai-changes.md`。
5. 请求体统一使用 `--input-file` 或 `--data`，不把复杂 JSON 字段展开成 CLI flags。
6. 新增命令默认按 app 身份开放；如文档明确支持 user 或已有 MCP 路由，再按现有身份分流方式实现。
7. 开放平台通用 query 保持现有规则：默认 `user_id_type=user_id`，支持 `--user-id-type` / `--user-id` 透传。
8. URL 中的 `{contract_id}`、`{payment_id}`、`{payment_record_id}`、`{payment_plan_uuid}`、`{process_instance_id}` 是路径占位符，参考仓库已有实现用位置参数或 flag 获取实际值后拼接，并做 `url.PathEscape`。

## 命令参数约定

- 采用“主操作对象 ID 用位置参数，父资源/上下文 ID 用 flag”的方式。
- `contract_id` 作为父资源时使用 `--contract <contract-id>`。
- `payment_id` 作为父资源时使用 `--payment <payment-id>`。
- `payment_plan_uuid` 使用 `--plan <payment-plan-uuid>`。
- `payment_id` 作为当前操作对象时使用位置参数 `<payment-id>`。
- `payment_record_id` 作为当前操作对象时使用位置参数 `<payment-record-id>`。
- `process_instance_id` 使用位置参数 `<process-instance-id>`。

## 参考已有命令

查询合同分享记录 GET `/open-apis/contract/v1/contracts/{contract_id}/share_records`  
https://ysi13ckdb9.feishu.cn/wiki/LIzmwfxb3iRjdBk6a4HcJwd0nue  
命令： `contract-cli contract share get`

## 需要开发的 P2 接口与命令

创建付款申请 POST `/open-apis/contract/v1/contracts/{contract_id}/payments`  
https://docs.qfei.cn/367580774e0  
命令： `contract-cli payment create`  
参数示例： `contract-cli payment create --contract <contract-id> --input-file payment.json`

更新付款信息 PATCH `/open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}`  
https://docs.qfei.cn/367581314e0  
命令： `contract-cli payment update`  
参数示例： `contract-cli payment update <payment-id> --contract <contract-id> --input-file payment-update.json`

查看付款信息 GET `/open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}`  
https://docs.qfei.cn/367582542e0  
命令： `contract-cli payment get`  
参数示例： `contract-cli payment get <payment-id> --contract <contract-id>`

查询付款申请列表 GET `/open-apis/contract/v1/contracts/{contract_id}/payments`  
https://docs.qfei.cn/367584456e0  
命令： `contract-cli payment list`  
参数示例： `contract-cli payment list --contract <contract-id>`

同步付款记录 POST `/open-apis/contract/v1/payment/notify`  
https://docs.qfei.cn/367585221e0  
命令： `contract-cli payment plan notify`  
参数示例： `contract-cli payment plan notify --input-file notify.json`

搜索付款计划 POST `/open-apis/contract/v1/payments/search`  
https://docs.qfei.cn/367585632e0  
命令： `contract-cli payment plan search`  
参数示例： `contract-cli payment plan search --input-file payment-plan-search.json`

创建付款记录 POST `/open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records`  
https://docs.qfei.cn/367586715e0  
命令： `contract-cli payment record create`  
参数示例： `contract-cli payment record create --contract <contract-id> --payment <payment-id> --input-file payment-record.json`

更新付款记录 PATCH `/open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records/{payment_record_id}`  
https://docs.qfei.cn/367590319e0  
命令： `contract-cli payment record update`  
参数示例： `contract-cli payment record update <payment-record-id> --contract <contract-id> --payment <payment-id> --input-file payment-record-update.json`

查询付款记录详情 GET `/open-apis/contract/v1/contracts/{contract_id}/payments/{payment_id}/payment_records/{payment_record_id}`  
https://docs.qfei.cn/367591242e0  
命令： `contract-cli payment record get`  
参数示例： `contract-cli payment record get <payment-record-id> --contract <contract-id> --payment <payment-id>`

根据付款计划id查询付款记录 GET `/open-apis/contract/v1/contracts/payments/{payment_plan_uuid}/payment_records`  
https://docs.qfei.cn/421418055e0  
命令： `contract-cli payment record list`  
参数示例： `contract-cli payment record list --plan <payment-plan-uuid>`

发起流程审批 POST `/open-apis/contract/v1/process_instances/{process_instance_id}/task_approval`  
https://docs.qfei.cn/367594253e0  
命令： `contract-cli contract approval start`  
参数示例： `contract-cli contract approval start <process-instance-id> --input-file approval.json`

查询审批实例详情 GET `/open-apis/contract/v1/process_instances/{process_instance_id}`  
https://docs.qfei.cn/367595371e0  
命令： `contract-cli contract approval get`  
参数示例： `contract-cli contract approval get <process-instance-id>`

## 实现建议

1. 新增顶层 `payment` 命令分发。
2. 新增 `internal/openplatform/payment` service，负责付款、付款计划、付款记录接口请求。
3. 在现有 `contract` 命令树下新增 `approval` 子命令。
4. POST/PATCH 命令必须要求 `--input-file` 或 `--data`。
5. GET/list 命令不接受请求体。
6. 所有 path 占位符值 trim 后再 `url.PathEscape`。
7. 所有新增命令补 `--help`。
8. 同步更新 `docs/cli-command-reference.md`、`docs/cli-test-plan.md`。
9. 同步更新 `skills/contract-cli-shared`，并新增或更新 payment / contract 对应 skill 文档。
10. 追加 `docs/ai-changes.md`。

## 测试要求

1. 先写失败测试。
2. CLI 测试覆盖每个命令的 method、path、body、query、Authorization。
3. 覆盖缺少 `--contract`、`--payment`、`--plan`、位置 ID、body 时的本地校验。
4. 覆盖 POST/PATCH 缺少 body 时不发 HTTP。
5. 覆盖 `--as user` 如果当前按 app-only 实现，应本地拒绝且不发 HTTP。
6. 覆盖 help topic。
7. 运行 `go test ./...`。
