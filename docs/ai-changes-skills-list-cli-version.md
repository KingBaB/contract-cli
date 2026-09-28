# AI 变更记录：skills list 展示 CLI 版本

## 2026-05-26 skills list 版本口径统一

变更摘要：`contract-cli skills list` 中每个内置 skill 展示的版本改为当前 CLI 运行时版本，与 `contract-cli --version` 保持一致。

涉及文件/模块：

- `internal/cli/skills_command.go`
- `internal/cli/app_test.go`
- `docs/app-command-development-guide.md`

关键逻辑/决策：

- `skills list` 不再使用 `SKILL.md` front matter 里的 `version` 作为展示版本。
- 展示版本统一从 `internal/build.Current().Version` 获取，和 `contract-cli --version` 同源。
- `SKILL.md` 的源文件 `version` 字段保留为 skill 文档内部元数据。
- `skills install` 安装时会将目标目录中 `SKILL.md` front matter 的 `version` 改写为当前 CLI 运行时版本，避免 agent 查看已安装文件时看到第二套版本口径。
- 不新增 `cli_version` 或 `skill_version` 字段，统一对外展示和安装结果中的版本口径。

验证策略：

- 更新 `skills list` 测试，临时设置 CLI 版本为 `1.2.3`，确认所有 skill 行都展示 `1.2.3`。
- 同一测试确认 fixture 中的 `SKILL.md version` 不再出现在列表版本列。
- 保留隐藏禁用 skill 的断言，避免 `contract-cli-api-call` 被展示。
- 执行 `go test ./...` 验证整体回归。

## 2026-06-02 skills install 改写安装版本

变更摘要：`contract-cli skills install` 在复制内置 skill 时，将目标目录内每个 `SKILL.md` front matter 的 `version` 改写为当前 CLI 运行时版本。

涉及文件/模块：

- `internal/cli/skills_command.go`
- `internal/cli/app_test.go`
- `docs/app-command-development-guide.md`

关键逻辑/决策：

- 只在安装输出阶段改写 `SKILL.md`，不修改仓库源文件里的 skill 文档内部版本。
- 非 `SKILL.md` 文件保持原样复制，已存在 skill 且未传 `--force` 时继续跳过，不改写用户本地已有文件。
- `--force` 覆盖安装时会重新复制并改写 `SKILL.md` 版本。

验证策略：

- 更新 `skills install`、`skills install --force` 和默认 Codex Home 安装测试，确认安装后的 `SKILL.md version` 等于测试注入的 CLI 版本。
- 执行 `go test ./internal/cli`、`go test ./...` 和 `make release-check` 验证整体回归。
