# 审查正文来源选择

## 2026-09-20

- 分支 `test-feature-contract-cli-add-review`，仅增量修改审查 Skill；保留已有未提交的测试命令预览和引用参考。
- 明确同一合同可有协商正文与任务中心正文，两者可以不同。双入口可信候选存在且用户未选择时，上传前先选择来源；已有明确选择不重复确认。
- 文件标识和版本从同一来源记录成组获取，不混用 `cooperation_file_id + version` 与 `file_id + text_id`。服务端拒绝后不自动切换来源或身份；用户更换文件需重新提取主体及确认立场。
- 不把状态、pc_url 路径、可读权限或上传成功当作审查授权；关联不足报告缺口。结果链接仍沿用合同服务 pc_url，不新增 URL 推断。
- 不修改 CLI 命令/API、历史 rule table、认证、环境配置或审查结果导航。

## 验证与边界

- `go test ./...`、`go vet ./...`、Skill `quick_validate.py`、`git diff --check` 通过。已有可执行引用上传/主体/导航/分发测试作为参数契约验证。
- 对双入口未选、已选任务正文、缺关联、服务端拒绝四类情况做人工离线规则检查；不是独立 Agent 测试或真实宿主端到端验证，不宣称自动行为测试闭环。
- 未更新本地安装版本、提交、推送或发布。需先部署 CLM 来源分流修复，再通过既有构建和 Skill 分发流程更新使用端；新 Skill 遇旧后端拒绝时不得换来源重试。
- 回滚仅移除本次来源选择增量，保留原测试命令预览。过程证据位于 `/Users/cuizongbao/.local/share/qfei-artifacts/clm-1f6964c4bacb/review-source-routing/20260920-side/`。
