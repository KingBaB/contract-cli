# AI 变更记录：飞书式更新提示缓存

## 变更摘要

- 将普通业务命令的自动升级提示改为飞书式模型：优先读取本地 `update-check.json` 决定是否注入 `_notice.update`，缓存缺失或过期时在当前命令内短超时刷新远端版本缓存。
- 保留手动 `contract-cli update check` 的同步远端检查行为。
- 保留 24 小时 TTL；cache fresh 时不访问 registry，因此刚发布的新包不会在当前命令中立刻提示。

## 关键逻辑

- 复用 `LoadCache`、`CacheFresh` 和 `NoticeFromCache` 做 fresh cache 判断与 notice 生成。
- cache 缺失、channel 不匹配或 stale 时，`maybePrepareUpdateNotice` 在 `1500ms` 超时内同步刷新远端并写入 cache。
- 刷新失败不阻断业务命令、不写失败缓存、不注入 stale notice。

## 验证

- 覆盖缓存可升级时当前命令立即注入 `_notice.update`。
- 覆盖 stale 缓存同步刷新并写入新版本。
- 覆盖 fresh 无更新缓存不访问 registry、不立即提示新包。
- 覆盖刷新失败不阻断业务命令。

## 2026-06-03 修复：对齐飞书当前同步刷新实现

- 变更摘要：修复纯 goroutine 后台刷新在真实 CLI 进程退出时可能丢失的问题，改为对齐飞书当前实现。
- 关键逻辑：普通命令先读取本地 `update-check.json`；cache fresh 且 channel 匹配时只从缓存注入 `_notice.update` 并跳过远端；cache 缺失、channel 不匹配或过期时，当前命令在 `1500ms` 超时内同步刷新远端版本并写入 cache。
- 失败策略：远端刷新失败不阻断业务命令、不注入 stale notice、不写失败 cache，仅记录 debug 日志。
- 验证策略：覆盖 fresh cache 有/无更新、stale cache 成功刷新、stale cache 刷新失败、cache 缺失成功写入，以及自动检查关闭/CI/help/version/update check 跳过场景。
