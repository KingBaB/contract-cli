# 审查 Skill 临时测试命令预览（已移除）

## 2026-09-22 移除

- 按用户要求完全移除临时命令预览：删除审查 Skill 的临时规则小节及专用参考文件。
- 现有业务正文、格式修复、CLI 命令和环境限制保持不变；下方为历史记录，不再是活动规则。
- 新增嵌入分发回归，确保 Skill 正文与资源均不再分发临时预览。

## 2026-09-20

- 分支：`test-feature-contract-cli-add-review`。
- 修改 `skills/contract-cli-review/SKILL.md`，新增 `references/test-command-preview.md`。仅实际环境为 test 且开放平台 origin 为 `https://test-open.qtech.cn` 时启用，不凭 profile 名称判定。
- 每次命令前先输出命令、对应开放平台接口和参数，覆盖帮助、本地辅助、预检、合同取数、上传、主体、清单、发起、分页和显式轮询。无远端请求时明确标注；内部编排说明涉及接口和运行时参数来源，不宣称拦截每个 HTTP 请求。
- 命令及参数两处都隐藏凭据、会话和 URL 签名，不输出合同正文/文件字节。仅在测试进度消息允许必要业务 ID 和字段，不改变正式候选列表和最终审查结果模板。
- 生产、其他或未确认环境禁止详细预览。迁移/发布生产 Skill 前必须移除测试入口段落和参考文件，不将此能力作为生产功能发布。该限制为 Skill 行为约束，不是新增 CLI 日志开关或生产发布流水线门禁。
- 不修改 CLI 命令、请求参数、用户权限、确认/重试规则、独立清单规则维护 Skill，不安装或更新本地二进制及 Skill 副本。

## 验证与交付边界

- 修改前离线场景可复现：test 状态查询没有先展示命令、接口和参数；行为断言失败后才编辑 Skill。
- 修改后五个独立离线场景通过：test 状态查询先预览、prod 不启用、首次帮助无远端接口、签名 URL 展示脱敏但执行参数不被替换、环境未知不凭 profile 名开启。
- `go test ./...`、`go vet ./...`、Skill `quick_validate.py`、`git diff --check` 通过；系统 Python 缺少 PyYAML，改用已有项目虚拟环境校验，未安装依赖。
- 离线模拟只验证约定场景，不证明所有 Agent 或宿主都能机械遵守；未执行真实上传或发起审查。
- 过程记录在仓外 `/Users/cuizongbao/.local/share/qfei-artifacts/contract-cli-770819e9cd80/review-test-command-preview/20260920/`。
- 本次仅修改项目源码；后续需构建并通过既有 Skill 分发流程刷新使用端才能使用，未提交、推送或发布。
