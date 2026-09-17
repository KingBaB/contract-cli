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

## 2026-09-17 协商引用上传最小方案

- 在 `dev-feature-contract-cli-add-review` 新增 `review cooperation get --cooperation-id`，以 user 身份复用个人协商详情接口；历史 app-only 协商命令不变。
- `review file upload` 新增 business-id、business-file-id、business-file-version、source-locator、source-revision-id、source-content-md5、channel-type。无文件时执行引用上传；原文件/stdin模式保留，并允许可选业务ID、业务文件ID、渠道。版本/来源/修订/MD5与文件流互斥。
- 按上传文档 revision 17 校验必填、版本长度、来源 JSON/枚举/字节长度；sourceLocator 中的 ID 只接受字符串，长 ID 和 MD5 原值透传，协商 ID 不复制为顶层表单字段，不引入 appType 输入。
- 更新 review Skill：合同查询 → 从真实 pc_url 取得协商ID → 查询个人协商 → 唯一当前主合同 → cooperation_file_id/version 映射 → 引用上传 → 主体/发起/查询；可选 MD5、修订号缺失不阻断。完成后 pc_url 设置 targetType=ai-review，不因链接缺失重审。
- 本轮只闭合协商取数；任务中心沿用结构化主体和引用协议，但不宣称其取数及页面端到端已验证。特殊禁止场景暂不新增，服务端身份/权限校验保留。上文“上传不接受 businessId”为历史记录，由本节覆盖。
- 先观察新增flag/命令及Skill示例缺失的失败测试，再实现；目标测试、全量 `go test ./...`、`go vet ./...`、`git diff --check` 通过。测试使用假传输和真实 App 入口，不上传真实文件或创建审查任务。未修改后端、环境配置，未全局安装、提交或推送。

## 2026-09-16 审查 Skill 落地（dev 分支）

- 本轮在用户确认的 `dev-feature-contract-cli-add-review` 开发，来源为 CLI 需求 §5.5（revision 235）及四份目标附件。上文 test 分支记录仅代表历史阶段。
- 新增 `skills/contract-cli-review/SKILL.md`：单文件预审、主体/强度/清单选择、一次任务创建与同任务恢复、真实风险展示；区分系统协商/任务中心目标流程和当前未闭合的正文映射、场景及详情导航。
- 按用户确认，将文档原提案名 `contract-cli-review-config` 改为 `contract-cli-review-rule`；职责仍覆盖清单、规则、分组。仅单项写入，保留准确影响确认、复查、预检、回读；引用/原子性不足时不执行级联删除或分组迁移。
- 只增补 shared/contract 的发现描述、路由及既有 Agent 默认提示；保留历史个人审批、文件下载、全部参考资源、授权及审批矩阵身份规则。未整份覆盖旧版附件。
- 主体提取示例适配为真实 flags；task/run JSON 示例通过实际 CLI 预检并验证 int64 文件 ID 不丢精度。新 Skill 不复制安装、授权、环境选择、其他产品归因或固定链接有效期。
- 新增公开入口测试：真实 embed 目录的 list/install、13项分发、包内引用解析、真实 JSON 示例预检、历史源码引用有效性。没有执行真实审查或配置写入，不证明宿主模型运行效果或服务器/网页场景已验收。
- 验证：`go test ./...`、`go vet ./...`、`node --test scripts/install.test.js`、两份新 Skill 与 contract 入口的 quick_validate、`git diff --check`。shared 的既有顶层 version 不在通用校验器白名单内，保留 CLI 历史契约并以真实分发测试及 YAML 解析验证，不改字段规避检查。
- 原有环境修改及其余资源共 99 个基线文件通过哈希保护检查。未改 CLI 业务代码、未全局安装 Skill、未提交推送或发布。

## 2026-09-16 文件审查分类引导补充

- 在 `dev-feature-contract-cli-add-review` 按用户确认修改 `contract-cli-review`：上传成功后展示真实分类及层级路径供用户选择，再收集主体、强度；清单步骤自动使用选定分类和 `--enabled=true` 筛选，不先展示未过滤清单。
- 保留全部类型清单的真实服务语义、完整分页、去重及显式跨类型选择；分类变更重新生成清单候选，不重复上传。分类缺失、映射不明或没有候选时不跳过筛选创建任务。系统合同仍使用已有分类，分类不加入审查任务 JSON。
- 新增测试从内嵌 Skill 读取并执行实际清单示例，验证 user 身份、GET、长整数字符串分类 ID、启用状态及分页参数。旧示例因未传 enabled 正确失败，修改后通过；交互顺序及恢复/空候选边界另作人工指令检查，不宣称真实宿主或服务端端到端验收。
- 验证通过：目标 Skill 测试、`go test ./internal/cli ./skills`、`go vet ./internal/cli`、Skill quick_validate、`git diff --check`。仅改 Skill、对应测试和本摘要；未改 CLI 业务代码、未重新安装或打包，未继续创建当前合同审查任务。

## 2026-09-16 系统内置审查清单置顶修正

- 按用户截图及参数说明，Skill 固定置顶 `0. 通用审查清单（系统内置）`，移除等待远端清单证明内置项存在的条件；自定义清单仍按分类、启用状态筛选并从 1 编号。
- 只选 0 时传 `matchContractTypeRulePackage=true` 并省略自定义 ID；组合选择时同时传两类来源；未选 0 不自动启用。成功查询但无自定义清单时仍可选择 0，查询失败不能冒充空列表。内置项与“全部类型”自定义清单明确区分。
- task start/run 示例覆盖组合方式，新测试执行内嵌真实示例的预检并验证内置单选、组合、自定义单选三种配置。旧示例以缺少 true 正确 RED，修改后 GREEN；固定置顶与空列表分支作人工指令检查，不宣称宿主或服务端端到端验证。
- 目标测试、`go test ./internal/cli ./skills`、`go vet ./internal/cli`、Skill quick_validate、`git diff --check` 通过。未修改 CLI 业务代码、未打包安装；未重试先前返回 503、结果不确定的审查请求，也未将内置项追加到用户此前选择。

## 2026-09-17 任务中心正文引用上传

- review Skill 增加任务中心合同 ID 取数：`contract get` 的同一个 `contract_files.contract_text` 中 file_id/text_id/file_name 对应 businessFileId/businessFileVersion/name，业务 ID 取当前合同，sourceFileType 固定为 2。
- text_id 缺失时报告后端字段缺口，不用 contract_version、file_id、审批任务 ID 或默认值替代，不降级为下载重传。不调用协商详情，不执行主体提取；立场来自合同结构化主体，沿用原任务与 pc_url 导航规则。
- 新增 CLI user 合同详情新字段透传测试，以及执行内嵌 Skill 示例的预检测试，覆盖不同文件/正文长 ID 和缺少版本拒绝。CLI 业务实现、历史命令与协商流程不改。
- 后端必须先提供正文 text_id；未部署前字段缺失仍应停止。本轮未安装 CLI/Skill、未部署或真实发起审查。

## 2026-09-17 系统合同审查链接规则修复

- 修改 review、contract、shared 三个 Skill：系统合同完成审查后从合同 `multi_url.pc_url` 设置 `targetType=ai-review` 并返回可点击链接，不依赖 info 的 url/reviewDetailUrl/documentLocator，也不以网页尚未联调为返回链接的前置条件。
- 恢复任务缺少合同地址时按原合同 ID、profile、user 重新查询；地址缺失或无效才报告链接缺口，不重新上传或发起审查。保留 from、业务定位、空参数和 fragment，不额外添加 AI taskId。外部独立预审仍不返回链接。
- 新增实际执行内嵌 Skill JavaScript 示例的回归测试（测试环境需要 Node.js），覆盖真实协商地址、无 query、重复 targetType/fragment、无效及非 HTTP 地址；先观察缺少示例的 RED，再修改为 GREEN。入口指令和恢复分支另作人工检查，不声称宿主模型或网页端到端验证。
- `go test ./internal/cli ./skills`、`go vet ./internal/cli`、review/contract quick_validate、`git diff --check` 通过。只改 Skill、回归测试及本摘要，未改业务代码、全局安装、提交或推送。

## 2026-09-17 协商审查链接退出来源调整

- 按用户确认，仅 Skill 返回的协商审查链接将 `/cooperation/detail` 的 `from=openPlatform` 改为 `from=allCoopList`；覆盖上节保留该值的旧规则。其他 from 值、缺失 from、任务中心链接及后端原始 pc_url 保持不变。
- 保留 `targetType=ai-review`、业务定位、空参数和 fragment；contract Skill 入口委托 review 的导航规则，避免冲突。未修改 CLI 业务实现、后端和共享接口。
- 真实内嵌 JavaScript 回归覆盖协商转换、已转换/其他来源/缺少来源、任务中心隔离、参数保留、幂等和原有无效 URL 拒绝。先观察两条旧值断言失败，再修改 Skill；不宣称网页端到端验收，不自动安装全局 Skill 或 CLI。
- 验证通过：`go test ./internal/cli ./skills`、`go vet ./internal/cli`、review/contract Skill quick_validate、`git diff --check`；格式校验通过单次 PYTHONPATH 复用已有 PyYAML 缓存，没有安装系统依赖。

## 2026-09-17 所有合同 ID 审查统一文件主体

- 按用户确认，协商及任务中心所有阶段均在引用上传后调用 `review subject extract`，使用同一审查文件的 counterparts 名称/角色供用户选择；覆盖上文任务中心使用合同结构化主体的旧规则。登记交易方和我方/对方不能替代文件主体。
- 更新 review 的任务中心章节、主体规则、阶段表、恢复和 run 边界，以及 contract 审查入口。提取失败、空候选、角色缺失或歧义时停留本步，不猜角色、不回退交易方。尚未发起的旧交易方选择重新提取确认；已有 taskId 只查询，不重建。
- 保留上传业务/版本映射、分类、强度、清单和导航；不改后端、CLI 业务实现与历史命令。合同 ID 流程不以 review run 绕过提取和确认，外部预审能力不变。
- 新增真实内嵌任务中心示例执行测试：先观察缺少提取步骤的 RED，再验证 POST subjects 携带上传后的 businessId/fileId/fileHash，以及返回名称/角色成对保留。交互语义、各阶段和恢复异常路径人工核查，不宣称模型或服务器端到端验收。
- `go test ./internal/cli ./skills`、`go vet ./internal/cli`、review/contract quick_validate、`git diff --check` 通过。未全局安装 CLI/Skill、未发起或重试真实审查、未提交推送。

## 2026-09-17 发起结果不确定时由用户决定重发

- 按用户确认，将未取得 taskId 的创建超时/连接中断从“禁止重发”改为“提示可能已创建及重复任务风险，由用户明确决定是否重发”；不自动重发，用户未确认、拒绝或未回复时保留上下文。
- 用户确认后复用有效上传结果和配置，仅重新调用 task start 一次；再次不确定重新询问，不重跑 run 或重新上传。已有 taskId 继续查询原任务；明确未创建时修正原因后可发起。
- 协商、任务中心及恢复段统一引用一个规则正文，保留首次发起不额外确认及原有后端鉴权。未修改 CLI 网络重试、后端幂等或任务找回能力。
- 新增嵌入 Skill 指令契约回归，先观察旧禁令导致的 RED，再验证 GREEN；人工核查各分支一致性。go test ./internal/cli ./skills、go vet ./internal/cli、git diff --check 通过。不宣称模型执行或线上端到端验证；未全局安装、未真实发起审查、未提交推送。
