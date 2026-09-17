# dev 环境支持

- 当前分支将 test 预设替换为 dev，保留 prod；参考 `lmc/202609-task-dev` 的 dev 端点、OAuth PKCE、scope 和回调约定。
- 保留当前完整重定向链检查，阻断不支持的 test 域名；prod/dev profile、凭据与端点不可混用。不内置或复用生产 Device client。
- 保留 review、历史 rule table、其他业务命令以及调用来源请求头的实现。
- 使用独立 `contract-dev` profile，不删除或覆盖已有 test 凭据。`config add` 沿用既有行为，将新 profile 设为当前项。
- 测试覆盖 dev 配置/PKCE/历史命令/审查身份与来源头、test 拒绝、跨环境重定向和凭据恢复。
- 本地使用：`contract-cli config add --env dev --name contract-dev`，然后 `contract-cli auth login --profile contract-dev --as user`。
- 回滚本地安装可恢复安装前二进制；服务端、数据库和部署配置不在本次修改范围。
