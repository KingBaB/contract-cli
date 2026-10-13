# 智书合同 CLI 使用指南（Skill 能力修订副本）

> 能力基线：GitHub Release **v1.9.5**；核对日期：2026-10-13。本文面向使用 AI Agent 的客户、实施和客户成功人员。

你可以直接告诉 Agent：“查一下上个月申请、目前已归档的采购合同”“下载这份合同的归档文件”“审查这份采购协议”“看看我有哪些合同审批待办”。Agent 根据随 CLI 安装的 skills，选择命令、查询真实候选、补齐必要信息，再在你的权限范围内执行。

`contract-cli` 提供连接智书合同服务的命令；配套 skills 告诉 Agent 如何把业务需求转成可执行步骤，包括查找人员与部门、识别字段、选择审查立场、处理分页，以及核验执行结果。你负责说明目标和确认关键业务选择，Agent 负责具体操作。

**个人授权（user）**以当前用户身份访问；**应用授权（app）**使用企业应用获授的访问能力。两种授权相互独立，应用授权不包含全部个人能力。能否完成任务，要同时满足“该操作支持这种身份、授权有效、具备相应业务权限”三个条件。

## 快速开始

### 先确认任务需要哪种授权

Skill 负责选择和编排业务操作，授权由“授权与身份管理”能力统一处理。同一种身份在当前租户、配置和有效凭证范围内可以复用，不需要每切换一个 Skill 就重新授权；个人授权和应用授权不能互相替代。

| 授权方式 | 由谁完成 | 可以使用哪些能力 |
| --- | --- | --- |
| **个人授权（user）** | 使用者按官方链接完成手机号验证、企业确认和授权 | 能力地图标为“仅个人授权”的能力，以及“两种授权”中支持个人身份的操作；实际范围受本人业务权限限制 |
| **应用授权（app）** | 管理员或集成人员配置企业应用凭据，并确认应用已获得所需权限 | 能力地图标为“仅应用授权”的能力，以及“两种授权”中支持应用身份的操作；实际范围受应用授权与接口业务校验限制 |

日常个人任务按下面的个人授权流程开始；明确使用应用身份的任务，可直接配置和检查应用授权，无须先完成个人授权。一个任务需要两种身份时，Agent 应先说明各步骤需要的授权，确认使用范围后分别执行。

### 1. 安装 CLI 和同版本 skills

把下面这段话发送给 Agent：

```text
请帮我安装或升级 @qfeius/contract-cli，并同步更新当前客户端使用的全部内置 skills。
先确认命令运行在本机还是沙箱，以及当前客户端实际读取的 skills 目录。
使用 npm 官方源安装，检查实际执行路径和版本，再用同一个 CLI 执行：
contract-cli skills install --target "<已确认的真实技能目录>" --force
按 skills list 核对名称、数量、SKILL.md 和附属文件，不能只补装缺少的 skill。
如果目录中有本地定制，请先识别并保留定制内容，再按更新方案处理。
最后分别报告 CLI 版本、安装目录、文件同步结果和客户端加载结果。
```

手动安装时执行：

```bash
npm install -g @qfeius/contract-cli@latest --registry https://registry.npmjs.org
contract-cli --version
contract-cli skills list
contract-cli skills install --target "<当前客户端实际读取的 skills 目录>" --force
```

将占位目录替换为已确认的真实路径。`--force` 会覆盖目标目录中的同名内置 skill，包含其附属文件；有本地定制时先保留。不要用 GitHub 默认分支上的 skills 替代当前 CLI 内置版本。

本文核对的 v1.9.5 内置 **16 个 skills**。安装 `latest` 后若版本与本文不同，应按该版本的实际能力核对。后续版本以实际 `skills list` 为准；该命令显示的是包内清单，不能单独证明客户端已经安装或加载。安装输出中的 skill 版本与 CLI 版本对齐，源码文件头中的独立版本号不作为安装版本依据。

更新后完全退出并重新启动客户端，再新建任务核对加载情况。仅文件已更新但尚未确认加载时，应报告“文件已更新，客户端加载待验证”。

### 2. 初始化并完成个人授权

v1.9.5 当前支持生产环境，按配套 skill 使用 `contract` profile：

```bash
contract-cli config add --env prod --name contract
contract-cli auth status --profile contract --as user
```

日常可直接告诉 Agent：

```text
请使用 contract-cli，以我的个人身份处理合同任务。
先复用已有有效授权；需要授权时请给我官方授权链接，完成后继续原任务。
```

在配套 skill 覆盖的豆包、WorkBuddy、Codex 场景中，Agent 使用手机号授权流程：需要授权时发起一次 `auth init`，展示授权链接和有效期；你在页面完成手机号验证、企业确认与授权后，回复“已授权”，Agent 再查询一次授权结果。

```bash
# 仅在需要新授权或处理可复用的 Device 授权时执行
contract-cli auth init --profile contract --output json
# 仅在等待授权且用户回复“已授权”后执行
contract-cli auth complete --profile contract --output json
```

已有有效授权，或 `auth init` 返回已授权时，直接继续原任务。等待授权时不反复查询；授权过期、被拒绝或结果不确定时，按实际状态处理。

本地桌面与云端沙箱的登录态复用范围不同，以 CLI 返回的凭证范围为准。旧版本的手动回调脚本不再作为本文推荐入口；其他宿主的授权适配需按实际环境核验。本地仍保留 `auth login --as user` 授权码登录方式。

### 3. 按任务需要配置应用授权

付款、固定汇率、主数据文件下载、事件出口 IP，以及部分合同管理操作仅支持应用授权。任务明确选择应用身份时，由管理员在对话外安全配置本机环境变量或凭据存储，并确认应用具有所需权限，再执行：

```bash
contract-cli auth login --profile contract --as app
contract-cli auth status --profile contract --as app
```

环境变量名为 `CONTRACT_CLI_APP_ID`、`CONTRACT_CLI_APP_SECRET`。不要把 App Secret、Token、Cookie 或密码发给 Agent。Agent 不能因为个人身份查询无结果或权限不足，就自行换成应用身份。

### 4. 从一个明确的查询开始

```text
帮我查合同编号包含“HT2026-001”的合同，以我的身份查询。
列出名称、编号和状态；如果有多个候选，先让我确认目标。
```

## 用户场景

以下按 v1.9.5 面向用户的 15 个 skills 分组，按各 Skill 的实际能力展开代表场景；合同管理覆盖创建、审批、文件、签署与协作等常用业务。中文能力名称与 Skill 的对应关系统一见下方能力地图；涉及其他能力时，Agent 再调用对应 Skill 协作。示例中的名称、编号和文件名请替换为实际对象。右侧展示关键命令；JSON 文件由 Agent 按对应 skill 的真实字段与已确认参数生成。多步骤场景仍需逐步查询、确认并核验，不能将整列命令直接连续执行。

### 1. 授权与身份管理

#### 场景 1.1：完成授权后继续处理合同

**授权支持：**个人

Agent 检查当前授权，必要时引导你在官方页面完成手机号验证和企业授权；收到“已授权”后查询一次结果，再交给合同搜索 skill 继续原任务。应用身份只使用已在本机安全配置好的凭据，不在对话中索取 App Secret。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请帮我完成智书合同的个人授权，然后继续查询合同“HT2026-001”。<br>如果已有有效授权就直接复用；需要重新授权时，给我官方授权链接。 | 先检查授权：<br>`contract-cli auth status --profile contract --as user`<br><br>仅在需要授权时：<br>`contract-cli auth init --profile contract --output json`<br><br>仅在用户完成授权后：<br>`contract-cli auth complete --profile contract --output json` |


### 2. 合同搜索

#### 场景 2.1：按关键词快速找到合同

**授权支持：**个人

Agent 在当前用户可见范围内按默认关键词范围查询，不擅自缩小为仅合同名称。你明确要求“合同文本”或“仅正文”时，按相应范围查询。结果说明分页情况，区分合同组与合同条目，不把同组成员全部当作独立命中。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 帮我搜索“采购”，范围与合同页面默认“全部”一致。<br>列出合同名称、编号和状态，并说明是否还有下一页。 | `contract-cli contract search --profile contract --as user --input-file keyword-search.json` |


#### 场景 2.2：按人员和自定义字段组合筛选

**授权支持：**个人

Agent 借助人员 skill 查候选，通过字段发现获取真实字段及选项，再按已支持的组合表达查询。“我的合同”页签不等于“我申请的合同”；我方主体精确 ID 筛选目前没有已验证的通用配方，复杂条件也不能承诺任意 AND/OR 都可表达。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 查张三申请的、自定义字段“项目区域”为“华东”的合同。<br>先查真实人员和字段选项；有重名或字段歧义时让我确认，查询时保留这两个条件。 | 查询人员候选：<br>`contract-cli employee list --name "张三" --profile contract --as user`<br><br>发现真实字段：<br>`contract-cli contract search-fields --keyword "项目区域" --profile contract --as user`<br><br>保留两个条件查询：<br>`contract-cli contract search --profile contract --as user --input-file combined-search.json` |


### 3. 合同创建、审批与全流程业务操作

#### 场景 3.1：查看合同并下载归档材料

**授权支持：**个人 / 应用

Agent 读取已知合同，按指定字段收集文件并下载。明确要求“全部文件”时，还会收集可访问的流程和评论附件，去重后逐项处理。该流程由 skill 编排，没有独立的全部下载命令；归属不明、字段不可见或读取失败时会说明材料缺口。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请查看合同 ID 为“&lt;contract-id&gt;”的详情和正文，整理交易双方、金额及状态。<br>把归档主文件和归档附件下载到我指定的目录，列出文件名、来源和下载结果。 | 读取合同：<br>`contract-cli contract get <contract-id> --profile contract --as user`<br><br>读取正文：<br>`contract-cli contract text <contract-id> --profile contract --as user --full-text`<br><br>按已确认文件逐个下载：<br>`contract-cli contract download-file <file-id> --contract <contract-id> --profile contract --as user --output-file ./归档文件.pdf` |


#### 场景 3.2：用已有文件创建合同

**授权支持：**个人 / 应用

把合同正文和业务信息交给 Agent，它会查询可用分类、核对签约主体、金额与期限，先列出缺项，再整理创建资料。正文和附件先上传取得平台文件 ID；分类简称使用与创建身份一致的查询结果。创建支持个人和应用身份，应用身份还需提供真实创建人。创建成功后返回合同 ID 并核对详情；提交审批另行处理。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请以我的个人身份，用我提供的“采购合同.docx”创建采购合同。<br>先核对分类、我方主体、交易方、金额和有效期，缺少的信息先问我；展示创建摘要，确认后创建，暂不提交审批。 | 查询可用分类：<br>`contract-cli contract category list --profile contract --as user --lang zh-CN`<br><br>资料确认后上传正文：<br>`contract-cli contract upload-file --profile contract --as user --file ./采购合同.docx --file-type text`<br><br>将真实分类、主体和返回的文件 ID 写入请求体后创建：<br>`contract-cli contract create --profile contract --as user --input-file contract-create.json`<br><br>用返回的合同 ID 核验：<br>`contract-cli contract get <contract-id> --profile contract --as user` |

#### 场景 3.3：选择标准模板并创建合同

**授权支持：**个人 / 应用

Agent 查询分类下的模板、读取模板字段，补齐变量后先生成模板实例，再使用返回的实例 ID 创建合同。模板实例与合同是两个对象；生成实例并不代表合同已经创建。支持个人和应用身份，分类、模板与创建保持同一身份。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请用系统里的标准采购模板创建一份合同。<br>先列出可用分类和模板让我选择，再按模板字段核对双方、金额和日期；展示填写结果，确认后创建合同并返回详情。 | 查询分类并取得真实分类编号：<br>`contract-cli contract category list --profile contract --as user --lang zh-CN`<br><br>列出模板并读取所选模板：<br>`contract-cli contract template list --profile contract --as user --category-number <category-number> --page-size 20`<br>`contract-cli contract template get <template-id> --profile contract --as user`<br><br>补齐变量后生成实例：<br>`contract-cli contract template instantiate --profile contract --as user --input-file template-instance.json`<br><br>将返回的实例 ID 写入创建资料：<br>`contract-cli contract create --profile contract --as user --input-file template-contract-create.json` |

#### 场景 3.4：基于原合同创建补充协议或终止申请

**授权支持：**个人 / 应用

Agent 先读取原合同，核对关联关系，再按合同变更或终止的创建配方整理资料。变更需要变更说明；终止需要终止原因、日期及日期类型，分类、签约主体和正文也按要求补齐。这里创建的是关联原合同的业务单据，创建成功不代表原合同已经完成变更或终止。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请基于合同“&lt;contract-id&gt;”和我提供的补充协议文件，创建一份合同变更申请。<br>先展示原合同、变更原因和新协议资料，确认后创建，暂不提交。 | 读取原合同：<br>`contract-cli contract get <contract-id> --profile contract --as user`<br><br>按场景 3.2 查询分类、上传协议，补齐变更资料后：<br>`contract-cli contract create --profile contract --as user --input-file contract-change-create.json` |
| 请基于这份合同准备终止申请，使用我提供的终止原因和生效日期。<br>缺项先问我，确认资料后创建申请。 | 读取原合同并按相同身份核对分类、正文和终止资料后：<br>`contract-cli contract get <contract-id> --profile contract --as user`<br>`contract-cli contract create --profile contract --as user --input-file contract-termination-create.json` |

#### 场景 3.5：提交合同审批，或处理退回后的流程

**授权支持：**应用

合同提交和重新提交目前仅支持应用身份。首次提交限开放平台创建的草稿合同；重新提交要求合同已撤回或被拒绝，且当前处于提交节点。应用流程接口还支持按权限处理同意、拒绝、撤回、重新提交及作废；其中作废仅限提交节点。Agent 会核对合同、流程、任务和操作人，按明确选择执行一项操作，再查询结果。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请用已授权的应用身份检查这份开放平台创建的草稿合同。<br>资料和状态满足要求后，按我的确认提交审批，并返回流程信息。 | 核对合同：<br>`contract-cli contract get <contract-id> --profile contract --as app`<br><br>确认提交后：<br>`contract-cli contract submit <contract-id> --profile contract --as app`<br><br>用返回的流程 ID 查询：<br>`contract-cli contract approval get <process-instance-id> --profile contract --as app` |
| 这份合同被退回了，请检查是否已回到提交节点。<br>资料已按意见处理完成，满足条件后重新提交。 | 核对合同与流程后执行重新提交：<br>`contract-cli contract resubmit <contract-id> --profile contract --as app` |
| 请用应用身份检查这个已有流程能否撤回，并展示将要操作的合同和任务。<br>我确认后再执行撤回；如果节点不允许，就说明原因。 | 读取已有流程：<br>`contract-cli contract approval get <process-instance-id> --profile contract --as app`<br><br>核对操作人、任务及动作后执行：<br>`contract-cli contract approval start <process-instance-id> --profile contract --as app --input-file approval-action.json`<br><br>该命令操作已有流程；不是创建新审批流程的入口。 |

#### 场景 3.6：查看我的审批待办、已办和知会，追踪流程与评论

**授权支持：**个人

Agent 按个人身份查询任务，结合流程详情和评论说明当前节点、处理进度及待补材料。流程实例 ID、任务实例 ID 与合同 ID 分别使用。评论可以新增、回复、@ 人员或关联附件；提交前按真实对象整理内容。流程数据不完整时会明确说明，不能将缺失字段判断为流程结束。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 列出我的合同审批待办；找到合同“CN-001”对应的任务，告诉我当前在哪个节点、最近有哪些评论。 | 查询待办：<br>`contract-cli contract approval task list --profile contract --as user --task-type todo --query "CN-001" --page-size 20`<br><br>读取所选流程与评论：<br>`contract-cli contract approval get <process-instance-id> --profile contract --as user`<br>`contract-cli contract approval comment list <process-instance-id> --profile contract --as user` |
| 我想回顾已经处理过的合同，以及抄送给我的合同。请分别列出结果，并说明是否还有下一页。 | 分别查询已办与知会：<br>`contract-cli contract approval task list --profile contract --as user --task-type done --page-size 20`<br>`contract-cli contract approval task list --profile contract --as user --task-type notice --page-size 20` |
| 请在我选定的流程中添加评论：“请补充付款条件说明”，并附上我提供的说明文件。 | 上传同用户的评论附件：<br>`contract-cli contract upload-file --profile contract --as user --file ./付款条件说明.pdf --file-type reviewAttachment`<br><br>将返回文件 ID 与已确认评论写入请求体后提交：<br>`contract-cli contract approval comment create <process-instance-id> --profile contract --as user --input-file approval-comment.json`<br><br>回查评论确认结果：<br>`contract-cli contract approval comment list <process-instance-id> --profile contract --as user` |

#### 场景 3.7：按我的意见通过或拒绝审批待办

**授权支持：**个人

Agent 先读取个人待办和流程，核对可办理任务，再按你明确选择的意见通过或拒绝；拒绝必须填写原因。需要审批附件时，以同一用户上传专用审批附件，并在授权有效期内使用。盖章、归档节点需回合同系统处理。办理结果不确定时先查状态，不能直接重复提交。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请打开我选择的待办，核对合同和当前节点。<br>我确认通过后，以“同意”办理；如需附上我提供的审批意见文件，先按审批附件上传。 | 核对任务与流程后，按需上传新附件：<br>`contract-cli contract upload-file --profile contract --as user --file ./审批意见.pdf --file-type approveAttachment`<br><br>明确通过后执行；仅在有有效新附件时加文件参数：<br>`contract-cli contract approval task approve <task-instance-id> --profile contract --as user --comment "同意" --file-id <file-id>`<br><br>无附件时省略 `--file-id`，完成后回查流程。 |
| 请拒绝我选定的这条审批待办，意见是“付款条件尚未确认，请补充后重新提交”。 | 核对选定任务可办理后执行：<br>`contract-cli contract approval task reject <task-instance-id> --profile contract --as user --comment "付款条件尚未确认，请补充后重新提交"`<br><br>回查处理结果：<br>`contract-cli contract approval get <process-instance-id> --profile contract --as user` |

#### 场景 3.8：补齐归档材料，或生成合同信息打印文件

**授权支持：**应用

应用身份可按合同状态补充规定的扫描件、归档主文件和归档附件。Agent 先读取合同，上传对应用途文件，再更新允许的文件字段；这不等于可任意修改合同名称、金额等基础字段，也不代替办理归档节点。生成打印文件会先返回平台文件 ID，随后下载到指定目录。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 这份合同正在归档，请把我提供的签署扫描件补为归档主文件。<br>先核对合同状态及文件用途，确认后更新，并告诉我结果。 | 读取合同并核对归档状态：<br>`contract-cli contract get <contract-id> --profile contract --as app`<br><br>上传扫描件：<br>`contract-cli contract upload-file --profile contract --as app --file ./签署扫描件.pdf --file-type scan`<br><br>将返回文件 ID 写入对应归档字段后更新：<br>`contract-cli contract patch <contract-id> --profile contract --as app --input-file contract-archive-patch.json` |
| 请生成这份合同的合同信息打印文件，并下载到我指定的目录。 | 生成打印文件：<br>`contract-cli contract print-file --profile contract --as app --input-file contract-print.json`<br><br>使用返回的文件 ID 下载：<br>`contract-cli contract download-file <file-id> --profile contract --as app --output-file ./合同信息.pdf` |

#### 场景 3.9：获取签署与实名认证链接，或调整签署方式

**授权支持：**应用

应用身份可获取合同签署链接，以及个人、企业实名认证链接，由相关人员在页面完成认证和签署。需要改为纸质签署时，先核对真实业务对象、业务类型和状态，再执行切换。获取链接不代表已经完成认证或签署。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请获取这份合同的签署链接，并根据我提供的签署主体信息，准备所需的个人或企业实名认证链接。 | 获取签署链接：<br>`contract-cli contract sign-url get <contract-id> --profile contract --as app`<br><br>按实际主体选择一种认证入口：<br>`contract-cli contract esign personal-auth-url --profile contract --as app --input-file personal-auth.json`<br>或：<br>`contract-cli contract esign org-auth-url --profile contract --as app --input-file org-auth.json` |
| 这份合同申请确定改为线下纸质签署。请先核对目标合同和当前签署方式，确认后执行切换。 | 核对合同申请业务类型及状态后：<br>`contract-cli contract sign switch-to-paper --profile contract --as app --business-id <contract-id> --business-type-code 0` |

#### 场景 3.10：查询分享记录，并按指定范围分享或授权

**授权支持：**应用

Agent 使用应用身份查询分享记录，按明确的合同、接收人员及授权期限执行分享或授予权限。姓名须先解析为真实人员，存在歧义时先确认；分享和授权是不同操作，按业务目标分别选择，不自动合并执行。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 查一下这份合同已经分享给谁；再按我确认的接收人名单分享，完成后核对记录。 | 查询现有分享：<br>`contract-cli contract share get <contract-id> --profile contract --as app`<br><br>接收人确认后分享：<br>`contract-cli contract share batch-create --profile contract --as app --input-file contract-share.json`<br><br>再次查询分享记录核验。 |
| 请给指定人员授予这份合同的访问权限，期限按我提供的开始和结束时间设置。执行前先展示对象和期限。 | 核对合同、人员及授权期限后：<br>`contract-cli contract authorization grant --profile contract --as app --input-file contract-authorization.json` |

#### 场景 3.11：查看协商进展，并取得协商链接和文件

**授权支持：**应用

应用身份可搜索协商合同、查询协商操作记录与入口链接，并取得协商文件信息后下载。Agent 按返回记录整理参与人、操作时间和进展；查询记录不代表能够代替参与人确认合同、定稿或接受修订。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请查找与“采购”相关的协商合同，打开我选定的一份，整理最近的协商记录并给我协商入口。 | 查询协商列表：<br>`contract-cli contract cooperation search --profile contract --as app --input-file cooperation-search.json`<br><br>查询所选合同的记录与入口：<br>`contract-cli contract cooperation record get <contract-id> --profile contract --as app`<br>`contract-cli contract cooperation link get <contract-id> --profile contract --as app` |
| 请查看这份合同可取得的协商文件，并把我选定的文件下载到指定目录。 | 取得文件信息：<br>`contract-cli contract cooperation file get <contract-id> --profile contract --as app`<br><br>使用返回的文件 ID 下载：<br>`contract-cli contract cooperation file download <file-id> --profile contract --as app --output-file ./协商文件.docx` |

#### 场景 3.12：核对合同字段配置，并维护下拉选项

**授权支持：**个人 / 应用

实施人员可通过应用身份读取指定分类和业务类型的流程字段，并维护已支持的下拉选项范围。个人身份另可查询合同枚举。字段选项配置与某份合同上的字段值不同；该能力不包含任意新增字段或任意修改合同数据。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请查出这个合同分类在合同申请流程中使用的字段，便于我核对录入要求。 | 使用真实分类 ID 查询合同申请的流程字段：<br>`contract-cli contract form attribute list --profile contract --as app --category-id <category-id> --business-type-code 0` |
| 请查询合同状态枚举；另外，把指定字段组下这个下拉字段的选项按我确认的完整清单更新。 | 个人身份查询枚举：<br>`contract-cli contract enum list --profile contract --as user --type contract_status`<br><br>应用身份按已确认字段和完整选项清单更新：<br>`contract-cli contract field update --profile contract --as app --input-file contract-field-options.json` |

#### 场景 3.13：同步合同用户分组

**授权支持：**个人 / 应用

Agent 可在已授权的个人或应用身份下触发合同用户分组同步，报告实际执行结果。该能力用于同步已有分组关系，不提供任意新建分组、修改通讯录或自行扩大合同权限的入口。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请以我的个人身份执行合同用户分组同步，并告诉我同步是否成功；如果返回失败，列出原因。 | `contract-cli contract sync-user-groups --profile contract --as user` |

#### 场景 3.14：清理确认不再使用的合同草稿

**授权支持：**应用

应用身份可删除同一应用下创建的草稿合同。Agent 先核对名称、ID、状态及应用归属，展示待删除对象；只有明确确认删除后才执行。此入口不适用于任意状态的合同。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 这份测试合同不再使用。请先检查它是不是本应用创建的草稿，展示名称和 ID，等我确认后再删除。 | 检查合同：<br>`contract-cli contract get <contract-id> --profile contract --as app`<br><br>满足限制且明确确认删除后：<br>`contract-cli contract delete <contract-id> --profile contract --as app` |

### 4. 人员查询

#### 场景 4.1：确认同名人员，避免选错申请人

**授权支持：**个人

Agent 使用个人身份查询真实人员候选，保留对应人员标识。确认后的候选可交给合同搜索 skill，用于申请人、需求人或自定义人员字段筛选。人员查询只读，不自动同步或修改通讯录。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 帮我查找姓名为“张三”的人员，列出所属部门和人员状态。<br>如果有多个同名人员，先让我选择，不要默认取第一位。 | `contract-cli employee list --name "张三" --profile contract --as user` |


### 5. 部门查询

#### 场景 5.1：定位实际业务部门

**授权支持：**个人

Agent 使用个人身份查询部门候选，保留真实部门标识，供后续人员查询或合同部门筛选使用。申请部门与需求人部门是不同业务角色，后续筛选需按你指定的角色处理；本 skill 不修改组织架构。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 帮我查找名称包含“华东销售”的部门，列出部门层级和状态。<br>如果有多个候选，先让我确认要使用哪个部门。 | `contract-cli department list --name "华东销售" --profile contract --as user` |


### 6. 合同智能审查

#### 场景 6.1：审查一份外部合同文件

**授权支持：**个人

Agent 告知文件将提交至智书合同审查服务，再按分类、主体、强度和清单推进。当前支持单份 DOC、DOCX、PDF，本地文件或原始文件流最大 50 MiB；主体来自同一文件的提取结果。外部预审无须创建合同或协商，结果在 Agent 中展示，不提供系统合同详情链接。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请审查我提供的《采购协议.docx》。<br>上传后先让我选择合同分类，再从文件提取的主体中确认保护立场。<br>审查强度用中立，清单从我有权限使用的清单中选择，最后展示风险和建议。 | 上传文件：<br>`contract-cli review file upload --profile contract --as user --file ./采购协议.docx --name 采购协议.docx`<br><br>提取同一文件的主体：<br>`contract-cli review subject extract --profile contract --as user --business-id <business-id> --file-id <file-id>`<br><br>分类、立场、强度和清单确定后：<br>`contract-cli review task start --profile contract --as user --input-file review-start.json` |


#### 场景 6.2：审查系统合同，或继续已有任务

**授权支持：**个人

Agent 必要时先调用搜索 skill 定位合同，再核实来源、权限及状态。已识别状态下，新发起范围为协商的待发起协商、协商中，以及任务中心的审批中、签订中、归档中、已归档；未知状态由服务端校验。多个正文来源需先选择；已有任务沿用原任务及真实业务关联，符合导航条件时返回“查看审查”链接。审查不自动修改正文、接受修订或通过审批。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请审查合同“HT2026-001”的当前正文，先核实合同和正文来源。<br>让我从该文件提取的主体中选择立场；如果本次审查已有任务，就继续查询原任务。 | 新审查先读取目标合同：<br>`contract-cli contract get <contract-id> --profile contract --as user`<br><br>已有任务只续查原任务：<br>`contract-cli review task result --task-id <task-id> --business-id <business-id> --profile contract --as user` |


### 7. 审查清单与规则维护

#### 场景 7.1：了解团队已有的审查标准

**授权支持：**个人

Agent 查询当前租户可见配置，按所需范围处理分页；信息不完整时说明缺口。仅在审查时选择已有清单，仍由智能审查 skill 处理，不需要进入配置维护流程。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请查询我有权限查看的采购合同审查清单。<br>展示适用分类、包含的规则及规则所属分组，区分内置清单和自定义清单。 | 按已确认分类查询清单：<br>`contract-cli review checklist list --profile contract --as user --contract-category <category-id> --page-index 1 --page-size 100`<br><br>查询规则分组：<br>`contract-cli review rule group list --profile contract --as user --page-index 1 --page-size 100`<br><br>读取相关分组规则：<br>`contract-cli review rule list --profile contract --as user --group-id <group-id> --page-index 1 --page-size 100` |


#### 场景 7.2：补充一条团队审查规则

**授权支持：**个人

Agent 按已支持的单项流程维护自定义配置。内置清单只读；删除须确认完整影响，规则或分组删除还需证明完整引用范围及服务端单次原子级联能力。当前不支持规则跨分组迁移，不承诺整套清单复制；审查规则与审批矩阵由不同 skill 负责。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 我想新增一条检查付款前置条件的审查规则。<br>请先整理规则内容、所属分组及需要引用它的自定义清单，逐项展示变更，确认后保存并回读。 | 确认内容并通过预检后创建：<br>`contract-cli review rule create --profile contract --as user --group-id <group-id> --input-file review-rule.json`<br><br>取得新规则 ID 并确认完整引用差异后：<br>`contract-cli review checklist update --profile contract --as user --id <checklist-id> --input-file checklist-update.json` |


### 8. 交易方管理

#### 场景 8.1：查询交易方资料并确认目标

**授权支持：**个人 / 应用

Agent 可用个人或应用身份查询列表与详情。按证件查询和全量分页为应用专属能力，不能因为个人查询无结果就自动改用应用身份。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 帮我查找“示例供应商”的交易方资料和状态。<br>有同名记录时列出候选供我选择，确认后再查看详情。 | 查询候选：<br>`contract-cli mdm vendor list --profile contract --as user --name "示例供应商"`<br><br>读取已确认交易方：<br>`contract-cli mdm vendor get <vendor-id> --profile contract --as user` |


#### 场景 8.2：局部更新交易方信息

**授权支持：**个人 / 应用

Agent 使用对应身份的局部更新能力。个人创建和局部更新不发起审批，不能写入编码、风险和系统字段；启用、停用使用独立的个人身份操作，附件仅引用已有主数据文件。旧版全量更新仍为应用身份能力。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请按我提供的新资料更新这家已确认交易方的联系人信息。<br>先展示修改前后的差异，确认后只更新这些字段，其他信息保持不变。 | 确认差异后局部更新：<br>`contract-cli mdm vendor patch <vendor-id> --profile contract --as user --input-file vendor-patch.json`<br><br>回读核验：<br>`contract-cli mdm vendor get <vendor-id> --profile contract --as user` |


### 9. 法人主体管理

#### 场景 9.1：确认本次签约的我方主体

**授权支持：**个人 / 应用

Agent 可用个人或应用身份查询列表及按 ID 查看详情，避免直接把名称当成唯一主体。该查询可为合同创建准备主体信息，但不代表合同搜索已经支持主体 ID 精确筛选。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请查询我可以使用的法人主体，帮我找到“示例科技有限公司”。<br>展示候选并让我确认，再读取所选主体的详情。 | 查询主体候选：<br>`contract-cli mdm legal list --profile contract --as user --name "示例科技有限公司" --page-size 20`<br><br>读取已确认主体：<br>`contract-cli mdm legal get <legal-entity-id> --profile contract --as user` |


#### 场景 9.2：通过企业应用维护法人主体

**授权支持：**应用

法人主体按编码查询、创建和更新使用应用身份。字段要求通过主数据字段 skill 核对，操作人按接口要求提供，不从个人查询权限推断具有应用写权限。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请使用已配置的应用身份，根据我提供的资料准备新法人主体的创建参数。<br>先核对字段要求和操作人，展示完整内容，确认后创建并查询结果。 | 核对字段：<br>`contract-cli mdm fields list --profile contract --as app --biz-line legal_entity`<br><br>确认参数后创建：<br>`contract-cli mdm legal create --profile contract --as app --user-id <operator-user-id> --input-file legal-create.json`<br><br>按创建结果回读：<br>`contract-cli mdm legal get <legal-entity-id> --profile contract --as app` |


### 10. 主数据字段配置

#### 场景 10.1：先看字段要求，再准备主数据

**授权支持：**个人 / 应用

Agent 查询配置并整理填写要求，不自行修改字段定义。个人身份支持交易方、法人主体和交易方风险字段配置；应用身份支持交易方和法人主体。合同搜索字段发现属于合同搜索 skill，不使用此入口代替。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请查看当前交易方的字段配置，整理字段名称、类型和必填要求。<br>对照我提供的资料列出缺少的信息，先不要创建交易方。 | `contract-cli mdm fields list --profile contract --as user --biz-line vendor` |


### 11. 固定汇率维护

#### 场景 11.1：查询指定生效日的固定汇率

**授权支持：**应用

Agent 按明确的源币种、目标币种和生效日期查询系统固定汇率，不将结果解释为实时市场行情。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请使用已配置的应用身份，查询 2026 年 10 月 1 日人民币到美元的固定汇率。<br>展示币种方向、生效日期和查询结果。 | `contract-cli mdm fixed-exchange-rate get --profile contract --as app --source-currency CNY --target-currency USD --effective-date 2026-10-01` |


#### 场景 11.2：更新已确认的汇率配置

**授权支持：**应用

Agent 在必要字段完整后执行配置更新，不自行推算或选择汇率值。该能力仅支持应用身份。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请按我提供的币种方向、生效日期和汇率值准备更新内容。<br>先展示完整配置，确认后用应用身份更新，并重新查询核验。 | 确认后更新：<br>`contract-cli mdm fixed-exchange-rate update --profile contract --as app --input-file fixed-exchange-rate.json`<br><br>按同一组参数回读：<br>`contract-cli mdm fixed-exchange-rate get --profile contract --as app --source-currency <source-currency> --target-currency <target-currency> --effective-date <effective-date>` |


### 12. 主数据文件下载

#### 场景 12.1：下载主数据记录中的附件

**授权支持：**应用

Agent 使用真实主数据文件 ID 下载，显式指定输出位置。该 skill 仅支持应用身份；合同正文和合同附件走合同 skill 的文件能力，不混用文件归属和入口。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请使用已配置的应用身份，把这条主数据记录中已确认文件 ID 的附件下载到我指定的目录。<br>完成后给出本地文件位置和下载结果。 | `contract-cli mdm file download <file-id> --profile contract --as app --output-file ./主数据附件.pdf` |


### 13. 付款申请、计划与记录

#### 场景 13.1：查询合同的付款申请

**授权支持：**应用

Agent 使用指定合同查询付款申请列表与详情；付款计划可按另行明确的条件搜索。结果以真实返回字段为准，不用申请状态推断已经完成银行划款。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请使用已配置的应用身份，查询合同“&lt;contract-id&gt;”的付款申请。<br>整理申请列表和返回的状态，继续读取我选中的申请详情。 | 查询申请列表：<br>`contract-cli payment list --contract <contract-id> --profile contract --as app --page-size 20`<br><br>读取所选申请：<br>`contract-cli payment get <payment-id> --contract <contract-id> --profile contract --as app` |


#### 场景 13.2：登记已确认的付款记录

**授权支持：**应用

该 skill 还支持付款申请和记录更新、付款计划同步及按计划查询记录，全部使用应用身份。它管理申请、计划和记录，不执行银行打款或退款。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请根据我提供的实际付款资料，准备指定合同和付款申请的付款记录。<br>先展示关联对象及填写内容，确认后创建，再查询核验结果。 | 确认后创建记录：<br>`contract-cli payment record create --contract <contract-id> --payment <payment-id> --profile contract --as app --input-file payment-record.json`<br><br>按创建结果回读：<br>`contract-cli payment record get <record-id> --contract <contract-id> --payment <payment-id> --profile contract --as app` |


### 14. 审批矩阵维护

#### 场景 14.1：将业务规则整理为矩阵导入计划

**授权支持：**个人 / 应用

Agent 支持个人或应用身份维护矩阵。文档、表格和 Excel 内容先由 Agent 转换为 CLI 支持的输入，再生成计划；不是直接导入任意 Excel。矩阵最多 2000 行，总列最多 12，其中优先级和备注占 2 列，条件与结果列合计最多 10 列。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请根据我提供的采购审批规则，整理目标矩阵的条件列、结果列和规则行。<br>先核实人员、部门及现有结构，展示导入计划；确认后写入并核验，暂不发布。 | 先生成不写业务数据的计划：<br>`contract-cli rule table import plan --profile contract --as user --product-id contract --group-id approve_matrix --table-id <table-id> --input-file matrix-import.json`<br><br>计划获确认后执行：<br>`contract-cli rule table import apply --profile contract --as user --plan-id <plan-id>` |


#### 场景 14.2：核验导入进度，再发布矩阵

**授权支持：**个人 / 应用

导入支持暂停、恢复、核验和取消后续执行，逐行写入，不保证整批原子成功，取消不回滚已有数据。计划有效期 24 小时，身份、结构或数据变化后可能需重新生成。发布是独立步骤；当前不能保证完整矩阵原样复制，按可见信息重建需补齐并确认缺失配置。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请查看刚才的矩阵导入进度，先核验结果不确定的行。<br>展示剩余内容，按我确认的范围继续；全部核验后展示发布影响，确认后再发布。 | 查看进度：<br>`contract-cli rule table import get --profile contract --as user --plan-id <plan-id>`<br><br>先核验未确认的行：<br>`contract-cli rule table import verify --profile contract --as user --plan-id <plan-id>`<br><br>剩余导入完成且获授权后预发布：<br>`contract-cli rule table pre-release --profile contract --as user --product-id contract --group-id approve_matrix --table-id <table-id>`<br><br>确认发布影响后正式发布：<br>`contract-cli rule table release --profile contract --as user --product-id contract --group-id approve_matrix --table-id <table-id>` |


### 15. 事件出口 IP 查询

#### 场景 15.1：为事件通知准备网络白名单资料

**授权支持：**应用

Agent 查询并汇总出口 IP；分页未完成时说明清单不完整。该 skill 不自动修改防火墙或其他网络配置。

| 用户输入提示词 | Agent 会使用的命令 |
| --- | --- |
| 请使用已配置的应用身份，查询智书合同事件通知的出口 IP 清单。<br>读取完整分页，整理给网络管理员配置白名单使用。 | 读取第一页：<br>`contract-cli event outbound-ip list --profile contract --as app --page-size 20`<br><br>仅在仍有下一页时：<br>`contract-cli event outbound-ip list --profile contract --as app --page-size 20 --page-token <page-token>` |

## 能力地图

下表按用户场景的顺序，统一列出中文简称、Skill 和授权支持。“仅个人授权”或“仅应用授权”表示当前业务命令支持的身份；“两种授权”表示两种身份均有可用操作，标注“按操作区分”时，具体能力不能互相替代。“授权与身份管理”负责授权本身，其他 Skill 复用对应身份。操作差异见下方速查表。

| 中文简称 / Skill | 能力范围 | 对应用户场景 | 授权支持与限制 |
| --- | --- | --- | --- |
| **授权与身份管理**<br>`auth` | 授权、状态查询、退出及默认身份切换 | 1.1 完成授权后继续任务 | **统一授权入口**<br>分别管理个人授权与应用授权；完成一种不会自动获得另一种 |
| **合同搜索**<br>`contract-cli-contract-search` | 关键词、编号、文本、日期、状态、金额、对象及自定义字段筛选，分页 | 2.1 关键词查找；2.2 组合筛选 | **两种授权，按操作区分**<br>搜索均支持；字段发现仅个人，V2 搜索仅应用；查询参数不同 |
| **合同管理**<br>`contract-cli-contract` | 详情与文件下载；文件/模板创建及变更、终止申请；提交/重提、流程操作、个人任务与评论；归档材料更新与打印；签署及认证链接、纸质签切换；分享授权、协商、字段选项、枚举、用户组同步与草稿删除 | 3.1 详情；3.2～3.4 创建与变更/终止；3.5～3.7 审批；3.8 归档与打印；3.9 签署；3.10 分享与授权；3.11 协商；3.12～3.14 配置、分组与草稿清理 | **两种授权，按操作区分**<br>详情、创建等均支持；个人任务/评论和枚举仅个人；提交/重提、归档更新、签署等仅应用 |
| **人员查询**<br>`contract-cli-employee` | 按姓名或部门查询人员及状态，取得真实候选 | 4.1 同名人员消歧 | **仅个人授权** |
| **部门查询**<br>`contract-cli-department` | 查询部门候选、层级和状态 | 5.1 定位业务部门 | **仅个人授权** |
| **合同智能审查**<br>`contract-cli-review` | 单份外部文件与系统合同审查，主体提取、任务续查、风险展示 | 6.1 外部文件审查；6.2 系统合同审查 | **仅个人授权**<br>还需具备相应合同、正文和审查权限 |
| **审查规则维护**<br>`contract-cli-review-rule` | 查询及单项维护审查清单、规则和分组；内置清单只读，删除受限制 | 7.1 查询审查标准；7.2 补充审查规则 | **仅个人授权**<br>维护操作还受清单/规则权限及内置配置限制 |
| **交易方管理**<br>`contract-cli-mdm-vendor` | 交易方查询、创建、局部更新、启停；另有应用身份全量更新/分页和证件查询 | 8.1 查询交易方；8.2 局部更新资料 | **两种授权，按操作区分**<br>查询、创建和局部更新均支持；启停仅个人；全量更新/分页及证件查询仅应用 |
| **法人主体管理**<br>`contract-cli-mdm-legal` | 法人主体列表及详情、按编码查询、创建与更新 | 9.1 确认签约主体；9.2 创建法人主体 | **两种授权，按操作区分**<br>列表和按 ID 查询均支持；按编码查询、创建和更新仅应用 |
| **主数据字段配置**<br>`contract-cli-mdm-fields` | 查询主数据字段定义，为填写资料提供依据 | 10.1 核对字段要求 | **两种授权，按操作区分**<br>交易方/法人字段均支持；交易方风险字段仅个人 |
| **固定汇率维护**<br>`contract-cli-mdm-exchange` | 固定汇率查询与更新 | 11.1 查汇率；11.2 更新配置 | **仅应用授权** |
| **主数据文件下载**<br>`contract-cli-mdm-file` | 根据主数据文件 ID 下载附件 | 12.1 下载主数据附件 | **仅应用授权** |
| **付款管理**<br>`contract-cli-payment` | 付款申请、付款计划和付款记录的查询与维护 | 13.1 查询付款申请；13.2 登记付款记录 | **仅应用授权** |
| **审批矩阵维护**<br>`contract-cli-rule` | 矩阵及列结构维护、规则行查询/维护/复制、导入计划与执行、预发布和发布 | 14.1 整理并导入规则；14.2 核验和发布 | **两种授权**<br>个人身份用于合同产品，并需规则管理权限；操作按所选身份留痕 |
| **事件出口 IP 查询**<br>`contract-cli-event` | 分页查询事件出口 IP | 15.1 准备白名单资料 | **仅应用授权** |

## 身份和操作范围速查

本表的“支持”只表示 v1.9.5 已开放该身份的命令路径，不代表当前账号或应用已获业务权限。授权有效后，仍要核对合同可见范围、角色权限及业务状态。缺少授权时完成对应授权；权限不足时按错误提示由管理员处理；命令不支持该身份时说明限制，不通过反复登录或自动切换身份尝试执行。

**同一个 Skill 内也可能需要不同授权。** 例如，个人授权支持创建合同和办理自己的审批待办；提交或重新提交合同审批需要应用授权。一个 Skill 同时支持两种授权，不等于其中每条命令都能使用两种身份。

**跨 Skill 的任务按步骤检查授权。** 例如，用应用身份分享合同时，如果还需按姓名查找接收人，“人员查询”需要个人授权；应用授权不会自动解锁人员查询。已有已核实的人员 ID 时，可按分享接口要求准备参数；缺少必要授权或真实对象时，先说明缺口。

**页面登录与 CLI 授权分别生效。** Agent 返回的合同、审查、签署或认证链接，仍按对应页面的登录和权限要求打开；取得链接不代表已经拥有网页访问权或完成认证、签署。

| 能力 / 命令族 | 个人授权 user | 应用授权 app | 说明 |
| --- | --- | --- | --- |
| `contract search` | 支持 | 支持 | 两者请求契约不同；user 有页面关键词、对象及字段引导 |
| `contract search-fields`、`employee list`、`department list` | 支持 | 不支持 | 搜索字段发现和人员、部门候选 |
| `contract search-v2` | 不支持 | 支持 | 编号模糊/批量查询，不能套用 user 请求体 |
| 合同详情、正文、创建、分类、模板、用户组同步 | 支持 | 支持 | 模板实例创建也支持两种身份 |
| 合同文件上传、下载 | 支持 | 支持 | user 下载必须关联合同；评论/审批附件使用规定用途 |
| 合同枚举 `contract enum list` | 支持 | 不支持 | 用于获取可用枚举 |
| 提交、重提、更新、删除、打印合同 | 不支持 | 支持 | 状态和请求字段按对应操作核验 |
| 已有流程操作 `contract approval start` | 不支持 | 支持 | 支持的审批动作受节点和权限限制；不是新建流程入口 |
| 流程详情 `contract approval get` | 支持 | 支持 | 需要流程实例 ID |
| 审批评论、个人任务查询和通过/拒绝 | 支持 | 不支持 | 需要当前用户任务权限；盖章/归档节点除外 |
| `review` 智能审查及清单/规则/分组 | 支持 | 不支持 | 按当前正文、业务场景和权限校验 |
| 交易方 list/get/create/patch | 支持 | 支持 | 两种身份参数与接口不同 |
| 交易方 enable/disable | 支持 | 不支持 | 启停操作 |
| 交易方 update/list-all/query-by-cert | 不支持 | 支持 | 旧版全量更新、全量分页、证件查询 |
| 法人主体 list/get（按 ID） | 支持 | 支持 | 先消歧再选定主体 |
| 法人主体 get --code/create/update | 不支持 | 支持 | 创建/更新按接口要求提供操作人 |
| 主数据字段配置 `mdm fields list` | 支持 | 支持 | user 还支持 vendor_risk；app 仅交易方和法人主体 |
| 付款、固定汇率、主数据文件、事件 IP | 不支持 | 支持 | 对应 payment、mdm、event 命令 |
| 合同字段选项、流程字段、签署、授权、实名认证链接、分享及协商命令 | 不支持 | 支持 | 智能审查的个人协商取数属于独立 review 路径 |
| 审批矩阵 `rule` / `approval-matrix` | 支持 | 支持 | 查询、结构维护、规则行、导入和发布；个人用于合同产品且需规则管理权限，按所选身份留痕 |

## 给实施和集成同学的命令示例

下面是入口示例，`<...>` 必须替换为真实值；JSON 文件按对应 skill 的字段文档准备，不直接执行占位内容。

```bash
# 合同查询与字段发现
contract-cli contract search --profile contract --as user --input-file search-user.json
contract-cli contract search-fields --keyword "项目区域" --profile contract --as user --output json
contract-cli employee list --name "张三" --profile contract --as user
contract-cli department list --name "法务" --profile contract --as user
contract-cli contract get <contract-id> --profile contract --as user
contract-cli contract text <contract-id> --profile contract --as user --full-text

# 文件与审批
contract-cli contract download-file <file-id> --contract <contract-id> --profile contract --as user --output-file ./合同附件.pdf
contract-cli contract approval task list --profile contract --as user --task-type todo --page-size 20
contract-cli contract approval get <process-instance-id> --profile contract --as user
contract-cli contract approval comment list <process-instance-id> --profile contract --as user

# 继续已有系统合同审查，复用原任务真实 task-id 和 business-id
contract-cli review task result --task-id <task-id> --business-id <business-id> --profile contract --as user

# 合同创建：参数已完整且已授权执行后使用
contract-cli contract category list --profile contract --as user --lang zh-CN
contract-cli contract create --profile contract --as user --input-file contract-create.json

# 应用身份的付款查询
contract-cli payment list --contract <contract-id> --profile contract --as app --page-size 10
```

审批矩阵的计划、写入和发布，以及审查配置的写操作，请按对应 skill 的完整流程执行。参数有疑问时，先查当前版本的 `--help` 和该 skill 的 `references/`。

## Agent 的结果应该包含什么

查询结果要列明实际条件、范围、分页情况和必要业务字段；文件任务要给出可访问的本地产物及成功/失败清单；审查任务要展示真实风险结果与任务状态；写操作要返回真实执行结果，并使用可用查询核验。

缺少关键业务信息时，Agent 先查真实候选，再请你消歧。请求失败不等于无数据，写入超时不等于未写入，未查完不等于总量已确定。结果不确定时先核实，不能直接重复创建、审批或发布。

本次对外范围以已发布的结构化命令及本文介绍的 15 个 skills 为准。通用 `api call` 尚未开放；没有在该 Release 分发的到期提醒、我的草稿、签署排行等独立业务 skills，不作为内置能力宣传。本文也不承诺持续监控、自动定时通知或所有合同状态下的审查。

## 常见问题

**安装完成后 Agent 仍不会用怎么办？**

检查实际 CLI 路径和版本、当前客户端的技能目录、同名旧 skill 是否已覆盖，以及新任务是否加载新版。`skills list` 是内置目录，不是客户端加载报告。

**个人已授权，为什么某项功能仍不能执行？**

先区分身份不支持和业务权限不足。应用专属操作需要已配置的应用授权；个人支持的操作仍受合同、任务、审查清单等权限约束。Agent 应说明具体缺口，不自动换身份。

**合同文本搜索与读取正文有什么区别？**

搜索用于找合同，可按正文或页面文本范围查询；`contract text` 用于读取已知合同正文。正文摘要、一般问答与调用智能审查服务也属于不同任务。

**为什么审查前要选择主体、强度和清单？**

它们决定保护哪一方、采用何种谈判立场和检查哪些规则。主体从当前文件中提取，不能只根据文件名或登记交易方猜测。

**可以把旧版 Aily 手动回调流程直接套到所有 Agent 吗？**

不能作为通用流程。本文优先按当前内置 auth skill 的宿主和授权规则接入，未覆盖宿主需要另行验证；不让用户通过聊天粘贴临时凭据来替代官方授权流程。

## 版本与参考

本修订稿以 [v1.9.5 Release](https://github.com/qfeius/contract-cli/releases/tag/v1.9.5) 和对应源码为准。

- [版本对应的全部 skills](https://github.com/qfeius/contract-cli/tree/v1.9.5/skills)
- [合同搜索 skill](https://github.com/qfeius/contract-cli/blob/v1.9.5/skills/contract-cli-contract-search/SKILL.md)
- [智能审查 skill](https://github.com/qfeius/contract-cli/blob/v1.9.5/skills/contract-cli-review/SKILL.md)
- [审查规则 skill](https://github.com/qfeius/contract-cli/blob/v1.9.5/skills/contract-cli-review-rule/SKILL.md)
- [审批矩阵 skill](https://github.com/qfeius/contract-cli/blob/v1.9.5/skills/contract-cli-rule/SKILL.md)
- [授权 skill](https://github.com/qfeius/contract-cli/blob/v1.9.5/skills/auth/SKILL.md)

GitHub 提供代码和版本参考入口；仓库当前标记为 `UNLICENSED`，不使用“开源许可”描述。

| 文档状态 | 日期 | 本次修订 |
| --- | --- | --- |
| 修订副本，基于原文档 v1.7 | 2026-10-13 | 介绍 CLI v1.9.5 面向用户的 15 个 skills；按能力范围配置场景，合同管理展开为 14 个场景，能力地图逐项对应；补充已发布能力并更新授权及身份边界。本稿为独立修订副本，原飞书文档保持不变。 |
