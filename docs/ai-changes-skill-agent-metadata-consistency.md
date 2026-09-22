# AI 变更记录：公开 Skill agents 元数据一致性

## 2026-09-22

- 为 `contract-cli-review` 和 `contract-cli-review-rule` 补齐 `agents/openai.yaml`，分别声明审查任务与审查配置的展示信息、默认调用提示和隐式调用策略。
- 将公开内置 Skill 的结构约定写入 `docs/app-command-development-guide.md`：含 `SKILL.md` 的分发 Skill 必须同时提供 `agents/openai.yaml`；不含 `SKILL.md` 的禁用内部描述保留例外。
- 新增元数据一致性测试，扫描所有公开 Skill，校验 `agents/openai.yaml` 的基本字段，并校验两个审查 Skill 的相互路由。
- 不修改 CLI 命令、接口、权限、身份、环境校验、审查正文规则或 `contract-cli-api-call` 禁用策略。

## 验证

- `go test ./internal/cli -run 'TestPublicSkillAgentMetadata|TestReviewSkillAgentMetadata' -count=1`
- 相关 Skill 回归测试、`npm pack --dry-run --json --ignore-scripts` 和 `git diff --check`
