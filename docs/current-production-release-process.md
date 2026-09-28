# 当前生产发版流程（2026-05）

本文整理 `contract-cli` 当前已经落地的正式发版流程、执行时需要准备的信息，以及这套流程当前已经达成的结果，供后续讨论“正式版后续怎么发”时作为基线。

## 1. 当前生产发版入口

当前正式发版入口是仓库脚本：

```bash
scripts/release.sh
```

配套能力：

- 本地检查：`make release-check`
- Release 附件构建：`make release-assets`
- GitHub Release 附件目录：`dist/release-assets/`

脚本目标是把正式版本一次性完成到以下两个发布面：

- GitHub 正式 Release
- npm `latest` dist-tag

## 2. 当前流程的标准步骤

### 2.1 版本号要求

正式版本必须是稳定语义版本：

```text
x.y.z
```

例如：

```text
0.1.3
1.0.0
```

`scripts/release.sh` 会拒绝带 `-beta.n` 的预发布版本号。

### 2.2 预演发布计划

先跑 dry-run：

```bash
scripts/release.sh --version 1.0.0 --dry-run
```

这一步会打印出计划执行的命令和预期附件名，但不会改文件，也不会触发远端发布。

### 2.3 本地准备

再跑本地准备模式：

```bash
scripts/release.sh --version 1.0.0
```

这一步会做三件事：

1. 更新 `package.json` 版本号
2. 执行 `make release-check`
3. 执行 `make release-assets`

执行完成后，脚本会停在“本地已准备完成”的状态，不会推送 GitHub，也不会发布 npm。

### 2.4 正式发布

确认本地准备无误后，再执行：

```bash
scripts/release.sh --version 1.0.0 --publish --yes
```

脚本按顺序完成：

1. 提交版本变更：`release: prepare <version>`
2. 打 tag：`v<version>`
3. 推送分支
4. 推送 tag
5. 创建或更新 GitHub 正式 Release
6. 上传 `dist/release-assets/` 下的附件
7. 执行 `npm publish --tag latest`
8. 校验 npm `latest` 指向目标版本

## 3. 当前流程需要准备的信息

### 3.1 必填信息

正式发版时，至少要准备这些信息：

- 目标版本号：例如 `0.1.3`
- 目标 Git remote：当前脚本默认读 `REMOTE`，未显式传值时默认 `origin`
- 目标分支：当前脚本默认取当前分支名，取不到时退回 `main`

### 3.2 发布权限

要满足至少一种 GitHub 授权方式：

- 本机已执行 `gh auth login`
- 或设置 `GITHUB_TOKEN`

要满足至少一种 npm 授权方式：

- 本机已执行 `npm login`
- 或设置 `NPM_TOKEN`

### 3.3 本地环境依赖

脚本会直接校验以下命令是否存在：

- `git`
- `node`
- `npm`
- `go`
- `zip`
- `make`
- 发布模式下额外需要 `gh`

### 3.4 工作区状态要求

脚本要求工作区可发布：

- 正常情况下必须是干净工作区
- 如果 `package.json` 已经先被改到目标版本，则允许只带少量白名单改动继续执行

当前白名单主要包括：

- `package.json`
- `package-lock.json`
- `*/.DS_Store`
- `scripts/release.sh`

其他无关改动会直接中断发布。

## 4. 当前流程的关键约定

### 4.1 Release 附件约定

正式发版会生成并校验以下附件：

- `checksums.txt`
- `contract-cli-<version>-darwin-amd64.tar.gz`
- `contract-cli-<version>-darwin-arm64.tar.gz`
- `contract-cli-<version>-linux-amd64.tar.gz`
- `contract-cli-<version>-linux-arm64.tar.gz`
- `contract-cli-<version>-windows-amd64.zip`
- `contract-cli-<version>-windows-arm64.zip`

这些附件是 npm 包安装脚本下载预编译二进制的来源之一，因此不是“可选产物”。

### 4.2 GitHub Release 语义

正式发版脚本创建的是 GitHub 正式 Release，而不是 pre-release。

脚本 dry-run 和测试里都已经约束：

- 正式版会带 `--latest`
- 正式版不会带 `--prerelease`

### 4.3 npm 发布语义

正式版固定发布到：

```text
latest
```

即：

```bash
npm publish --tag latest
```

发布完成后，脚本会检查：

```bash
npm view @qfeius/contract-cli@latest version --registry https://registry.npmjs.org
```

如果远端 `latest` 没有指向目标版本，脚本会直接失败。

## 5. 当前流程已经达成的结果

### 5.1 已有的一键化能力

目前正式发版已经具备这些自动化能力：

- 统一版本号写入
- 统一本地发版检查
- 统一多平台 release assets 构建
- 统一 GitHub Release 创建/覆盖上传
- 统一 npm `latest` 发布
- 发布后自动校验远端版本结果

### 5.2 已有的回归保护

仓库里已经有正式发版脚本的 dry-run 测试：

- [tests/release/release-script.sh](/Users/lyy/contract-cli/tests/release/release-script.sh)

当前测试明确保护了这些行为：

- dry-run 不修改 `package.json`
- 正式版附件名正确
- 正式版 GitHub Release 使用 `gh release create v<version>`
- 正式版带 `--latest`
- 正式版 `npm publish --tag latest`
- 正式版不会错误打成 `--prerelease`

### 5.3 与安装链路的打通结果

当前正式发版流程不只是“把包发出去”，它已经和安装链路打通：

- `make release-assets` 生成的附件可以被 npm 安装脚本消费
- npm 包元信息、GitHub Release 产物、二进制命名规则已经对齐
- 正式包安装文案和 `update check --channel latest` 语义已收敛到生产渠道

## 6. 当前已知边界和讨论点

虽然正式发版流程已经可用，但目前仍有这些边界，后续讨论时建议重点看：

- 远端默认值仍是 `origin`，如果 `origin` 不是 GitHub，需要手动传 `REMOTE=github`
- 分支默认值是“当前分支”，这意味着当前流程允许从非 `main` 分支发布，需要明确这是不是预期策略
- npm 远端校验是发布后即时读取，偶发传播延迟时可能出现“发布成功但首查未命中”
- 当前没有单独沉淀正式发版 skill；beta 已有本地 skill，正式版还没有同等封装
- 当前流程偏人工触发，尚未收敛成 CI/CD 固定入口

## 7. 建议作为后续讨论输入的问题

后续讨论“怎么发”时，可以直接围绕下面几个问题收敛：

1. 正式版是否必须只允许从 `main` 发布？
2. 正式版 remote 是否应该固定成 `github`，避免误发到其他远端？
3. 正式版是否也要沉淀成一个本地 skill，和 beta 一样只输入版本号与 key 即可执行？
4. 未来正式版是否要迁移到 GitHub Actions / CI，而不是本机手工触发？
5. 正式版发布后，是否需要增加安装回归或 smoke 验收的外部检查步骤？
