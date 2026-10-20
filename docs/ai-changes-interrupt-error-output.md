# 中断取消与结构化错误出口

## 2026-10-20 补齐文本取消的升级提示

- 用户复现 `auth login --profile trace-dev-app --as user` 中断后缺少升级提示；模拟缓存只准备通知，旧 Execute 文本错误分支未显示它。
- TDD：先用真实 OAuth 回调等待取消链路、模拟 registration Transport 复现默认/debug 两个红例；修复仅在文本错误摘要后输出本次已有的非空 `notice.update.message`，没有新网络请求，也没有新增锁。
- 无新版或关闭自动检查时不追加提示；JSON 错误继续只输出 JSON，成功输出、登录参数、取消退出码 130 与原错误链保持。
- 同步 README 和命令参考说明；验证与安装证据保存到 `/tmp/contract-cli-text-update-notice-20261020/`，仍保留原回滚检查点。
- 验证：默认/dev 全量、Go/Node、race、vet 均通过，CLI 覆盖率 80.3%。候选版真实 dev `auth login --profile trace-dev-app --as user --no-open-browser` 到授权等待后发送 SIGINT，退出 130，stderr 先取消摘要、后单条模拟 99.0.0 升级提示；配置文件保持不变，显式 JSON 仍可解析且保留 notice。99.0.0 仅为用户授权设置的测试缓存，不代表真实发布版本。

## 2026-10-20 回滚检查点与 master 同步

- 用户授权先提交既有 dev、命令级 Trace、日志开关及验证文档；检查点提交为 `e2ab4ebb504ebd023bf7694fb5cc3478e06379df`，本地标记为 `codex/pre-interrupt-20261020`，不含本次统一中断处理。
- git fetch、ls-remote 与 merge 均已核对：origin/master 仍为 `e7bb2bc7d084dec4d9f5dcfb2f6785e66f68e1fd`，已被检查点包含，merge 返回 Already up to date，没有遗漏新的 master 提交，也没有制造空合并。
- 以后可要求“回滚到中断改造前的 codex/pre-interrupt-20261020”；本需求与该检查点分开保存，不推送或发布。

## 2026-10-20 实现与 TDD

- 用户明确选择只有显式 `--output json` 时错误输出 JSON；其他模式继续文本，成功 stdout 业务结构保持。JSON 错误按错误链给出类型、Trace、HTTP 状态和取消标志，附本次已有更新通知，错误渲染不发新网络请求。
- 主入口统一 SIGINT：首次取消并清理后退出 130，第二次恢复系统默认强退；Inspect 删除内部硬退出，沿 parent context 回收私有进程组，内部 inspection-helper 仍在 CLI 初始化前分发。
- 沿既有 resolveRawBody/resolveJSONObjectBody 调用链透传 context；Unix 由单个读取方拥有非阻塞 fd，Darwin 使用 Select、其他 Unix 使用 Poll，每次等待最多 100ms 并检查取消；退出前恢复原 O_NONBLOCK，仅关闭自己打开的 fd。其他平台保留 AfterFunc 关闭自有文件的取消方式，不关闭全局 stdin。
- import apply 保留原有状态保存、输出和锁范围；发送中取消保留 UncertainWriteError、uncertain 行，回读取消保留 unverified 行，返回取消链并停止后续写入。没有新增锁、重试或补偿机制。
- TDD 先确认导入取消被吞掉、真实进程缺取消 JSON 的红例；标准 FlagSet 文本诊断污染 JSON 也先失败，再统一交由 Execute 输出。覆盖值内的开关、无效 debug 参数、重复执行 notice 隔离、超时/HTTP/不确定写入分类，以及正常/取消 FIFO。
- macOS 半段 stdin 取消、原 stdin 标志恢复、空 writer 和空闲 writer 关闭后的正常 EOF 都补了红例后修复；Select 避开 [Go 已记录的 Darwin FIFO kqueue EOF 问题](https://github.com/golang/go/issues/24164)。Darwin 自有 fd 达到 FD_SETSIZE（1024）时明确报错，避免 fd_set 越界；Linux/Windows 仅完成交叉编译，未做实机验证。
- 验证：默认/dev 全量、Go/Node、vet、受影响模块 race 均通过；Darwin 真实进程 stdin/FIFO/探测回收及第二 SIGINT 通过，Linux/Windows 交叉编译通过。平台仅编译的结果不等同于实机运行。

## 2026-10-20 验证证据与交付

- 本地证据目录：`/tmp/contract-cli-interrupt-error-20261020/`，保留红绿日志、真实进程、交叉编译及实测记录。最终代码的 make test、dev 全量、两种构建 vet、受影响模块 race 全通过，CLI 单包覆盖率 80.3%；独立复审通过。
- 功能改动仅在 contract-cli，后端代码、配置和业务数据不改；测试中的写入取消仅使用本地 HTTP Transport 模拟。
- 候选版本 `1.9.2-beta.1-local.interrupt.dev`：真实进程读取半段 stdin 后 SIGINT，文本/JSON/debug JSON 三种模式均退出 130（本轮约 4–131ms）；显式 JSON 保留 Trace 和已有更新通知。更新通知使用合成缓存验证，不代表真实存在 99.0.0 新版本，也没有取消后联网检查。
- 2026-10-20 16:36 dev import plan：默认、debug、debug JSON 都成功且本地计划已取消，未执行 apply；默认 stderr 0 字节。两种 debug 的列头/规则行接口各自共享同一 Trace，分别为 `83ad395b4077f617c5d715f3471eae9a`、`4628f44348ec6ef2df2e3defb186ecb9`，OpenObserve 查到每条 Trace 的两个 open-platform server span 均 HTTP 200（仅核对 Trace，不声称核对 MDC 日志）。
- 交付时从本次功能提交重建并安装同一版本；保留原已安装二进制备份及 `codex/pre-interrupt-20261020`，中断功能与旧检查点分开提交，方便按检查点回滚。
