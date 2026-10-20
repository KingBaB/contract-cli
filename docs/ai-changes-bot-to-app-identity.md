# 统一将 bot 身份改名为 app

## 变更摘要

- 将推荐身份值从 `bot` 改为 `app`，新文档和 help 统一展示 `--as app`。
- 为避免影响老用户，CLI 继续兼容旧参数 `--as bot`，运行时按 `app` 身份处理。
- 配置文件写入新字段 `app_token_endpoint`、`identities.app`、`default_identity=app`。
- 兼容读取旧配置里的 `bot_token_endpoint`、`identities.bot`、`default_identity=bot`，保存后归一为 app 字段。
- 新增 `CONTRACT_CLI_APP_ID` / `CONTRACT_CLI_APP_SECRET`，旧 `CONTRACT_CLI_BOT_APP_ID` / `CONTRACT_CLI_BOT_APP_SECRET` 继续作为 fallback。

## 关键逻辑

- `ParseIdentityKind("bot")` 作为旧别名返回 `IdentityApp`，但新文档不再主推旧称。
- app 登录仍使用 `app_id/app_secret` 换取 `tenant_access_token/internal`，开放平台接口路径不变。
- app-only 命令沿用原有身份限制语义，只把错误提示更新为 `only supports --as app`。
- 旧 `profile.bot.app_secret` 仍可通过旧配置中的 `secret_ref` 读取；下一次 `auth login --as app` 会写入新的 `profile.app.app_secret`。

## 验证

- 新增和更新测试覆盖身份解析、`--as bot` 旧别名、旧配置迁移、app 登录、旧环境变量 fallback、app-only 拦截和文档口径。
- 运行 `go test ./...` 验证整体回归。

## 追加：MDM skill agent 身份口径对齐

### 变更摘要

- 修复 MDM 内置 skill 的 `agents/openai.yaml` 仍使用 `user-only` / `user-authorized` 旧描述的问题。
- 将 `contract-cli-mdm-vendor`、`contract-cli-mdm-legal` 的 agent 元信息改为 user/app 双身份查询口径。
- 将 `contract-cli-mdm-fields` 的 agent 元信息改为 user/app 字段查询口径，并明确 app 当前只覆盖 `vendor` / `legal_entity`，`vendor_risk` 仍走 user/MCP。
- 同步清理 MDM 参数参考中关于 `mdm fields` 仍保持 user-only 的过期描述。

### 验证

- 新增静态测试覆盖 MDM agent metadata 和 reference 文档，不允许残留过期身份口径。
