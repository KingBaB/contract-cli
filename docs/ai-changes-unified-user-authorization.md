# 豆包与 WorkBuddy 统一用户授权变更记录

## 变更目标

保留旧 `auth login --as user` Authorization Code 模式，新增适配豆包与 WorkBuddy 的一次发起、一次查询 Device 授权模式，并将凭证和重试边界固定在 CLI 内。

## 主要改动

- 新增 `auth init` 和 `auth complete`；两者都单次请求后退出，不持续轮询。
- `auth init` 在发起远程请求前先验证 CredentialStore 可用；二维码和 pending transaction 保存成功后才切换 profile，前置失败不会破坏旧授权模式和旧 Token。
- Device client、scope、business 和 resource 来自明确 profile 配置；缺失直接失败，不使用默认 scope。
- `auth init` 仅输出完整 HTTPS 授权链接、二维码路径和过期时间；`auth complete` 返回 pending/succeeded/denied/expired/uncertain/restart_required/busy。
- 服务端返回 `slow_down` 时 `auth complete` 仍只返回一次 pending；返回 `invalid_grant` 时保留终态，只有用户明确同意后才能通过 `auth init --restart` 创建新会话。
- Device Code、Access Token、Refresh Token 和密钥不输出到 stdout、stderr 或日志。
- 豆包使用 `$SKILL_SESSION_WORKSPACE` 下的 AES-256-GCM 加密凭证文件，密钥必须由 `CONTRACT_CLI_CREDENTIAL_KEY_V1` 明确注入。
- WorkBuddy 必须提供 `CODEBUDDY_SESSION_ID`，凭证分别写入 macOS Keychain、Windows Credential Manager 或 Linux Secret Service；无安全存储时直接失败。
- Access Token 距过期不足 5 分钟时，使用 profile + task 粒度的跨进程非阻塞文件锁刷新，锁内二次检查。
- 仅 HTTP 401 同时包含可信 `X-Qfei-Open-Platform-Auth-Error: token_expired` 响应头和 `data.error_type=token_expired` 响应体时触发一次刷新；body-only、header-only 和非 401 均不会刷新。可重放请求才会再发一次，流式请求体不会被错误重放。
- Refresh Token 返回 `invalid_grant` 时，在凭证操作锁内清理失效 Token、保留 pending 授权并明确要求先获得用户同意，再查询真实授权状态并决定复用会话、普通 init 或 restart；临时错误不删除凭证，服务端轮换成功但安全存储失败时不自动重试。
- 读请求仅对明确临时网络错误重试一次；写请求的网络错误、超时或 5xx 均不重试，返回结果不确定错误。
- Skill 授权指引与业务恢复规则已更新；正式发布时要求豆包只允许固定结构化工具，不暴露 Shell、Python、raw CLI 或任意 API path。
- 安装器强制校验 SHA-256 checksum，不回退本地源码编译。

## 验证

- `go test ./...`：通过。
- `node --test scripts/install.test.js`：通过。
- `make release-check`：通过，覆盖 Go 单测、CLI smoke、npm package dry-run、本地安装、发布脚本和构建参数检查。
- macOS Intel、macOS Apple Silicon、Windows amd64/arm64、Linux amd64/arm64 六目标 `CGO_ENABLED=0` 构建通过。
- npm package dry-run 和本地安装校验通过。

## 发布阻塞项

- 仓库当前没有豆包平台可直接导入的官方 manifest 模板或导出物，因此尚不能用自造格式宣称“固定二进制、固定子命令、JSON Schema 参数数组执行”已经由平台配置强制生效。
- 正式发布前必须从豆包 Skill 后台导出一份官方 manifest/工具定义模板，纳入本仓库并增加校验测试；该项不阻塞 CLI 与服务端功能开发，但阻塞豆包正式发布。

## 2026-08-06 Review 修复：豆包 profile 安全恢复

- 在现有加密 `DeviceCredential` 中新增可选 `device_profile` 快照，只保存 Device 运行必需的环境、资源、client、scope 和 endpoint 元数据。
- `auth init` 成功后，将清理过敏感字段的快照与 pending transaction 一起写入 CredentialStore；不会把 Access Token、Refresh Token、device code、App 身份、手机号、企业 ID 或业务参数写入快照。
- 豆包沙箱重建且临时 HOME 中 profile 缺失时，仅在存在 `SKILL_SESSION_WORKSPACE`、命令显式传入 profile 名称且加密快照完整匹配时恢复。
- 恢复结果只写入当前临时 HOME 的 `config.json`；工作区仍只有 AES-256-GCM 密文，不写入明文 profile 或 `secrets.json`。
- 快照缺失、名称不匹配或关键字段缺失时明确失败，提示重新执行 `config add` 和 `auth init`，不猜测环境、scope、client 或 endpoint。
- WorkBuddy、旧 Authorization Code profile、App 身份和已有本地 profile 的读取优先级保持不变。
- `skills/auth/SKILL.md` 已补充恢复规则，并要求合同业务调用始终显式携带 `--profile contract`。

### 本轮验证

- `go test ./internal/config ./internal/credential ./internal/cli`：通过。
- 新增豆包相同工作区、全新 HOME 的授权完成与合同查询回归用例。
- 新增快照缺失、名称不匹配和关键字段缺失的失败关闭用例，均确认不会发起 HTTP 请求。

## 2026-08-06 Review 修复：授权完成后的安全存储失败

- Device Token 兑换成功但本地 CredentialStore 保存失败时，不再次请求 Token Endpoint。
- 错误会明确说明当前授权状态可能已经失效，并要求用户处理安全存储问题后重新执行 `auth init`。
- 错误和日志只包含 profile 与底层存储错误，不输出 device code、Access Token 或 Refresh Token。
- 新增回归测试，验证 Token Endpoint 仅调用一次、错误包含重新授权指引且所有凭证均不会进入输出。

## 2026-08-06 WorkBuddy 重复授权修复

- Pending transaction 新增显式状态和完整授权链接；旧凭证缺少状态时按 `pending` 读取，缺少链接时失败关闭。
- `auth init` 遇到有效 pending 时只返回原链接并重新生成二维码，不再覆盖或创建新的服务端会话。
- 服务端创建 Device 会话后先安全持久化 pending，再生成二维码；二维码文件写入失败也不会丢失已创建的授权会话并产生新 code。
- 新增 `auth init --restart`；只有用户明确同意重新授权后，Skill 才能用它替换异常或终态会话。
- `auth complete` 请求 Token Endpoint 前先持久化 `checking`；网络异常、响应解析失败或未知服务端错误统一转为 `uncertain`，后续调用不再访问 Token Endpoint。
- denied、expired 和 invalid_grant 保留为终态，不自动清理并重新发起；成功后仍原子保存 Token 并清除 pending。
- Device 授权与 Token 刷新共用 task + profile 的非阻塞凭证操作锁，避免两类操作并发保存时互相覆盖；并发调用只允许一个进程进入 Token Endpoint，锁竞争直接返回 `busy`。
- Skill 明确要求展示授权链接后立即结束当前轮次，只有收到新的用户确认消息后才能执行一次 complete；OAuth 请求完全排除在合同读请求重试规则之外。
- `auth status` 明确展示 pending、uncertain、denied、expired 和 restart_required，不支持自动附加 `--output`。
- 收紧 WorkBuddy 授权展示契约：`auth init` 的最终回复必须同时包含可点击的 `verification_uri_complete`、`qr_code_path` 对应 PNG 产物和过期时间，禁止只返回链接或只返回二维码。二维码附件失败时仍返回链接并明确告知展示失败，不得虚假声称已展示。

### 本轮验证

- 覆盖重复 init、显式 restart、请求前 checking、超时 uncertain、终态保留、并发 complete、旧 pending 状态兼容和 Skill 文案约束。
- `go test ./internal/cli ./internal/credential ./internal/oauth -count=1`、`go test ./... -count=1`、`go mod verify` 和 `git diff --check` 通过。
- 本需求授权状态机、凭证锁和 CredentialStore 的定向 `-race` 测试通过。
- 完整 `go test -race ./internal/cli ./internal/credential -count=1` 仍会命中存量 `App.updateNotice` 并行测试竞态（`internal/cli/app.go:171`）；该点不在本次授权改动路径上，本轮不通过修改高扇出全局状态规避。

## 2026-08-06 WorkBuddy 二维码内联展示修复

- `auth init --output json` 的 `pending` 响应新增可选 `qr_code_data_uri`，格式为 `data:image/png;base64,...`；原有 `qr_code_path`、授权链接和过期时间保持不变。
- CLI 使用现有 `go-qrcode` 一次生成 PNG 字节，同一份内容同时写入 0600 权限的二维码文件并编码为 Data URI，避免两套二维码内容不一致。
- Data URI 只作为本次 CLI JSON 输出使用，不写入 CredentialStore、profile 或日志；非 `pending` 状态不输出该字段。
- WorkBuddy 的二维码文件名摘要包含 `CODEBUDDY_SESSION_ID`，不同任务即使使用同一个 profile 也不会互相覆盖 PNG 降级产物；豆包会话工作区路径规则保持不变。
- WorkBuddy Skill 在 `auth init` 前加载展示说明，随后使用 `show_widget` HTML 模式内联展示二维码，并在正文中继续提供可点击授权链接和过期时间。
- `show_widget` 失败时才降级为授权链接和 PNG 产物卡片，并明确提示内联展示失败；不得虚假声称二维码已展示。
- 豆包继续使用原有 `qr_code_path` 展示 PNG，不依赖 WorkBuddy 的 `show_widget`。
- 授权状态机保持不变：展示完成后立即结束当前轮次，仍需收到用户新的“已授权”消息后才允许执行一次 `auth complete`。

### 本轮验证

- 新增 Data URI 格式、PNG 尺寸、Data URI 与落盘文件字节一致、pending 复用、WorkBuddy 跨任务二维码路径隔离、非 pending 状态不输出图片数据及 Skill 展示边界测试。
- `go test ./internal/cli -count=1`、`go test ./... -count=1`、`go mod verify` 和 `git diff --check` 通过。
- 本需求 `go test -race ./internal/cli -run 'TestAuthDeviceInit|TestDeviceAuthSkillsEnforceWorkBuddyTurnBoundaries' -count=1` 通过。
- 完整 `go test -race ./internal/cli -count=1` 仍命中存量 `App.updateNotice` 并行测试竞态，堆栈位于 `internal/cli/app.go:171`，与本次二维码生成和 Skill 展示路径无关。

## 2026-08-06 WorkBuddy 授权文案与二维码尺寸优化

- WorkBuddy 授权提示改为面向用户的固定语义，不再展示 `user 身份未授权`、CLI 命令或内部状态。
- 回复顺序统一为授权说明、可点击链接、扫码提示、过期时间和“已授权”后续指引。
- 二维码 PNG 仍保持 320×320 原始清晰度，`show_widget` 中的页面显示尺寸限制为 220×220，减少占屏同时保留扫码所需的完整边界。
- 共享 Skill 同步引用同一展示契约，豆包展示、Device 状态机、Token 和服务端接口均不变。
- 新增静态回归断言，锁定新文案、220×220 显示尺寸，并防止旧技术化文案和 320×320 过大显示回归。

## 2026-08-06 WorkBuddy 授权时间与位置文案修正

- `auth init` 和授权成功输出新增 `expires_at_display`，固定为 `YYYY-MM-DD HH:mm（北京时间）`；原 `expires_at` RFC3339 字段保留，不破坏旧调用方。
- WorkBuddy 最终回复只展示 `expires_at_display`，禁止把带 `T`、`+08:00` 和秒数的机器时间直接呈现给用户，也不使用代码样式渲染时间。
- 授权文案取消“上方二维码”和“下方二维码”等位置假设，统一为“也可以使用手机扫描二维码”，适配 WorkBuddy 工具卡片的实际排版。
- 增加新建与复用 pending 会话的北京时间输出回归，非 pending 状态不输出授权截止时间展示值。

## 2026-08-06 WorkBuddy 授权三步文案定稿

- 按最终业务文案先展示“继续查询前，需要先完成智书合同授权。”，再展示三步编号列表：点击蓝色授权链接或扫描本消息中的二维码、在截止时间前完成手机号验证与企业确认授权、授权后回复“已授权”并继续合同查询；不使用“上方/下方”等依赖客户端布局的相对方位。
- “已授权”固定使用 Markdown `**已授权**` 加粗，逗号不纳入加粗范围。
- `expires_at_display` 最终格式调整为 `YYYY-MM-DD HH:mm:ss`，不显示 RFC3339 的 `T` 和 `+08:00`；原 `expires_at` 继续保留供机器读取。
- 新增静态 Skill 断言锁定三步文案、链接标题、加粗时间和加粗“已授权”；新建与复用 pending 会话的展示时间均包含秒。

## 2026-08-07 auth init 连接前安全重试

- `auth init` 首次调用只有在错误链明确包含 TCP `dial` 失败、能确认 HTTP 请求尚未发出时，才由 CLI 内部自动重试一次。
- 两次 TCP `dial` 均失败时立即返回最终错误，不增加第三次请求、循环重试、轮询或通用重试框架。
- 请求已发送后的读超时、HTTP 5xx、响应中断、响应解析失败、上下文取消或截止均不自动重试。
- `auth complete`、Refresh Token 和 Revoke 请求继续保持不自动重试，Device 授权状态机和服务端接口不变。
- Skill / 模型层仍只允许执行一次 `auth init` 命令；命令最终失败后禁止自行执行 `curl`、`auth status` 或其他探测，必须如实告知用户并等待新的重新授权确认。

## 2026-08-07 豆包 macOS 本地 Skill 适配

- 新增豆包本地运行环境识别：没有 `SKILL_SESSION_WORKSPACE` 时，使用 `DOUBAO_SESSION_ID` 作为当前对话的稳定凭证命名空间。
- 豆包本地 Device Token 使用系统安全存储；macOS 落入 Keychain，不写入普通配置文件或明文 Token 文件。
- `DOUBAO_TASK_ID` 仅代表单次工具调用，不参与凭证、授权锁或二维码文件隔离，保证同一对话多轮调用复用同一授权。
- 豆包本地与 WorkBuddy 使用独立的 Keychain 账户名前缀和文件锁命名空间；即使两端出现相同 session 字符串也不会串用凭证。
- WorkBuddy 既有 Keychain 账户名和锁命名空间保持不变，避免现有已授权任务失效。
- AgentKit 云端豆包继续优先使用 `SKILL_SESSION_WORKSPACE` 和 `CONTRACT_CLI_CREDENTIAL_KEY_V1` 加密文件，不改变云端沙箱重建恢复逻辑。
- 同时出现 `DOUBAO_SESSION_ID` 和 `CODEBUDDY_SESSION_ID` 时失败关闭，禁止猜测当前本地 Agent。
- 豆包本地二维码文件名按 `DOUBAO_SESSION_ID + profile` 隔离，二维码正文仍必须与完整 HTTPS 授权链接和过期时间同时展示。
- 移除 `auth` 和 `contract-cli-shared` FrontMatter 中社区 SKILL.md 校验不允许的 `version` 字段，保证豆包本地扫描器可按标准元数据加载。

## 2026-08-07 豆包本机冗余适配清理

- 根据客户实际使用链路完成纠偏：豆包工作任务使用云端会话工作区和加密凭证，WorkBuddy 使用客户本机任务标识和操作系统安全存储。
- 删除建立在错误运行假设上的豆包本机环境识别、Keychain 账户前缀、二维码命名空间和授权锁命名空间分支；相关本机环境变量不再启用 Device 凭证能力。
- 豆包云端继续优先使用 `SKILL_SESSION_WORKSPACE` 和 `CONTRACT_CLI_CREDENTIAL_KEY_V1`，保留 AES-256-GCM 密文格式、凭证路径和 profile 快照恢复，不迁移现有数据。
- WorkBuddy 的 `CODEBUDDY_SESSION_ID`、`sessionID:profile` 安全存储账户名、二维码摘要和授权锁摘要保持不变，现有已授权任务无需重新授权。
- Device Grant 命令、状态机、Token、二维码 JSON 输出、业务请求、旧 OAuth、App 身份、安装器和发布流程均未修改。
- 授权 Skill 删除豆包本机与平台强制结构化工具的错误断言，保留链接与二维码同时展示、单次 `auth init` / `auth complete` 和 OAuth 禁止自动重试等交互约束。

## 2026-08-07 正式发布环境收敛

- 正式 CLI 恢复为只内置 `prod` 环境预设，删除测试阶段临时加入的 `dev` 入口和开发环境服务地址。
- `config add --env dev` 在任何网络请求前明确拒绝；不做 dev 到 prod 的自动迁移，避免开发环境 Token 与生产服务混用。
- CLI 帮助、Device profile 恢复提示和 Auth Skill 统一只引导创建生产 profile；默认环境和默认 profile 仍分别为 `prod`、`contract`。
- 授权 Skill 恢复定稿的二维码位置无关文案，明确豆包云端工作任务与 WorkBuddy 客户本机的真实运行边界，并删除平台强制结构化工具的错误断言。
- Device Grant 状态机、CredentialStore、Token 刷新、业务请求、旧 OAuth 和 App 身份行为保持不变。
- 本地 npm 安装校验改为给临时构建的当前版本制品生成独立 checksum，不再依赖工作区里恰好存在同版本历史制品，保证后续正式版本号升级仍能验证安装器验签链路。

## 2026-08-09 豆包普通工作任务授权支持

- 在 AgentKit 和 WorkBuddy 之后新增豆包普通工作任务运行时：仅当前两种平台标识都缺失时，才使用 `SESSION_ID` 识别当前任务，不改变原有优先级。
- 不使用当前环境中指向无效目录的 `WORKSPACE`；以真实当前工作目录和 `SESSION_ID` 摘要建立任务级凭证、二维码和文件锁目录。
- Device Token、Refresh Token、device code 和 profile 快照继续使用 AES-256-GCM 密文及原子写入；原始 `SESSION_ID` 不进入文件名、日志或 CLI 输出。
- 同一任务多轮调用可复用 pending 会话和 Token，新任务使用新命名空间并必须重新授权；临时 HOME 丢失时可从同任务密文快照恢复 Device profile。
- 普通工作任务没有平台 Secret API 或 Keyring；该密文方案用于避免明文落盘和正常对话泄露，不宣称能抵御同沙箱内恶意 Shell。AgentKit 显式密钥和 WorkBuddy 操作系统安全存储行为保持不变。

## 2026-08-09 豆包普通工作任务二维码产物交付修复

- 根据豆包普通工作任务实测结果，将该运行时的二维码 PNG 从隐藏会话目录调整到任务初始工作目录根部，使用 `contract-cli-device-auth-<摘要>.png` 可见文件名，便于平台识别并交付图片产物。
- 文件摘要继续包含任务命名空间和 profile，不使用原始 `SESSION_ID`；不同任务和不同 profile 不共享二维码文件。
- 凭证密文和授权锁继续保留在 `.contract-cli/sessions/<摘要>` 隐藏目录，未改变 Token、pending、profile 快照或任务隔离规则。
- 授权 Skill 强制要求将 `qr_code_path` 作为图片附件或图片产物交付；只有对话中实际出现图片缩略图或产物卡片时，才能声称二维码已经展示。
- AgentKit 的二维码路径以及 WorkBuddy 的二维码路径、文件名摘要和 `show_widget` 展示逻辑保持不变。

## 2026-08-09 豆包普通工作任务授权展示降级

- 根据普通工作任务实测结果，本地 PNG 无法通过稳定的平台接口转换为对话图片产物，因此该运行模式授权时只展示完整可点击 HTTPS 链接和过期时间，不再尝试复制、重新编码或交付二维码。
- 授权 Skill 禁止普通工作任务读取或处理 `qr_code_path`、`qr_code_data_uri`，也禁止为二维码调用代码执行、图片处理或图片交付工具；返回链接后立即结束当前轮次。
- WorkBuddy 继续通过 `show_widget` 展示二维码，AgentKit 继续使用平台的 `qr_code_path` 能力，二者行为不变。
- CLI 继续生成并返回 `qr_code_path` 和 `qr_code_data_uri`，保持 JSON 接口、Device Grant 状态机和后续恢复二维码展示的兼容性。

## 2026-08-10 统一授权扩展至智审平台

- 生产环境 Device profile 的固定授权范围调整为 `contract:full contract-review:full`，顺序固定；`auth init` 将该组合 scope 原样提交给统一授权服务。
- 继续复用既有 Device client、business type、resource、Token 存储和刷新流程，不新增 Token 导出、通用 HTTP 调用或智审业务命令。
- Auth Skill 补充说明一次授权同时包含合同与智审平台访问范围，但不提供智审业务操作说明或智审 Skill。
- WorkBuddy、豆包普通工作任务和 AgentKit 的运行时识别、凭证隔离、授权展示与单次 `auth init` / `auth complete` 状态机保持不变。

## 2026-08-18 智书 Skill 凭证与接口调用边界收紧

- 智书共享 Skill 禁止无业务边界地枚举或批量调用全部接口；用户未明确业务目标、环境、允许范围和读写类型时只做澄清，不执行授权探测或业务命令。
- Agent 不得主动索要或接收 Token、Refresh Token、AK/SK、Cookie、Session、App Secret、device code、密码等原始敏感凭证；用户主动发送时不复述、不使用，并提示立即撤销或轮换。
- user 身份缺失时继续使用既有 Device Grant；app 身份保留原 CLI 能力，但 Agent 只能使用用户已在本机安全配置的环境变量或 CredentialStore，不得让用户在对话中粘贴 App Secret。
- 接口文档不能替代明确的业务范围和调用授权，未覆盖接口继续明确为暂不支持，不回退到未开放的通用 `api call`。
- 本次只调整智书 Skills 与静态回归测试，不修改 CLI 公共命令、OAuth、CredentialStore、开放平台接口或全局 Agent 安全策略。

## 2026-08-20 WorkBuddy 授权二维码改为 PNG 附件

- 根据 WorkBuddy 实测，模型将 `qr_code_data_uri` 复制到内联展示工具时可能截断 Base64，导致图片数据损坏和二维码无法加载。
- WorkBuddy 主路径改为将 `qr_code_path` 对应的原始 PNG 通过 `present_files` 交付为图片附件/产物卡片，不再让模型复制或重新编码 Base64。
- 图片交付失败时不重试、不调用其他图片处理工具，仍返回完整可点击授权链接和过期时间。
- CLI 继续返回 `qr_code_path` 和 `qr_code_data_uri`，二维码内容、文件权限、命名空间、Device Grant 状态机、Token 和 CredentialStore 行为均不变。
- 豆包 AgentKit 继续使用平台 `qr_code_path` 能力，豆包普通工作任务继续只展示授权链接，两者不受影响。

## 2026-08-24 正式发布版本更新为 1.8.0

- 将 npm 正式包版本从已发布的 `1.7.0` 更新为 `1.8.0`，用于发布本次豆包、WorkBuddy 统一用户授权能力。
- 本次只调整发布版本元数据，不修改 Device Grant、CredentialStore、Skills、业务命令或各环境访问边界。

## 2026-08-24 正式包非生产残留状态防护

- 正式 CLI 新增集中式生产 profile 校验，授权、身份切换、Token 刷新、注销及所有业务请求均在发送网络请求前拒绝非 `prod` 环境或非生产服务地址。
- 对 Device profile 快照、待兑换会话的 Token Endpoint 和完整授权链接同步执行生产边界校验，历史非生产凭证不兑换、不刷新、不撤销也不自动迁移。
- 用户显式使用同名 profile 重新执行生产 `config add` 时，清除当前运行环境中遗留的 user/app 身份、App Secret 和 Device 凭证，再保存干净的生产配置。
- 正式包 HTTP Transport 增加最后一道防线，直接拒绝已知的非生产开放平台和授权域名，防止未来新调用链遗漏 profile 校验。
- 正式 Skills 固定使用 `contract` profile 和 `prod` 环境，用户 Prompt 不能覆盖该规则；WorkBuddy 更新 Skill 后必须完全退出并新建任务，旧任务不会热加载新规则。
- npm 正式包版本更新为 `1.8.3`；本次不修改 org-v2、org-fe、open-platform、Nacos 或其他服务端配置。
