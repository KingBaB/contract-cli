# 调用通道与参数约定

## 环境和授权

本 Skill 保留 test 下受控验证范围，不据此开放 prod 或 blue 业务写入。仓库当前 CLI 的生产构建只支持 prod；它不能执行这里的 test 配方。遇到该限制应说明不兼容并停止，不改地址、换凭证或改用 prod 执行原请求。只有已有且明确支持 test 的授权通道才可继续；不要求安装某个版本或指定来源包。

沿用同一环境、同一租户、同一用户授权。CLI 授权步骤按 [auth](../../auth/SKILL.md) 中实际支持的流程执行；授权文档不构成扩大环境范围的依据。无需用户在对话里提供 Token、Cookie、Session 或身份头。没有可用通道时说明缺少的能力及环境，不自行安装来源不明的包。

## 选择入口

先检查现有 MCP 的实际 schema 是否支持通用 Open Web 请求；否则检查 CLI `api call --help` 是否支持 `--channel web`、query 与输入文件。通道必须覆盖本场景的读取、写入与结果回查。能力检查不固定 CLI 版本，不因同名命令存在就认定兼容。

只有入口缺失或参数能力不兼容时才选择另一通道。业务拒绝、已发送请求或结果不确定时不换通道重放。所有业务调用限本 Skill 明确列出的 9 个 Open Web 端点，不枚举其他接口；MCP 未提供的 raw_body 等参数不能自行添加。

## CLI 调用形式

使用已确认的 profile 和 user 身份；以下是替换占位符后的调用形式，不是要求执行的测试请求：

```sh
contract-cli api call --channel web --profile <已确认profile> --as user GET /open-web/web/clm/employee/info
contract-cli api call --channel web --profile <已确认profile> --as user --query 'keywords=接收人姓名' --query 'includeInvalid=false' GET /open-web/web/clm/employee/findEmployeesByKeywords
contract-cli api call --channel web --profile <已确认profile> --as user --input-file <请求文件.json> POST /open-web/web/clm/permission/handover
```

- GET 参数使用 query，不拼接未知 URL；重复 query 的实际编码按命令帮助处理。query 中的布尔值使用 true/false。
- POST 使用 UTF-8 JSON 文件，避免 shell 拼接用户文本。body 中整数类型码保持整数、布尔保持布尔，Long ID 从响应到请求保留十进制字符串；不要经 JavaScript Number 转换。MCP 如需 raw_body 保真，必须确认实际 schema 支持。
- 内部 employeeId、contractId、groupId、assemblyId 与开放 user_id 含义不同，不混用；不发送 tenantId、employeeId 等额外身份覆盖字段，业务请求中的 targetTenantId 仅是已核实接收人的租户。
- 所有接口预期下游 body 为 code=200 且 success=true，再读取 data。HTTP/通道失败与业务失败分别判断；CLI 只返回 body 时不编造外层 HTTP 状态。data 为空不是交接成功。
- 只保存必要的本次对象、写前快照和核对结果，不记录凭证或无关个人信息。
