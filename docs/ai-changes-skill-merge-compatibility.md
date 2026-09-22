# Release 合并后的 Skill 兼容修复

## 范围

基于 `test-feature-contract-cli-add-review` 的合并提交 `19c7d8b`，修复已确认的 Skill 冲突和两处编译阻断。不切换分支、不发布、不更新全局安装。

## 修改

- 合并 `contract-cli-contract`、`contract-cli-shared` 的重复 description；同时保留 release 新增的搜索、人员、部门能力及已有审查入口。
- 合同搜索只路由到独立搜索 Skill，去掉共享 Skill 中重复的旧路由和模块条目。
- auth/shared/contract 中授权检查、授权与业务操作沿用同一 test profile；明确 app 任务不先走 user 授权、新 test profile 默认 OAuth、Device 需要 test 专属 client。
- 本地 Device 复用与恢复的两处旧方法名接回现有 `loadEnvironmentDeviceCredential`，保留凭证环境校验。
- 修正合并带入的本地授权正向测试中生产地址夹具，以及 Skill 安装数量断言（13 → 16）。不修改环境拒绝测试。
- 增加源文件和安装产物均可运行的 Skill 元数据、路由、profile 一致性测试，并纳入 `npm test`。

## 保留

- test-only、prod 暂时禁用、dev 屏蔽。
- 历史 CLI 身份规则、审批矩阵 `rule table`、审查 user-only。
- 原有未提交审查 Skill 及两个变更摘要不变；审查文件 SHA-256 仍为 `18b138cbdcc9e52d4c28d4771d77dc808c4c9b1b1ba942a6bd1ba4ac017c1960`。

## 验证

- 新增测试先复现重复 description、搜索双路由、profile 混用；独立复核再发现自定义 profile 被固定名称覆盖、app 示例遗漏 profile，补充测试修复后 6/6 通过。
- `npm test`、`go vet ./...`、`go build ./cmd/contract-cli`、`git diff --check` 通过。
- 隔离目录安装 16 个 Skill 后，同一组 6 项一致性检查通过；未操作用户真实凭证或在线业务。
- 合并导致的编译失败单列为基线故障，不将其冒充行为测试 RED。

## 未纳入本次范围

标准 Skill 校验器接受 contract Skill，但拒绝 auth/shared 历史顶层 `version` 字段；CLI 安装器仍使用原有版本写入协议。此兼容问题非此次合并引入，本次没有宣称所有 Skill 均通过标准格式校验。后续如统一到 `metadata.version`，须同时修改安装器和相应兼容测试，不能只改 Markdown。

独立复核另提示 Codex 的 Device pending 二维码展示说明不如 WorkBuddy/豆包完整；本次不新增宿主工具或二维码强制展示协议，保留已有链接授权与跨轮次完成边界。

过程证据：`/Users/cuizongbao/.local/share/qfei-artifacts/contract-cli-770819e9cd80/skill-merge-compatibility/20260922/`。

回滚仅反向应用本次改动，不重置工作树，不覆盖已有未提交审查修改；没有数据迁移或外部写操作需要撤销。
