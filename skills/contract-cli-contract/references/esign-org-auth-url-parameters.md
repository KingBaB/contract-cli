# contract esign org-auth-url Parameters

本页专用于 `contract-cli contract esign org-auth-url`。

- 接口：`POST /open-apis/esign/auth/orgAuthUrl`
- 身份：仅 `app`
- 请求体：JSON 必填，`--input-file` 与 `--data` 二选一
- 官方 OpenAPI：[获取机构认证&授权页面链接](https://docs.qfei.cn/367694046e0.md)
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
| --input-file | $body | JSON file | 二选一必填 | 从文件读取 JSON；与 `--data` 互斥。 |
| --data | $body | JSON string | 二选一必填 | 内联 JSON；与 `--input-file` 互斥。 |
| --profile | 本地上下文 | string | 可选 | 不传时使用当前 profile。 |
| --as | 本地上下文 | enum | 可选 | 仅支持 `app`；不传时使用 profile 默认身份。 |
| --user-id-type | $query.user_id_type | string | 可选 | 不传时 CLI 默认发送 `user_id`。 |
| --user-id | $query.user_id | string | 可选 | 传入时透传；MDM create/update 除外。 |
| --output | CLI 输出 | enum | 可选 | `json`、`yaml` 或 `table`；默认 `json`。 |
| --raw | CLI 输出 | boolean | 可选 | 原样输出服务端响应 body。 |

## 请求体字段

字段名、类型和服务端必填性来自官方 OpenAPI；“CLI 必填/禁止”是结构化命令的额外本地校验。父对象可选时，其内部必填字段标记为“父对象存在时必填”。

| JSON 路径 | 类型 | 必填性 | 说明 |
| --- | --- | --- | --- |
| orgAuthConfig | object | 可选 | 组织机构认证配置项<br>- 传入机构信息后获取该机构的授权认证或者实名认证页面链接（若不传机构信息，则由用户自主在页面任意填写机构名称、证件号等认证信息）；<br>- 常规场景：组织机构名称orgName与机构账号orgId二者选一项传入即可。 |
| orgAuthConfig.orgName | string | 可选 | 组织机构认证配置项<br>- 传入机构信息后获取该机构的授权认证或者实名认证页面链接（若不传机构信息，则由用户自主在页面任意填写机构名称、证件号等认证信息）；<br>- 常规场景：组织机构名称orgName与机构账号orgId二者选一项传入即可。 |
| orgAuthConfig.orgInfo | object | 可选 | 组织机构身份附加信息 |
| orgAuthConfig.orgInfo.orgIDCardNum | string | 可选 | 组织机构证件号 |
| orgAuthConfig.orgInfo.orgIDCardType | string | 可选 | 组织机构证件类型，可选值如下：<br>- CRED_ORG_USCC - 统一社会信用代码<br>- CRED_ORG_REGCODE - 工商注册号 |
| orgAuthConfig.orgInfo.legalRepName | string | 可选 | 法定代表人姓名 |
| orgAuthConfig.orgInfo.legalRepIDCardNum | string | 可选 | 法定代表人证件号 |
| orgAuthConfig.orgInfo.legalRepIDCardType | string | 可选 | 法定代表人证件类型，可选值如下：<br>- CRED_PSN_CH_IDCARD - 中国大陆居民身份证<br>- CRED_PSN_CH_HONGKONG - 香港来往大陆通行证（回乡证）<br>- CRED_PSN_CH_MACAO - 澳门来往大陆通行证（回乡证）<br>- CRED_PSN_CH_TWCARD - 台湾来往大陆通行证（台胞证）<br>- CRED_PSN_PASSPORT - 护照 |
| orgAuthConfig.orgAuthPageConfig | object | 可选 | 机构实名认证页面配置项 |
| orgAuthConfig.orgAuthPageConfig.orgDefaultAuthMode | string | 可选 | 指定页面中默认选择的机构认证方式：<br>- ORG_BANK_TRANSFER - 对公账户打款认证<br>- ORG_ALIPAY_CREDIT - 法人快捷认证（必须操作人为法定代表人本人场景才会显示，需要法人支付宝刷脸授权完成认证）<br>- ORG_LEGALREP_INVOLVED - 法定代表人认证/法人授权书认证（如操作人为法定代表人本人操作则为法定代表人认证，如非法定代表人本人则为法人授权书认证） |
| orgAuthConfig.orgAuthPageConfig.orgAvailableAuthModes | array<string> | 可选 | 设置页面中可选择的机构认证方式，若不传此参数，则可选择全部认证方式<br>- ORG_BANK_TRANSFER - 对公账户打款认证<br>- ORG_ALIPAY_CREDIT - 法人快捷认证（必须操作人为法定代表人本人场景才会显示，需要法人支付宝刷脸授权完成认证）<br>- ORG_LEGALREP_INVOLVED - 法定代表人认证/法人授权书认证（如操作人为法定代表人本人操作则为法定代表人认证，如非法定代表人本人则为法人授权书认证） |
| orgAuthConfig.orgAuthPageConfig.orgEditableFields | array<string> | 可选 | 设置页面中可编辑的信息，不传此参数，页面默认不允许编辑机构信息。<br>- orgNum - 机构证件号（如果账号已实名，传了该字段，页面也是不可编辑更改的，因为证件号是唯一标识）<br>- legalRepName - 法定代表人姓名<br>- orgBankAccountNum - 企业对公打款银行账户 |
| orgAuthConfig.transactorInfo | object | 可选 | 经办人身份信息<br>（如果不传经办人个人的账号信息，则需要经办人自主在页面填写手机号/邮箱进行验证码回填注册） |
| orgAuthConfig.transactorInfo.psnAccount | string | 可选 | 经办人账号标识（手机号或邮箱） |
| orgAuthConfig.transactorInfo.psnInfo | object | 可选 | 经办人身份信息 |
| orgAuthConfig.transactorInfo.psnInfo.psnName | string | 可选 | 经办人姓名 |
| orgAuthConfig.transactorInfo.psnInfo.psnIDCardNum | string | 可选 | 经办人证件号 |
| orgAuthConfig.transactorInfo.psnInfo.psnIDCardType | string | 可选 | 经办人证件类型，可选值如下：<br>- CRED_PSN_CH_IDCARD - 中国大陆居民身份证<br>- CRED_PSN_CH_HONGKONG - 香港来往大陆通行证（回乡证）<br>- CRED_PSN_CH_MACAO - 澳门来往大陆通行证（回乡证）<br>- CRED_PSN_CH_TWCARD - 台湾来往大陆通行证（台胞证）<br>- CRED_PSN_PASSPORT - 护照<br>【注】CRED_PSN_CH_IDCARD 类型同时兼容港澳台居住证（81、82、83开头18位证件号）、外国人永久居住证（9开头18位证件号） |
| orgAuthConfig.transactorInfo.psnInfo.psnMobile | string | 可选 | 经办人手机号（运营商实名登记手机号或银行卡预留手机号，仅用于认证） |
| authorizeConfig | object | 可选 | 机构授权配置项<br>- 不传此参数默认页面仅实名认证，不需要用户授权；<br>- 实名认证模式下，如用户之前已实名，接口报错："企业用户已实名"；授权认证模式下，如用户之前已实名，正常获取授权链接，需要经办人做个人认证，然后直接授权通过或向企业管理员发起授权审批。 |
| authorizeConfig.authorizedScopes | array<string> | 可选 | 设置页面中权限范围，参数值如下：<br>- 授权当前应用AppId获取用户的账号基本信息：<br>get_org_identity_info - 授权允许获取企业/组织的基本信息（需要授权获取的只有企业的法定代表人证件号，其他信息不授权可直接获取）<br>get_psn_identity_info - 授权允许获取经办人个人用户的账号信息（姓名、手机号/邮箱、证件号等）<br>- 授权当前应用AppId代用户发起合同签署：<br>org_initiate_sign - 授权允许代表企业/组织用户发起合同签署以及查询合同签署详情<br>psn_initiate_sign - 授权允许代表经办人个人用户发起合同签署以及查询合同签署详情<br>- 授权当前应用AppId获取用户资源管理权限：<br>manage_org_member - 授权允许获取企业/组织用户的组织成员的查询、新增、编辑、删除权限<br>manage_org_seal - 授权允许获取企业/组织用户的印章的查询、新增、编辑、授权、删除权限<br>manage_org_template -授权允许获取企业/组织用户的模板的查询、新增、编辑、复制、删除权限<br>use_org_template - 授权允许获取企业/组织用户的模板的使用权限<br>manage_org_resource - 授权允许获取企业/组织用户的印章、组织成员等资源的管理权限（不包含用印权限）<br>manage_psn_resource - 授权允许获取经办人个人用户的印章等资源的管理权限<br>- 授权当前应用AppId存储用户的合同文件：<br>psn_sign_file_storage - 授权允许个人合同文件存储到平台应用的本地服务器<br>org_sign_file_storage - 授权允许企业/组织合同文件存储到平台应用的本地服务器<br>- 授权当前应用AppId获取用户的用印审批信息：<br>org_approval_info - 授权允许获取企业/组织用户的用印审批信息<br>- 授权当前应用AppId获取用户订单使用权限：<br>use_org_order - 授权允许获取企业/组织用户套餐订单的使用权限 |
| redirectConfig | object | 可选 | 认证完成重定向配置项 |
| redirectConfig.redirectUrl | string | 可选 | 认证完成后跳转页面（除app和小程序端集成外，地址需符合 https /http 协议地址）<br>【注】贵司的重定向域名需要在e签宝提前放行，否则会报错：“您即将访问的页面可能有安全风险”。 |
| clientType | string | 可选 | 指定客户端类型，默认值 ALL（注意参数值全部为英文大写）<br>- ALL - 自动适配移动端或PC端<br>- H5 - 移动端适配<br>- PC - PC端适配 |
| notifyUrl | string | 可选 | 接收回调通知的Web地址，通知开发者用户认证和授权的完成以及变更情况 |

## 枚举与约束

- 目前支持对接 e签宝 平台
- `orgAuthConfig.orgInfo.orgIDCardType`（string，可选）：组织机构证件类型，可选值如下：<br>- CRED_ORG_USCC - 统一社会信用代码<br>- CRED_ORG_REGCODE - 工商注册号
- `orgAuthConfig.orgInfo.legalRepIDCardType`（string，可选）：法定代表人证件类型，可选值如下：<br>- CRED_PSN_CH_IDCARD - 中国大陆居民身份证<br>- CRED_PSN_CH_HONGKONG - 香港来往大陆通行证（回乡证）<br>- CRED_PSN_CH_MACAO - 澳门来往大陆通行证（回乡证）<br>- CRED_PSN_CH_TWCARD - 台湾来往大陆通行证（台胞证）<br>- CRED_PSN_PASSPORT - 护照
- `orgAuthConfig.orgAuthPageConfig.orgDefaultAuthMode`（string，可选）：指定页面中默认选择的机构认证方式：<br>- ORG_BANK_TRANSFER - 对公账户打款认证<br>- ORG_ALIPAY_CREDIT - 法人快捷认证（必须操作人为法定代表人本人场景才会显示，需要法人支付宝刷脸授权完成认证）<br>- ORG_LEGALREP_INVOLVED - 法定代表人认证/法人授权书认证（如操作人为法定代表人本人操作则为法定代表人认证，如非法定代表人本人则为法人授权书认证）
- `orgAuthConfig.orgAuthPageConfig.orgAvailableAuthModes`（array<string>，可选）：设置页面中可选择的机构认证方式，若不传此参数，则可选择全部认证方式<br>- ORG_BANK_TRANSFER - 对公账户打款认证<br>- ORG_ALIPAY_CREDIT - 法人快捷认证（必须操作人为法定代表人本人场景才会显示，需要法人支付宝刷脸授权完成认证）<br>- ORG_LEGALREP_INVOLVED - 法定代表人认证/法人授权书认证（如操作人为法定代表人本人操作则为法定代表人认证，如非法定代表人本人则为法人授权书认证）
- `orgAuthConfig.orgAuthPageConfig.orgEditableFields`（array<string>，可选）：设置页面中可编辑的信息，不传此参数，页面默认不允许编辑机构信息。<br>- orgNum - 机构证件号（如果账号已实名，传了该字段，页面也是不可编辑更改的，因为证件号是唯一标识）<br>- legalRepName - 法定代表人姓名<br>- orgBankAccountNum - 企业对公打款银行账户
- `orgAuthConfig.transactorInfo.psnInfo.psnIDCardType`（string，可选）：经办人证件类型，可选值如下：<br>- CRED_PSN_CH_IDCARD - 中国大陆居民身份证<br>- CRED_PSN_CH_HONGKONG - 香港来往大陆通行证（回乡证）<br>- CRED_PSN_CH_MACAO - 澳门来往大陆通行证（回乡证）<br>- CRED_PSN_CH_TWCARD - 台湾来往大陆通行证（台胞证）<br>- CRED_PSN_PASSPORT - 护照<br>【注】CRED_PSN_CH_IDCARD 类型同时兼容港澳台居住证（81、82、83开头18位证件号）、外国人永久居住证（9开头18位证件号）
- `authorizeConfig`（object，可选）：机构授权配置项<br>- 不传此参数默认页面仅实名认证，不需要用户授权；<br>- 实名认证模式下，如用户之前已实名，接口报错："企业用户已实名"；授权认证模式下，如用户之前已实名，正常获取授权链接，需要经办人做个人认证，然后直接授权通过或向企业管理员发起授权审批。
- `authorizeConfig.authorizedScopes`（array<string>，可选）：设置页面中权限范围，参数值如下：<br>- 授权当前应用AppId获取用户的账号基本信息：<br>get_org_identity_info - 授权允许获取企业/组织的基本信息（需要授权获取的只有企业的法定代表人证件号，其他信息不授权可直接获取）<br>get_psn_identity_info - 授权允许获取经办人个人用户的账号信息（姓名、手机号/邮箱、证件号等）<br>- 授权当前应用AppId代用户发起合同签署：<br>org_initiate_sign - 授权允许代表企业/组织用户发起合同签署以及查询合同签署详情<br>psn_initiate_sign - 授权允许代表经办人个人用户发起合同签署以及查询合同签署详情<br>- 授权当前应用AppId获取用户资源管理权限：<br>manage_org_member - 授权允许获取企业/组织用户的组织成员的查询、新增、编辑、删除权限<br>manage_org_seal - 授权允许获取企业/组织用户的印章的查询、新增、编辑、授权、删除权限<br>manage_org_template -授权允许获取企业/组织用户的模板的查询、新增、编辑、复制、删除权限<br>use_org_template - 授权允许获取企业/组织用户的模板的使用权限<br>manage_org_resource - 授权允许获取企业/组织用户的印章、组织成员等资源的管理权限（不包含用印权限）<br>manage_psn_resource - 授权允许获取经办人个人用户的印章等资源的管理权限<br>- 授权当前应用AppId存储用户的合同文件：<br>psn_sign_file_storage - 授权允许个人合同文件存储到平台应用的本地服务器<br>org_sign_file_storage - 授权允许企业/组织合同文件存储到平台应用的本地服务器<br>- 授权当前应用AppId获取用户的用印审批信息：<br>org_approval_info - 授权允许获取企业/组织用户的用印审批信息<br>- 授权当前应用AppId获取用户订单使用权限：<br>use_org_order - 授权允许获取企业/组织用户套餐订单的使用权限
- `clientType`（string，可选）：指定客户端类型，默认值 ALL（注意参数值全部为英文大写）<br>- ALL - 自动适配移动端或PC端<br>- H5 - 移动端适配<br>- PC - PC端适配

## 示例

```bash
contract-cli contract esign org-auth-url --input-file request.json --profile contract --as app
```

官方请求体示例（动态字段接口仍须以当前租户配置为准）：

```json
{
  "orgAuthConfig": {
    "orgName": "******公司",
    "orgInfo": {
      "orgIDCardNum": "9133010****8110212",
      "orgIDCardType": "CRED_ORG_USCC",
      "legalRepName": "这里是法定代表人的姓名",
      "legalRepIDCardNum": "110101********1001",
      "legalRepIDCardType": "CRED_PSN_CH_IDCARD"
    },
    "orgAuthPageConfig": {
      "orgDefaultAuthMode": "ORG_BANK_TRANSFER",
      "orgAvailableAuthModes": [
        "ORG_BANK_TRANSFER",
        "ORG_LEGALREP_INVOLVED"
      ],
      "orgEditableFields": [
        "orgNum"
      ]
    },
    "transactorInfo": {
      "psnAccount": "153****0000",
      "psnInfo": {
        "psnName": "这里是经办人的姓名",
        "psnIDCardNum": "110102*****0000",
        "psnIDCardType": "CRED_PSN_CH_IDCARD",
        "psnMobile": "151****0050"
      }
    }
  },
  "authorizeConfig": {
    "authorizedScopes": [
      "get_org_identity_info",
      "get_psn_identity_info"
    ]
  },
  "redirectConfig": {
    "redirectUrl": "https://www.xxx.cn/"
  },
  "clientType": "ALL",
  "notifyUrl": "http://******/notify"
}
```

## 来源差异说明

- 官方规格路径：`POST /open-apis/esign/auth/orgAuthUrl`
- 本页描述请求参数；响应 envelope 和输出格式遵循共享 Skill。
