# AI 变更记录：预审开发分支使用 dev 环境

## 2026-10-27

- 变更摘要：在已包含预审能力的 `dev-feature-cli-support-pre-audit` 分支，将默认且唯一环境从 prod 切换为 dev，默认 profile 名保持 `contract`。
- 配置范围：开平 API、OAuth 资源、账号元数据、app token 端点，以及既有 profile/Device/pending/HTTP 环境校验；当前 dev 配置已注册相同的公开 Device client 和合同、审查 scope，继续使用这些值。
- 隔离规则：prod/test/blue profile、Device 快照和 pending 在使用前拒绝；用户显式重新初始化同名 dev profile 时，沿用旧认证清理逻辑，清理原 Token、App Secret 和 Device 凭据后重新授权。传输层允许 dev，阻断已知 prod/test/blue 域名及其跳转。
- 配套调整：同步帮助、README、命令参考、内置 Skill/参考文档、构建环境说明及测试夹具；保留现有内部 production 前缀辅助名称以控制变更范围。
- 实现边界：仅适配 CLI 开发分支环境，不涉及预审业务接口和后端改造；本地配置由操作者显式重新初始化，不自动安装或发布。
- 验证策略：隔离配置目录和 mock HTTP，先观察 dev 默认配置、环境限制、文档指引与构建约束的有效 RED，再执行定向和完整 Go/Skill/安装器回归、静态检查、构建及 CLI smoke。
- 验证结果：基线通过；dev 默认/显式配置、prod 凭据清理、HTTP 主机限制、Skill/帮助与构建约束均观察到有效 RED。修改后定向和完整 `go test ./...`、Node Skill 契约、安装器/Skill 同步测试、构建 flag 检查、`go vet ./...`、宿主构建、CLI smoke 和 `git diff --check` 均通过。
