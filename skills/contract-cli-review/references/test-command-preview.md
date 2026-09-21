# 临时测试命令预览

仅用于入口已确认的 test 环境，禁止在生产启用或随生产 Skill 发布。环境依据是已解析的实际配置及请求 origin，不是 profile 名称；核实时只读取必要的非敏感字段，不输出完整配置、凭据文件或环境变量。

## 每次命令之前

先在用户可见的进度消息中输出下列内容，再调用执行工具；不能仅把命令写在工具参数里，也不能等执行完成后补打印。每次新命令、分页和显式轮询都重复此步骤，不以“同上”省略本次参数。

```text
[仅测试：执行前预览]
命令：<即将执行的完整命令及参数，敏感值替换为 [REDACTED]>
接口：<HTTP 方法 + 实际开放平台 origin + 路径；本地命令写“无远端接口”>
参数：
  CLI：<本次 flags、位置参数及非敏感实际值>
  Path / Query：<实际字段和值；没有则写“无”>
  Body / Form：<实际 JSON / multipart 字段和值；没有则写“无”>
  自动注入：<已确认的默认值、枚举转换和固定字段；未知值说明来源，不猜测>
```

- 包括首次 `review --help`、子命令帮助、合同搜索/详情/分类查询、协商详情、上传、主体提取、清单查询、发起和结果查询；流程内的本地脚本、文件处理命令也适用，无 HTTP 调用时不能虚构接口。
- `--input-file` / `--input` / `--data` 不能只展示文件路径：还须展示本次将提交的脱敏结构化内容。不要为了预览新增远端请求或执行写命令；需要读取输入的本地命令本身也先预览。
- `--dry-run` / `--print-input` 是本地预检：明确“本次不发送 HTTP”，可另列正式执行时的目标接口；归一化结果未取得时只展示已知输入，之后正式调用前再展示已核实的请求参数。不将 CLI 原始预检输出未经脱敏直接转发。
- 明确区分 CLI 输入与 HTTP 实际参数，例如 `reviewStrength` 中文输入转为 0/1/2、`fileId` 序列化、固定 `usageReportContext.reportBusinessCode=clm_100_openApi_cli`；以当前版本真实契约为准，不伪装成已抓取的网络请求。
- `review run` 或 `task result` 内部可能发多个请求：命令前列出涉及的接口、执行条件、轮询参数与已知输入。依赖上传响应的 ID/hash 等标为“运行时从上传响应取得”，不填假值，也不声称已逐个拦截内部 HTTP。此规则要求每次命令调用前预览，不修改 CLI 内部编排。
- 预览后使用原始、已校验参数安全执行，不能将含 `[REDACTED]` 的展示副本当作实际请求。预览不覆盖原 Skill 的创建确认、重试边界、user 身份和服务端权限要求；不为展示参数切换到裸 API。

## 脱敏与输出边界

- 命令和参数两处都必须脱敏：Authorization、Cookie、访问/刷新令牌、密码、app/client secret、会话、CSRF、OAuth code、签名及预签名下载 URL 的敏感 query 值，一律 `[REDACTED]`。URL 的用户名/密码同样不得展示；无法可靠区分敏感 query 时将整个 query 隐去。动态鉴权值标注“由 CLI 注入，已隐藏”，不读取它们来补齐预览。
- 文件参数仅展示路径、文件名、已知大小/类型及业务关联字段，不展示文件字节、Base64、合同正文或全文解析内容。JSON/脚本中若嵌入正文，以 `[合同内容已省略]` 代替该值，保留必要字段结构。
- 必要合同/文件/任务 ID、版本及业务配置可在测试预览中展示；正式选项列表仍不展示内部 ID。未经验证的字段写“尚未取得/待确认来源”，不能猜测真实参数或接口。
- 最终审查结果仍严格遵循主 Skill 模板，不追加命令、接口、参数或诊断记录。进入生产或失去可信 test 上下文后立即停止详细预览；原业务流程仍受其自身环境准入约束。

## 常用接口映射（个人 user 身份）

以下是当前 CLI 发出的开放平台入口，不是浏览器或下游服务的内部 URI。版本发生变化时以当前命令实现/已验证契约为准；没有证据的映射不得编造。

| 命令 | 方法与路径 |
| --- | --- |
| `contract search` | `POST /open-apis/contract/v1/mcp/contracts/search` |
| `contract get <id>` | `GET /open-apis/contract/v1/mcp/contracts/{id}` |
| `contract category list` | `GET /open-apis/contract/v1/mcp/contract_categorys` |
| `review cooperation get` | `GET /open-apis/contract/v1/mcp/cooperations/{id}` |
| `review file upload` | `POST /open-apis/contract/v3/review/file/upload` |
| `review file upload-url` | `POST /open-apis/contract/v3/review/file/uploadByUrl` |
| `review subject extract` | `POST /open-apis/contract/v3/review/subjects` |
| `review task start` | `POST /open-apis/contract/v3/review/startReview` |
| `review task status` | `GET /open-apis/contract/v3/review/status` |
| `review task info` | `GET /open-apis/contract/v3/review/info` |
| `review checklist list` | `GET /open-apis/contract/v1/review-rules/review-checklists` |
| `review task result` | 条件性循环 `GET .../review/status`，成功后 `GET .../review/info` |
| `review run` | 上传、可选主体提取、发起；按 wait 参数决定是否查询状态/详情，分别对应上表接口 |

未列出的命令先核实真实映射，再给预览；若无法核实，说明缺口，不把猜测当作接口信息。本参考不扩大本 Skill 到清单/规则维护或授权管理。
