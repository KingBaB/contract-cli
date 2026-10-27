# test 环境 CLI 支持变更记录

## 2026-10-10

- 变更摘要：将当前分支的环境预设调整为 prod/test，默认仍为 prod。
- 关键逻辑/决策：`config add --env test` 写入 `https://test-open.qtech.cn`、`https://test-myaccount.qtech.cn` 和 test app token endpoint；profile 校验按声明环境做同源校验，避免 prod/test 跨环境复用端点或凭证；传输层继续阻断历史 dev/blue 环境域名。
- 影响范围：配置初始化、profile/Device 快照校验、app 登录提示、Open Platform 请求上下文、帮助文案、README、命令文档和相关 Skill 环境说明。
- 验证：新增/更新 test profile、test 请求上下文、test Device 快照恢复、dev/blue 拦截、帮助文案和 Skill 安全边界测试。


## 2026-10-27 预审测试分支使用 test 环境

- 变更摘要：在已合入预审功能的 `test-feature-cli-support-pre-audit` 分支，将默认且唯一环境从 prod 切换为 test；默认 profile 名保持 `contract`。
- 配置范围：开平 API、OAuth 资源、账号元数据、app token 端点及现有 profile/Device/HTTP 环境校验；复用此前 test 预设的 Device client 与 scope。
- 隔离规则：旧 prod profile、Device 快照和 pending 不直接使用；显式重新初始化同名 test profile 时，沿用现有逻辑清理旧 Token、App Secret 与 Device 凭据，再重新授权。传输层允许 test，阻断已知 prod/dev/blue 域名及其跳转。
- 配套调整：同步帮助、README、命令参考、内置 Skill、构建说明和相关测试夹具；现有内部 production 前缀辅助函数名称保留，避免扩大重命名范围。
- 验证策略：隔离配置目录与 mock HTTP，先观察 test 默认配置、环境隔离、文案和构建断言的 RED，再执行定向与完整回归、静态检查和构建；不真实登录或发起业务审查。
- 验证结果：原基线通过；配置、网络、指引及构建约束均观察到有效 RED。完成后 `go test ./...`、Node Skill 契约、安装器/Skill 同步测试、`build-flags.sh`、`go vet ./...`、宿主构建及 CLI smoke 均通过。完整回归中发现的旧摘要 Skill 断言和内置参考文档的 prod 链接已同步修正。
