# AI 变更记录：每次执行检查版本更新

## 2026-05-26 自动版本检查去除 30 分钟节流

变更摘要：将交互终端下的自动版本检查从“最多每 30 分钟检查一次”调整为“每次符合条件的普通命令都检查一次”，发现新版本时每次都提示升级命令。

涉及文件/模块：

- `internal/cli/update_command.go`
- `internal/cli/app.go`
- `internal/update/update.go`
- `internal/cli/app_test.go`
- `internal/update/update_test.go`
- `README.md`
- `docs/cli-command-reference.md`
- `docs/cli-test-plan.md`

关键逻辑/决策：

- 自动检查入口不再读取 `update-check.json` 判断缓存 freshness，也不再通过 30 分钟间隔抑制请求和提示。
- 自动检查失败仍不阻断原命令，但不再写入失败缓存；下一次符合条件的命令会重新尝试远端检查。
- 成功检查后仍写入 `update-check.json` 作为最近检查结果记录；该缓存不再影响后续自动检查频率。
- 保留既有跳过条件：`help`、`version`、`update` 命令不触发自动检查，非交互终端不触发，`CONTRACT_CLI_NO_UPDATE_CHECK=1` 可关闭。

验证策略：

- 更新自动提示测试，连续两次 `skills list` 应发起两次远端请求，并在发现新版本时两次都输出升级提示。
- 更新失败场景测试，连续两次网络失败应发起两次远端请求，且不输出升级提示、不阻断原命令。
- 保留环境变量关闭自动检查和手动 `update check` 的既有测试。
- 执行 `go test ./...` 验证整体回归。
