# contract esign personal-auth-url Parameters

本页专用于 `contract-cli contract esign personal-auth-url`。

- 接口：`POST /open-apis/esign/auth/psnAuthUrl`
- 身份：仅 `app`
- 请求体：JSON 必填，`--input-file` 与 `--data` 二选一
- 官方 OpenAPI：[获取个人认证&授权页面链接](https://docs.qfei.cn/367691966e0.md)
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
| psnAuthConfig | object | 必填 | 个人实名认证配置项<br>- 传入个人账号信息后获取对应个人用户的授权认证或者实名认证页面链接；<br>- 如果不传个人的账号信息（psnAccount/psnId），则需要个人用户自主在页面填写手机号/邮箱进行验证码回填注册； |
| psnAuthConfig.psnAccount | string | 可选 | 个人用户账号标识（手机号或邮箱）<br>【注】psnAccount与psnId二选一传值即可，未知用户psnId时，直接传此字段 |
| psnAuthConfig.psnInfo | object | 可选 | 个人身份附加信息 |
| psnAuthConfig.psnInfo.psnName | string | 可选 | 姓名 |
| psnAuthConfig.psnInfo.psnIDCardNum | string | 可选 | 证件号码 |
| psnAuthConfig.psnInfo.psnIDCardType | string | 可选 | 证件类型，可选值如下：<br>CRED_PSN_CH_IDCARD - 中国大陆居民身份证<br>CRED_PSN_CH_HONGKONG - 香港来往大陆通行证（回乡证）<br>CRED_PSN_CH_MACAO - 澳门来往大陆通行证（回乡证）<br>CRED_PSN_CH_TWCARD - 台湾来往大陆通行证（台胞证）<br>CRED_PSN_PASSPORT - 护照<br>【注】CRED_PSN_CH_IDCARD 类型同时兼容港澳台居住证（81、82、83开头18位证件号）、外国人永久居住证（9开头18位证件号） |
| psnAuthConfig.psnAuthPageConfig | object | 可选 | 个人实名认证页面配置项 |
| psnAuthConfig.psnAuthPageConfig.psnDefaultAuthMode | string | 可选 | 设置页面中默认选择的实名认证方式，可选值如下：<br>PSN_FACE - 人脸识别认证（默认值）<br>PSN_MOBILE3 - 手机运营商三要素认证<br>PSN_BANKCARD4 - 银行卡四要素认证<br>【注】使用iframe内嵌集成不支持对接刷脸方式 |
| psnAuthConfig.psnAuthPageConfig.psnAvailableAuthModes | array<string> | 可选 | 设置页面中可选择的个人认证方式范围，若不传此参数，则可选择全部认证方式。<br>- PSN_FACE - 人脸识别认证<br>- PSN_MOBILE3 - 手机运营商三要素认证<br>- PSN_BANKCARD4 - 银行卡四要素认证<br>【注】使用iframe内嵌集成不支持对接刷脸方式 |
| authorizeConfig | object | 可选 | 个人授权配置项<br>- 不传此参数默认页面仅实名认证，不需要用户授权；<br>- 实名认证模式下，如用户之前已实名，接口报错："个人用户已实名"；授权认证模式下，如用户之前已实名，正常获取授权链接，需要用户获取验证码进入页面后，直接授权成功。 |
| authorizeConfig.authorizedScopes | array<string> | 可选 | 设置页面中权限范围，参数值如下：<br>- 授权当前应用AppId获取用户的账号基本信息：<br>get_psn_identity_info - 授权允许获取个人用户的账号信息（姓名、手机号/邮箱、证件号等）<br>- 授权当前应用AppId代用户发起合同签署：<br>psn_initiate_sign - 授权允许代表个人用户发起合同签署以及查询合同签署详情<br>- 授权当前应用AppId获取用户资源管理权限：<br>manage_psn_resource - 授权允许获取个人用户的印章等资源的管理权限<br>- 授权当前应用AppId存储用户的合同文件：<br>psn_sign_file_storage - 授权允许个人合同文件存储到平台应用的本地服务器 |
| notifyUrl | string | 可选 | 接收回调通知的Web地址，通知开发者用户认证和授权的完成以及变更情况 |
| clientType | string | 可选 | 指定客户端类型，默认值 ALL（注意参数值全部为英文大写）<br>ALL - 自动适配移动端或PC端<br>H5 - 移动端适配<br>PC - PC端适配 |
| redirectConfig | object | 可选 | 认证完成重定向配置项 |
| redirectConfig.redirectUrl | string | 可选 | 认证完成后跳转页面（除app和小程序端集成外，地址需符合 https /http 协议地址）<br>【注】贵司的重定向域名需要在e签宝提前放行，否则会报错：“您即将访问的页面可能有安全风险”。 |

## 枚举与约束

- 目前支持对接 e签宝 平台
- `psnAuthConfig.psnInfo.psnIDCardType`（string，可选）：证件类型，可选值如下：<br>CRED_PSN_CH_IDCARD - 中国大陆居民身份证<br>CRED_PSN_CH_HONGKONG - 香港来往大陆通行证（回乡证）<br>CRED_PSN_CH_MACAO - 澳门来往大陆通行证（回乡证）<br>CRED_PSN_CH_TWCARD - 台湾来往大陆通行证（台胞证）<br>CRED_PSN_PASSPORT - 护照<br>【注】CRED_PSN_CH_IDCARD 类型同时兼容港澳台居住证（81、82、83开头18位证件号）、外国人永久居住证（9开头18位证件号）
- `psnAuthConfig.psnAuthPageConfig.psnDefaultAuthMode`（string，可选）：设置页面中默认选择的实名认证方式，可选值如下：<br>PSN_FACE - 人脸识别认证（默认值）<br>PSN_MOBILE3 - 手机运营商三要素认证<br>PSN_BANKCARD4 - 银行卡四要素认证<br>【注】使用iframe内嵌集成不支持对接刷脸方式
- `psnAuthConfig.psnAuthPageConfig.psnAvailableAuthModes`（array<string>，可选）：设置页面中可选择的个人认证方式范围，若不传此参数，则可选择全部认证方式。<br>- PSN_FACE - 人脸识别认证<br>- PSN_MOBILE3 - 手机运营商三要素认证<br>- PSN_BANKCARD4 - 银行卡四要素认证<br>【注】使用iframe内嵌集成不支持对接刷脸方式
- `authorizeConfig`（object，可选）：个人授权配置项<br>- 不传此参数默认页面仅实名认证，不需要用户授权；<br>- 实名认证模式下，如用户之前已实名，接口报错："个人用户已实名"；授权认证模式下，如用户之前已实名，正常获取授权链接，需要用户获取验证码进入页面后，直接授权成功。
- `authorizeConfig.authorizedScopes`（array<string>，可选）：设置页面中权限范围，参数值如下：<br>- 授权当前应用AppId获取用户的账号基本信息：<br>get_psn_identity_info - 授权允许获取个人用户的账号信息（姓名、手机号/邮箱、证件号等）<br>- 授权当前应用AppId代用户发起合同签署：<br>psn_initiate_sign - 授权允许代表个人用户发起合同签署以及查询合同签署详情<br>- 授权当前应用AppId获取用户资源管理权限：<br>manage_psn_resource - 授权允许获取个人用户的印章等资源的管理权限<br>- 授权当前应用AppId存储用户的合同文件：<br>psn_sign_file_storage - 授权允许个人合同文件存储到平台应用的本地服务器
- `clientType`（string，可选）：指定客户端类型，默认值 ALL（注意参数值全部为英文大写）<br>ALL - 自动适配移动端或PC端<br>H5 - 移动端适配<br>PC - PC端适配

## 示例

```bash
contract-cli contract esign personal-auth-url --input-file request.json --profile contract --as app
```

官方请求体示例（动态字段接口仍须以当前租户配置为准）：

```json
{
  "psnAuthConfig": {
    "psnAccount": "183****0101",
    "psnInfo": {
      "psnName": "赵四",
      "psnIDCardNum": "130204********1001",
      "psnIDCardType": "CRED_PSN_CH_IDCARD"
    },
    "psnAuthPageConfig": {
      "psnDefaultAuthMode": "PSN_MOBILE3",
      "psnAvailableAuthModes": [
        "PSN_BANKCARD4",
        "PSN_MOBILE3",
        "PSN_FACE"
      ]
    }
  },
  "authorizeConfig": {
    "authorizedScopes": [
      "get_psn_identity_info"
    ]
  },
  "notifyUrl": "http://xx.xx.xx.172:8081/CSTNotify/asyn/notify",
  "clientType": "ALL",
  "redirectConfig": {
    "redirectUrl": "https://www.xxx.cn/"
  }
}
```

## 来源差异说明

- 官方规格路径：`POST /open-apis/esign/auth/psnAuthUrl`
- 本页描述请求参数；响应 envelope 和输出格式遵循共享 Skill。
