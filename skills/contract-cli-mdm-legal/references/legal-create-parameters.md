# mdm legal create Parameters

本页专用于 `contract-cli mdm legal create`。

- 接口：`POST /open-apis/mdm/v1/legal_entities`
- 身份：仅 `app`
- 请求体：JSON 必填，`--input-file` 与 `--data` 二选一
- 官方 OpenAPI：[创建法人实体](https://docs.qfei.cn/373516226e0.md)
- 校验基线：`contract-cli` 当前分支；CLM `master@334c18fc2e` 存在对应实现时，以 Controller、DTO 和业务校验补充官方规格。

## 目录

- [CLI 参数映射](#cli-参数映射)
- [请求体字段](#请求体字段)
- [枚举与约束](#枚举与约束)
- [示例](#示例)
- [来源差异说明](#来源差异说明)

## CLI 参数映射

| 参数 | 请求位置 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- | --- |
| --user-id | $query.user_id | string | 必填（CLI 本地校验） | 用户id，如果使用的是 tenant_access_token，那么需要填入user_id<br>示例值："123123123123" |
| --user-id-type | $query.user_id_type | string | 可选，默认 `user_id` | 用户 ID 类型，参考 用户身份体系 |
| --input-file | $body | JSON file | 二选一必填 | 从文件读取 JSON；与 `--data` 互斥。 |
| --data | $body | JSON string | 二选一必填 | 内联 JSON；与 `--input-file` 互斥。 |
| --profile | 本地上下文 | string | 可选 | 不传时使用当前 profile。 |
| --as | 本地上下文 | enum | 可选 | 仅支持 `app`；不传时使用 profile 默认身份。 |
| --output | CLI 输出 | enum | 可选 | `json`、`yaml` 或 `table`；默认 `json`。 |
| --raw | CLI 输出 | boolean | 可选 | 原样输出服务端响应 body。 |

## 请求体字段

字段名、类型和服务端必填性来自官方 OpenAPI；“CLI 必填/禁止”是结构化命令的额外本地校验。父对象可选时，其内部必填字段标记为“父对象存在时必填”。

| JSON 路径 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- |
| legalEntity | string | CLI 禁止 | 法人实体编码(根据配置会有不同的生成规则)<br>示例值："L00002002"<br>数据校验规则：最大长度：32 字符 |
| legalEntityText | string | 可选 | 法人实体名称<br>示例值："法人22"<br>数据校验规则：<br>最大长度：128 字符 |
| shortText | string | 可选 | 法人实体英文名称<br>示例值："legal_person"<br>数据校验规则：<br>最大长度：40 字符 |
| certificationType | string | 可选 | 证件类型<br>示例值："0"<br>可选值有：<br>- 0：统一社会信用代码(中国大陆)<br>- 1：中国大陆居民身份证(中国大陆)<br>- 2：注册号(海外)<br>- 3：税号(海外)<br>- 4：驾驶证(海外)<br>- 5：身份证(海外)<br>- 6：护照<br>- 8：港澳居民往来大陆通行证<br>- 9：台湾居民往来大陆通行征<br>- 10：香港永久性居民身份证<br>- 11：澳门特别行政区永久性居民身份证<br>- 12：台湾身份证<br>- 13：外国人永久居留证 |
| certificationId | string | 可选 | 证件id<br>示例值："91310120MA1H23N81AX"<br>数据校验规则：<br>最大长度：64 字符 |
| legalPerson | string | 可选 | 法人<br>示例值："张三"<br>数据校验规则：<br>最大长度：50 字符 |
| country | string | 可选 | 国家<br>示例值："CN"<br>数据校验规则：<br>最大长度：16 字符 |
| province | string | 可选 | 省份<br>示例值："MDPS00004000"<br>数据校验规则：<br>最大长度：16 字符 |
| city | string | 可选 | 城市<br>示例值："MDCY00006000"<br>数据校验规则：<br>最大长度：16 字符 |
| address | string | 可选 | 地址<br>示例值："地址"<br>数据校验规则：<br>最大长度：128 字符 |
| taxpayerType | string | 可选 | 纳税人类型<br>示例值："1"<br>可选值有：<br>- 1：一般纳税人<br>- 2：小规模纳税人 |
| telephone | string | 可选 | 联系电话<br>示例值："010-58341796"<br>数据校验规则：<br>最大长度：32 字符 |
| bankId | string | 可选 | 银行内部Id<br>示例值："MDBK00072319"<br>数据校验规则：<br>最大长度：100 字符 |
| bankName | string | 可选 | 开户银行名称<br>示例值："中原银行股份有限公司南阳华瑞支行"<br>数据校验规则：<br>最大长度：64 字符 |
| bankAccount | string | 可选 | 开户行账号<br>示例值："644666446"<br>数据校验规则：<br>最大长度：32 字符 |
| status | integer | 必填 | 状态<br>示例值：1<br>可选值有：<br>- 1：有效<br>- 0：无效 |
| appendix | array<object> | 可选 | 附件列表<br>数据校验规则：最大长度：10 |
| appendix[].tenantId | string | 可选 | — |
| appendix[].fileId | string | 可选 | 文件id(文件的唯一标识)<br>示例值："5c7237e821a8409d9b8b2e1041cdc9a4" |
| appendix[].fileName | string | 可选 | 文件名称<br>示例值："附件" |
| appendix[].fileType | string | 可选 | 文件类型<br>示例值："DOX"<br>可选值有：<br>- DOC：DOC<br>- DOCX：DOCX<br>- XLS：XLS<br>- XLSX：XLSX<br>- PNG：PNG<br>- JPG：JPG<br>- JPEG：JPEG<br>- PDF：PDF<br>- ZIP：ZIP<br>- RAR：RAR |
| appendix[].fileSize | integer | 可选 | 文件大小<br>示例值：1024 |
| appendix[].downloadUrl | string | 可选 | 文件下载地址<br>示例值："http://download.com/xxxxx" |
| legalEntityBanks | array<object> | 可选 | 银行账户列表<br>数据校验规则：<br>最大长度：100 |
| legalEntityBanks[].companyCode | string | 可选 | 公司编码<br>示例值："1002"<br>数据校验规则：<br>最大长度：50 字符 |
| legalEntityBanks[].bankId | string | 可选 | 银行Id<br>示例值："MDBK00131739"<br>数据校验规则：<br>最大长度：100 字符 |
| legalEntityBanks[].bankCode | string | 可选 | 银联号<br>示例值："001755053005"<br>数据校验规则：<br>最大长度：64 字符 |
| legalEntityBanks[].bankName | string | 可选 | 银行名称<br>示例值："中国人民银行丽江市中心支行"<br>数据校验规则：<br>最大长度：64 字符 |
| legalEntityBanks[].bankAcronym | string | 可选 | 总行英文缩写<br>示例值："PBC"<br>数据校验规则：<br>最大长度：100 字符 |
| legalEntityBanks[].country | string | 可选 | 国家<br>示例值："CN"<br>数据校验规则：<br>最大长度：16 字符 |
| legalEntityBanks[].accountName | string | 可选 | 账户名称<br>示例值："账户名称"<br>数据校验规则：<br>最大长度：64 字符 |
| legalEntityBanks[].bankAccount | string | 可选 | 银行账号<br>示例值："122907287xxxxx9"<br>数据校验规则：<br>最大长度：64 字符 |
| legalEntityBanks[].swiftCode | string | 可选 | 银行SWIFT编码<br>示例值："CMBCCNBS"<br>数据校验规则：<br>最大长度：100 字符 |
| legalEntityBanks[].bankControlCode | string | 可选 | 银行控制码<br>示例值："40001xxxxxxx00313261"<br>数据校验规则：<br>最大长度：100 字符 |
| legalEntityBanks[].extendInfo | array<object> | 可选 | 扩展字段相关信息列表<br>数据校验规则：<br>最大长度：100 |
| legalEntityBanks[].extendInfo[].fieldType | integer | 父对象存在时必填 | 字段类型<br>示例值：0<br>可选值有：<br>- 0：单行文本框<br>- 1：多行文本框<br>- 2：数字<br>- 3：单选框<br>- 4：多选框<br>- 5：下拉单选<br>- 6：下拉多选<br>- 7：日期<br>- 8：日期区间<br>- 12：附件 |
| legalEntityBanks[].extendInfo[].fieldValue | string | 可选 | 字段类型为 单行文本框(0)、多行文本框(1)、单选框(3)、下拉单选框(5) 时的值<br>示例值："文本值" |
| legalEntityBanks[].extendInfo[].options | array<string> | 可选 | 字段类型为 多选框(4) 下拉多选(6) 时的值<br>示例值：[""""]<br>数据校验规则：<br>最大长度：100 |
| legalEntityBanks[].extendInfo[].num | number | 可选 | 字段类型为 数字(2) 时的值<br>示例值：1.11 |
| legalEntityBanks[].extendInfo[].date | string | 可选 | 字段类型是 日期(7)时候的值<br>示例值："2021-10-14" |
| legalEntityBanks[].extendInfo[].rangeDate | array<string> | 可选 | 字段类型是 日期区间(8) 时候的值 数组长度为2 0-startTime 1-endTime<br>示例值：[""""]<br>数据校验规则：长度范围：2 ～ 2 |
| legalEntityBanks[].extendInfo[].fieldCode | string | 父对象存在时必填 | 字段编码<br>示例值："X00000001" |
| legalEntityBanks[].extendInfo[].appendix | array<object> | 可选 | 附件列表 字段类型是 附件(12) 时候的值<br>数据校验规则：<br>最大长度：10 |
| legalEntityBanks[].extendInfo[].appendix[].tenantId | string | 可选 | — |
| legalEntityBanks[].extendInfo[].appendix[].fileId | string | 可选 | 文件id(文件的唯一标识)<br>示例值："5c7237e821a8409d9b8b2e1041cdc9a4" |
| legalEntityBanks[].extendInfo[].appendix[].fileName | string | 可选 | 文件名称<br>示例值："附件" |
| legalEntityBanks[].extendInfo[].appendix[].fileType | string | 可选 | 文件类型<br>示例值："DOX"<br>可选值有：<br>- DOC：DOC<br>- DOCX：DOCX<br>- XLS：XLS<br>- XLSX：XLSX<br>- PNG：PNG<br>- JPG：JPG<br>- JPEG：JPEG<br>- PDF：PDF<br>- ZIP：ZIP<br>- RAR：RAR |
| legalEntityBanks[].extendInfo[].appendix[].fileSize | integer | 可选 | 文件大小<br>示例值：1024 |
| legalEntityBanks[].extendInfo[].appendix[].downloadUrl | string | 可选 | 文件下载地址<br>示例值："http://download.com/xxxxx" |
| legalEntityBanks[].ibanAccount | string | 可选 | IBAN账号<br>示例值："6446777"<br>数据校验规则：<br>最大长度：100 字符 |
| extendInfo | array<object> | 可选 | 扩展字段相关信息列表<br>数据校验规则：<br>最大长度：100 |
| extendInfo[].fieldType | integer | 父对象存在时必填 | 字段类型<br>示例值：0<br>可选值有：<br>- 0：单行文本框<br>- 1：多行文本框<br>- 2：数字<br>- 3：单选框<br>- 4：多选框<br>- 5：下拉单选<br>- 6：下拉多选<br>- 7：日期<br>- 8：日期区间<br>- 12：附件 |
| extendInfo[].fieldValue | string | 可选 | 字段类型为 单行文本框(0)、多行文本框(1)、单选框(3)、下拉单选框(5) 时的值<br>示例值："文本值" |
| extendInfo[].options | array<string> | 可选 | 字段类型为 多选框(4) 下拉多选(6) 时的值示例值：[""""]<br>数据校验规则：<br>最大长度：100 |
| extendInfo[].num | number | 可选 | 字段类型为 数字(2) 时的值<br>示例值：1.11 |
| extendInfo[].date | string | 可选 | 字段类型是 日期(7)时候的值<br>示例值："2021-10-14" |
| extendInfo[].rangeDate | array<string> | 可选 | 字段类型是 日期区间(8) 时候的值 数组长度为2 0-startTime 1-endTime示例值：[""""]<br>数据校验规则：<br>长度范围：2 ～ 2 |
| extendInfo[].fieldCode | string | 父对象存在时必填 | 字段编码<br>示例值："X00000001" |
| extendInfo[].appendix | array<object> | 可选 | 附件列表 字段类型是 附件(12) 时候的值<br>数据校验规则：<br>最大长度：10 |
| extendInfo[].appendix[].tenantId | string | 可选 | — |
| extendInfo[].appendix[].fileId | string | 可选 | 文件id(文件的唯一标识)<br>示例值："5c7237e821a8409d9b8b2e1041cdc9a4" |
| extendInfo[].appendix[].fileName | string | 可选 | 文件名称<br>示例值："附件" |
| extendInfo[].appendix[].fileType | string | 可选 | 文件类型<br>示例值："DOX"<br>可选值有：<br>- DOC：DOC<br>- DOCX：DOCX<br>- XLS：XLS<br>- XLSX：XLSX<br>- PNG：PNG<br>- JPG：JPG<br>- JPEG：JPEG<br>- PDF：PDF<br>- ZIP：ZIP<br>- RAR：RAR |
| extendInfo[].appendix[].fileSize | integer | 可选 | 文件大小<br>示例值：1024 |
| extendInfo[].appendix[].downloadUrl | string | 可选 | 文件下载地址<br>示例值："http://download.com/xxxxx" |
| accountName | string | 可选 | 账户名称<br>示例值："账户名称"<br>数据校验规则：<br>最大长度：64 字符 |
| bankControlCode | string | 可选 | 银行控制码<br>示例值："40001xxxxxxx00313261"<br>数据校验规则：<br>最大长度：100 字符 |
| swiftCode | string | 可选 | 银行SWIFT编码<br>示例值："CMBCCNBS"<br>数据校验规则：<br>最大长度：100 字符 |
| clearingAccount | string | 可选 | 清算科目编码<br>示例值："ASD123" |

## 枚举与约束

- 根据该接口，来创建一个法人实体，字段是否必填是根据后台动态配置的，如想获取配置，主数据的配置开放文档里获取。参数均采用驼峰式
- 字段必填性受租户动态配置影响；调用前使用 `mdm fields list --biz-line legal_entity`。
- CLI 禁止请求体包含 `legalEntity` 或 `legal_entity`。
- `legalEntity`（string，CLI 禁止）：法人实体编码(根据配置会有不同的生成规则)<br>示例值："L00002002"<br>数据校验规则：最大长度：32 字符
- `legalEntityText`（string，可选）：法人实体名称<br>示例值："法人22"<br>数据校验规则：<br>最大长度：128 字符
- `shortText`（string，可选）：法人实体英文名称<br>示例值："legal_person"<br>数据校验规则：<br>最大长度：40 字符
- `certificationType`（string，可选）：证件类型<br>示例值："0"<br>可选值有：<br>- 0：统一社会信用代码(中国大陆)<br>- 1：中国大陆居民身份证(中国大陆)<br>- 2：注册号(海外)<br>- 3：税号(海外)<br>- 4：驾驶证(海外)<br>- 5：身份证(海外)<br>- 6：护照<br>- 8：港澳居民往来大陆通行证<br>- 9：台湾居民往来大陆通行征<br>- 10：香港永久性居民身份证<br>- 11：澳门特别行政区永久性居民身份证<br>- 12：台湾身份证<br>- 13：外国人永久居留证
- `certificationId`（string，可选）：证件id<br>示例值："91310120MA1H23N81AX"<br>数据校验规则：<br>最大长度：64 字符
- `legalPerson`（string，可选）：法人<br>示例值："张三"<br>数据校验规则：<br>最大长度：50 字符
- `country`（string，可选）：国家<br>示例值："CN"<br>数据校验规则：<br>最大长度：16 字符
- `province`（string，可选）：省份<br>示例值："MDPS00004000"<br>数据校验规则：<br>最大长度：16 字符
- `city`（string，可选）：城市<br>示例值："MDCY00006000"<br>数据校验规则：<br>最大长度：16 字符
- `address`（string，可选）：地址<br>示例值："地址"<br>数据校验规则：<br>最大长度：128 字符
- `taxpayerType`（string，可选）：纳税人类型<br>示例值："1"<br>可选值有：<br>- 1：一般纳税人<br>- 2：小规模纳税人
- `telephone`（string，可选）：联系电话<br>示例值："010-58341796"<br>数据校验规则：<br>最大长度：32 字符
- `bankId`（string，可选）：银行内部Id<br>示例值："MDBK00072319"<br>数据校验规则：<br>最大长度：100 字符
- `bankName`（string，可选）：开户银行名称<br>示例值："中原银行股份有限公司南阳华瑞支行"<br>数据校验规则：<br>最大长度：64 字符
- `bankAccount`（string，可选）：开户行账号<br>示例值："644666446"<br>数据校验规则：<br>最大长度：32 字符
- `status`（integer，必填）：状态<br>示例值：1<br>可选值有：<br>- 1：有效<br>- 0：无效
- `appendix`（array<object>，可选）：附件列表<br>数据校验规则：最大长度：10
- `appendix[].fileId`（string，可选）：文件id(文件的唯一标识)<br>示例值："5c7237e821a8409d9b8b2e1041cdc9a4"
- `appendix[].fileType`（string，可选）：文件类型<br>示例值："DOX"<br>可选值有：<br>- DOC：DOC<br>- DOCX：DOCX<br>- XLS：XLS<br>- XLSX：XLSX<br>- PNG：PNG<br>- JPG：JPG<br>- JPEG：JPEG<br>- PDF：PDF<br>- ZIP：ZIP<br>- RAR：RAR
- `legalEntityBanks`（array<object>，可选）：银行账户列表<br>数据校验规则：<br>最大长度：100
- `legalEntityBanks[].companyCode`（string，可选）：公司编码<br>示例值："1002"<br>数据校验规则：<br>最大长度：50 字符
- `legalEntityBanks[].bankId`（string，可选）：银行Id<br>示例值："MDBK00131739"<br>数据校验规则：<br>最大长度：100 字符
- `legalEntityBanks[].bankCode`（string，可选）：银联号<br>示例值："001755053005"<br>数据校验规则：<br>最大长度：64 字符
- `legalEntityBanks[].bankName`（string，可选）：银行名称<br>示例值："中国人民银行丽江市中心支行"<br>数据校验规则：<br>最大长度：64 字符
- `legalEntityBanks[].bankAcronym`（string，可选）：总行英文缩写<br>示例值："PBC"<br>数据校验规则：<br>最大长度：100 字符
- `legalEntityBanks[].country`（string，可选）：国家<br>示例值："CN"<br>数据校验规则：<br>最大长度：16 字符
- `legalEntityBanks[].accountName`（string，可选）：账户名称<br>示例值："账户名称"<br>数据校验规则：<br>最大长度：64 字符
- `legalEntityBanks[].bankAccount`（string，可选）：银行账号<br>示例值："122907287xxxxx9"<br>数据校验规则：<br>最大长度：64 字符
- `legalEntityBanks[].swiftCode`（string，可选）：银行SWIFT编码<br>示例值："CMBCCNBS"<br>数据校验规则：<br>最大长度：100 字符
- `legalEntityBanks[].bankControlCode`（string，可选）：银行控制码<br>示例值："40001xxxxxxx00313261"<br>数据校验规则：<br>最大长度：100 字符
- `legalEntityBanks[].extendInfo`（array<object>，可选）：扩展字段相关信息列表<br>数据校验规则：<br>最大长度：100
- `legalEntityBanks[].extendInfo[].fieldType`（integer，父对象存在时必填）：字段类型<br>示例值：0<br>可选值有：<br>- 0：单行文本框<br>- 1：多行文本框<br>- 2：数字<br>- 3：单选框<br>- 4：多选框<br>- 5：下拉单选<br>- 6：下拉多选<br>- 7：日期<br>- 8：日期区间<br>- 12：附件
- `legalEntityBanks[].extendInfo[].options`（array<string>，可选）：字段类型为 多选框(4) 下拉多选(6) 时的值<br>示例值：[""""]<br>数据校验规则：<br>最大长度：100
- `legalEntityBanks[].extendInfo[].rangeDate`（array<string>，可选）：字段类型是 日期区间(8) 时候的值 数组长度为2 0-startTime 1-endTime<br>示例值：[""""]<br>数据校验规则：长度范围：2 ～ 2
- `legalEntityBanks[].extendInfo[].appendix`（array<object>，可选）：附件列表 字段类型是 附件(12) 时候的值<br>数据校验规则：<br>最大长度：10
- `legalEntityBanks[].extendInfo[].appendix[].fileId`（string，可选）：文件id(文件的唯一标识)<br>示例值："5c7237e821a8409d9b8b2e1041cdc9a4"
- `legalEntityBanks[].extendInfo[].appendix[].fileType`（string，可选）：文件类型<br>示例值："DOX"<br>可选值有：<br>- DOC：DOC<br>- DOCX：DOCX<br>- XLS：XLS<br>- XLSX：XLSX<br>- PNG：PNG<br>- JPG：JPG<br>- JPEG：JPEG<br>- PDF：PDF<br>- ZIP：ZIP<br>- RAR：RAR
- `legalEntityBanks[].ibanAccount`（string，可选）：IBAN账号<br>示例值："6446777"<br>数据校验规则：<br>最大长度：100 字符
- `extendInfo`（array<object>，可选）：扩展字段相关信息列表<br>数据校验规则：<br>最大长度：100
- `extendInfo[].fieldType`（integer，父对象存在时必填）：字段类型<br>示例值：0<br>可选值有：<br>- 0：单行文本框<br>- 1：多行文本框<br>- 2：数字<br>- 3：单选框<br>- 4：多选框<br>- 5：下拉单选<br>- 6：下拉多选<br>- 7：日期<br>- 8：日期区间<br>- 12：附件
- `extendInfo[].options`（array<string>，可选）：字段类型为 多选框(4) 下拉多选(6) 时的值示例值：[""""]<br>数据校验规则：<br>最大长度：100
- `extendInfo[].rangeDate`（array<string>，可选）：字段类型是 日期区间(8) 时候的值 数组长度为2 0-startTime 1-endTime示例值：[""""]<br>数据校验规则：<br>长度范围：2 ～ 2
- `extendInfo[].appendix`（array<object>，可选）：附件列表 字段类型是 附件(12) 时候的值<br>数据校验规则：<br>最大长度：10
- `extendInfo[].appendix[].fileId`（string，可选）：文件id(文件的唯一标识)<br>示例值："5c7237e821a8409d9b8b2e1041cdc9a4"
- `extendInfo[].appendix[].fileType`（string，可选）：文件类型<br>示例值："DOX"<br>可选值有：<br>- DOC：DOC<br>- DOCX：DOCX<br>- XLS：XLS<br>- XLSX：XLSX<br>- PNG：PNG<br>- JPG：JPG<br>- JPEG：JPEG<br>- PDF：PDF<br>- ZIP：ZIP<br>- RAR：RAR
- `accountName`（string，可选）：账户名称<br>示例值："账户名称"<br>数据校验规则：<br>最大长度：64 字符
- `bankControlCode`（string，可选）：银行控制码<br>示例值："40001xxxxxxx00313261"<br>数据校验规则：<br>最大长度：100 字符
- `swiftCode`（string，可选）：银行SWIFT编码<br>示例值："CMBCCNBS"<br>数据校验规则：<br>最大长度：100 字符

## 示例

```bash
contract-cli mdm legal create --user-id <operator-user-id> --input-file request.json --profile contract --as app
```

官方请求体示例（动态字段接口仍须以当前租户配置为准）：

```json
"{\r\n    \"legalEntity\": \"L00002002\",\r\n    \"legalEntityText\": \"法人22\",\r\n    \"shortText\": \"legalPerson\",\r\n    \"certificationType\": \"0\",\r\n    \"certificationId\": \"\",\r\n    \"legalPerson\": \"张三\",\r\n    \"country\": \"CN\",\r\n    \"province\": \"MDPS00004000\",\r\n    \"city\": \"MDCY00006000\",\r\n    \"address\": \"地址\",\r\n    \"taxpayerType\": \"1\",\r\n    \"telephone\": \"\",\r\n    \"bankId\": \"MDBK00072319\",\r\n    \"bankName\": \"中原银行股份有限公司南阳华瑞支行\",\r\n    \"bankAccount\": \"644666446\",\r\n    \"status\": 1,\r\n    \"appendix\": [\r\n                    {\r\n                        \"tenantId\": \"6977354570259330000\",\r\n                        \"fileId\": \"609a128628ad4eaebd3063c59928a103\",\r\n                        \"fileName\": \"xxxx.xlsx\",\r\n                        \"fileType\": \"XLSX\",\r\n                        \"fileSize\": 13367,\r\n                        \"downloadUrl\": \"https://xxxx.qfei.cn/downloadxxxxxxxxxx\"\r\n                    }\r\n                ],\r\n    \"legalEntityBanks\": [\r\n        {\r\n            \"companyCode\": \"1002\",\r\n            \"bankId\": \"MDBK00131739\",\r\n            \"bankCode\": \"001755053005\",\r\n            \"bankName\": \"中国人民银行丽江市中心支行\",\r\n            \"bankAcronym\": \"PBC\",\r\n            \"country\": \"CN\",\r\n            \"accountName\": \"账户名称\",\r\n            \"bankAccount\": \"122907287xxxxx9\",\r\n            \"swiftCode\": \"CMBCCNBS\",\r\n            \"bankControlCode\": \"40001xxxxxxx00313261\",\r\n            \"extendInfo\": [\r\n                {\r\n                    \"fieldType\": 0,\r\n                    \"fieldValue\": \"文本值\",\r\n                    \"options\": [\r\n                        \"\"\r\n                    ],\r\n                    \"num\": 1.11,\r\n                    \"date\": \"2021-10-14\",\r\n                    \"rangeDate\": [\r\n                        \"\"\r\n                    ],\r\n                    \"fieldCode\": \"LXX0000001\",\r\n                    \"appendix\": [\r\n                                     {\r\n                                         \"tenantId\": \"6977354570259330000\",\r\n                                         \"fileId\": \"609a128628ad4eaebd3063c59928a103\",\r\n                                         \"fileName\": \"xxxx.xlsx\",\r\n                                         \"fileType\": \"XLSX\",\r\n                                         \"fileSize\": 13367,\r\n                                         \"downloadUrl\": \"https://xxxx.qfei.cn/downloadxxxxxxxxxx\"\r\n                                     }\r\n                                 ]\r\n                }\r\n            ],\r\n            \"ibanAccount\": \"6446777\"\r\n        }\r\n    ],\r\n    \"extendInfo\": [\r\n        {\r\n            \"fieldType\": 0,\r\n            \"fieldValue\": \"文本值\",\r\n            \"options\": [\r\n                \"\"\r\n            ],\r\n            \"num\": 1.11,\r\n            \"date\": \"2021-10-14\",\r\n            \"rangeDate\": [\r\n                \"\"\r\n            ],\r\n            \"fieldCode\": \"T0001\",\r\n            \"appendix\": [\r\n                             {\r\n                                 \"tenantId\": \"6977354570259330000\",\r\n                                 \"fileId\": \"609a128628ad4eaebd3063c59928a103\",\r\n                                 \"fileName\": \"xxxx.xlsx\",\r\n                                 \"fileType\": \"XLSX\",\r\n                                 \"fileSize\": 13367,\r\n                                 \"downloadUrl\": \"https://xxxx.qfei.cn/downloadxxxxxxxxxx\"\r\n                             }\r\n                         ]\r\n        }\r\n    ]\r\n    \"accountName\": \"12\",\r\n    \"bankControlCode\": \"12\",\r\n    \"swiftCode\": \"123\",\r\n    \"clearingAccount\": \"12333\"\r\n}"
```

## 来源差异说明

- 官方规格路径：`POST /open-apis/mdm/v1/legal_entities`
- 本页描述请求参数；响应 envelope 和输出格式遵循共享 Skill。
