# AI 变更记录：升级提示 JSON notice 与 24 小时缓存

## 2026-05-27 飞书式升级提示改造

变更摘要：将自动升级提示从交互终端 stderr 文本提示改为 JSON 输出中的 `_notice.update`，并将自动远端检查改为 24 小时缓存节流。

关键逻辑/决策：

- 普通命令不再输出旧的 `A new contract-cli version is available` stderr 文本。
- JSON object 输出会在发现新版本时注入 `_notice.update`，其中 `command` 继续使用 npm 安装命令。
- `update-check.json` 作为 24 小时缓存：fresh cache 不再请求 npm registry，但仍可用于注入 notice。
- `contract-cli update check` 调整为飞书式手动校验：默认输出文本，带 `--json` 时返回顶层 `ok/action/current_version/latest_version/message` 等结构；手动检查不注入 `_notice.update`。
- `--raw`、yaml、table、纯文本命令不注入 `_notice.update`；CI 环境跳过自动远端检查。
- `contract-cli-shared` skill 的触发描述补充 “看到 JSON 输出中的 `_notice` / `_notice.update`”，对齐飞书 `lark-shared` 的 agent 触发方式。
- `contract-cli-shared` skill 补充普通查询不要默认加 `--raw`；说明 `--raw` 会绕过 JSON renderer，不会注入 `_notice.update`。

测试覆盖：

- 新增 update cache freshness、cache notice 生成测试。
- 新增 renderer 注入、合并和跳过 notice 的测试。
- 更新 CLI 自动检查测试，覆盖 JSON notice、fresh cache、禁用环境变量、手动 `update check` 默认文本输出和 `--json` 飞书式结构输出。

## 2026-05-28 review 修复：输出流与过期缓存

变更摘要：修复本地 review 发现的两个升级提示边界问题，确保手动检查正常结果走 stdout，并避免过期缓存刷新失败时注入旧 `_notice.update`。

关键逻辑/决策：

- `contract-cli update check` 默认文本结果改为写入 stdout；stderr 只保留日志和错误流，便于脚本重定向正常输出。
- 只有 fresh cache 才允许生成缓存 notice；过期 cache 必须先尝试刷新，刷新失败时不注入旧 `_notice.update`。
- 远端刷新成功后仍按最新结果重建 notice；无更新或跳过检查时不注入 notice。

测试覆盖：

- 新增 `update check` 默认文本输出的 stdout 回归测试，确认正常结果不写 stderr。
- 新增 stale cache + registry failure + JSON 业务命令测试，确认不会注入过期 `_notice.update`。

## 2026-05-28 13:00 qfei-code-review
- 变更摘要：已使用 `qfei-code-review` 完成本地代码审查，结论：Approve。
- 涉及文件/模块：当前 MR 或本地代码 diff（已排除 AI 变更记录文件）。
- 关键逻辑/决策：按共享 review 规则检查正确性、回归风险、测试覆盖；线上机器人按下方凭证判断是否触发 review。

<!-- qfei-code-review
record_version: 1
requirement: update-notice-json-cache
base_sha: 5a39986f6667e100ef112e509ff4381ba89b8fc3
start_sha: 5a39986f6667e100ef112e509ff4381ba89b8fc3
head_sha: f7fa1bf28b089d51028b4c9067cd8447f3559db6
diff_hash: sha256:4d7429c1a1a05c304a95ac67270e59e508bdf73fbd939c2956dace402dced505
reviewer: wanghengju@qfei.cn
reviewed_at: 2026-05-28T13:00:46+08:00
review_status: completed
verdict: Approve
-->
