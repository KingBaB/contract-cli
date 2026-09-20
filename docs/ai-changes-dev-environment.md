# 环境支持变更记录

## 2026-09-20：测试开发分支仅启用 test

- 适用分支：`test-feature-contract-cli-add-review`。新增默认 test 预设；prod 预设保留但由编译期常量暂时禁用，dev 不再提供可用预设。
- 修改统一 profile 校验和 HTTP 守卫：已有 prod/dev profile、Device 快照、直连及重定向请求都不得向禁用环境发送凭据；保持 test 的 OAuth PKCE 和独立凭据，不复用生产 Device client。
- 保留历史命令、身份权限、请求头及审查 Skill 的 `pc_url` 导航规则；不硬编码或替换结果域名。
- 同步帮助、README、命令参考及环境测试夹具，保留历史业务回归并新增 test-only 公开入口测试。
- 本次仅源码变更，不删除本机配置、不安装、不发布；未来恢复 prod 需调整开关并重新验证环境及凭据隔离。

## 历史：dev 开发分支（不再代表当前测试分支）

- 当前分支将 test 预设替换为 dev，保留 prod；参考 `lmc/202609-task-dev` 的 dev 端点、OAuth PKCE、scope 和回调约定。
- 保留当前完整重定向链检查，阻断不支持的 test 域名；prod/dev profile、凭据与端点不可混用。不内置或复用生产 Device client。
- 保留 review、历史 rule table、其他业务命令以及调用来源请求头的实现。
- 使用独立 `contract-dev` profile，不删除或覆盖已有 test 凭据。`config add` 沿用既有行为，将新 profile 设为当前项。
- 测试覆盖 dev 配置/PKCE/历史命令/审查身份与来源头、test 拒绝、跨环境重定向和凭据恢复。
- 本地使用：`contract-cli config add --env dev --name contract-dev`，然后 `contract-cli auth login --profile contract-dev --as user`。
- 回滚本地安装可恢复安装前二进制；服务端、数据库和部署配置不在本次修改范围。
