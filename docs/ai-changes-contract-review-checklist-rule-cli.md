# AI 变更摘要：合同审查、清单与清单规则 CLI 能力

## 范围

- 在 `test-feature-contract-cli-add-review` 分支增加与 Everyline CLI 对齐的合同审查命令入口：文件上传、URL 上传、主体提取、任务发起、状态查询和详情查询。
- 增加审查清单的查询、创建、批量创建、更新、批量更新、删除和批量删除。
- 在 `review` 下新增 `review checklist` 和 `review rule` 命令组，提供清单、审查规则分组及分组内规则的 CRUD/批量操作；历史 `rule table` 审批矩阵规则命令保持原形态不变。

## API 路径

- 合同审查：`/open-apis/contract/v3/review/**`。
- 审查清单与审查规则：`/open-apis/contract/v1/review-rules/**`。
- `review` 下的审查、清单和规则命令仅允许 user 身份；未传 `--as` 时自动使用 user。
- 审查发起和一键流程通过 `--input` 或 `--data` 传入 JSON 对象，兼容 `--input-file`；这些输入别名仅作用于 review 命令树。
- 批量删除支持重复 `--id` 或 JSON ID 数组，并统一发送为 `{"ids":[...]}`。
- 删除操作要求显式提供 `--yes`。

## 验证

- `go test ./...`
- `go vet ./...`
- `git diff --check`
- CLI help smoke checks for `review`, `checklist`, and `rule`.

## 2026-09-16 八类命令参数对齐

- 对照 EveryLine `test` 分支 `a76df4d6cf289cab5b6db4e9536eb8dcaf688e06`。
- 文件上传增加 `--stdin`，与 `--file` 互斥；限制 2 MiB、doc/docx/pdf。URL 上传校验 HTTP/HTTPS 和业务文件名。
- 主体提取使用 `--business-id`、`--file-id`、可选 `--file-hash`，发送 JSON 对象及字符串 fileId。
- 任务发起改为严格 JSON 对象输入，校验必填项、未知字段、指纹和规则来源；中文审查强度转换为 0/1/2。
- `review run` 执行上传、可选主体提取、发起任务和可选等待；支持 source/config/extractSubjects/wait 和 `--interval`、`--deadline`。
- `task status/info/result` 使用 task-id/business-id/visibility-scope；result 轮询到终态后取详情，失败/取消/超时输出已有快照。review 新增的 `--app-type` 不再暴露。
- 清单、规则和规则分组查询对齐 CSV/重复多值、分页和时间范围校验；零分页/时间参数省略。保留包含逗号的排序项示例：`--sort '"name,desc"'`。
- 帮助同步业务参数、示例、user-only 限制及真实工作流说明。全局解析器恢复历史行为，review 专用 `--input` 不影响历史命令。
- review 及所有子命令的 Notes 不展示接口路径和 HTTP 方法；保留参数约束、身份限制及使用说明。
- 按 cli-command-reference 的命令组/叶子规范补全 review 帮助：叶子统一 Name、中文说明、Usage、Flags、Examples、Notes；采用 input-file 主展示和 input 兼容别名，补齐清单/规则写操作输入及示例。帮助按实际解析能力标出规则写操作限制，未修改业务解析或请求行为。
- 本轮不扩展清单/规则写操作；不复制 EveryLine 专属用量归因编码，不修改认证、环境或网关路由。
- 按用户确认，task start 和 review run 的 startReview 请求固定注入 `usageReportContext.reportBusinessCode=clm_100_openApi_cli`；输入中的归因字段被拒绝，上传、主体提取和查询不注入该字段。
- 公开命令测试通过受控 HTTP 校验路径、请求体、ID 精度、调用顺序、失败和超时边界；未调用真实审查接口。

## 2026-09-16 第一阶段写命令修复

- 规则单项更新和删除统一使用 `--group-id` 与 `--rule-id`，修复两层参数解析冲突。
- 清单、规则及分组写命令支持 `--dry-run` / `--print-input`；预览不访问远端，只有实际删除需要 `--yes`。
- 批量删除支持重复及 CSV ID，去除首尾空白、稳定去重，拒绝空 ID、空列表和 ID/JSON 混用。
- 单项写入必须是 JSON 对象，批量写入必须是非空对象数组；更新时检查可选 body.id 与路径 ID 一致。大整数、未知业务字段及同值 body.id 保留原样。
- 所有 review 远端请求使用审查专属业务码判断及默认 30 秒请求超时，保留更短的上层 deadline、共享鉴权及写操作不确定语义，不改变成功输出结构。
- 同步帮助参数、示例和删除确认说明，移除旧的不可用提示；Notes 不展示接口路径。
- 公开 App.Run 回归测试与本机 HTTP 服务验证请求契约；全仓 Go 测试、目标 race 测试及 go vet 通过。未调用真实审查或清单写接口。
- 第一阶段只补输入形态与命令基础能力；业务字段必填/类型模型、上传 businessId 透传、文件限制等仍属于第二阶段，不宣称完整对齐。

## 边界（保持不变）

- 2026-09-16 第二阶段文件限制：用户明确将审查文件上限由 2 MiB 改为 50 MiB（52,428,800 字节），支持 doc/docx/pdf；本地文件、stdin、review run 本地 source 共用新上限，JSON 输入仍保留 2 MiB，URL 文件由服务端下载校验。
- 新增公开入口边界测试：2 MiB+1、50 MiB、50 MiB+1，本地与 run 预检、实际 stdin multipart、超限无 HTTP、格式及 JSON 限制隔离。未在真实环境上传大文件。
- 第二阶段只读核对：当前上传入口不接受 businessId，不复制 EveryLine 的该上传字段。清单写入模型不接收 reviewStage/enabled，分组/规则写入不接收 sourceType；规则支持 riskTips 和 riskLevel=0..3。是否严格拒绝未知写入字段待用户确认，本次未修改这些行为。
- 当前后端 config 还支持 auditerName/reviewRules 等超出 EveryLine CLI 白名单的字段；后续需独立确定开放范围，不把当前 CLI 声称为服务端全量能力。

- 本次只修改 CLI，不修改 open-platform、网关、部署配置、数据库或生产环境。
- 当前分支用于 test 环境开发，后续上线迁移到 `feature-contract-cli-add-review` 时需重新执行对应构建与环境验证。
