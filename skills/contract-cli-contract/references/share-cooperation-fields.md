# Share And Cooperation Fields

这份文档用于 `contract share get`、`contract cooperation link get` 和 `contract cooperation record get`。

官方来源：

- [查询合同分享记录](https://docs.qfei.cn/430216002e0)
- [合同协商操作记录信息查询](https://docs.qfei.cn/367662486e0)

## 1. CLI 命令面

```bash
contract-cli contract share get <contract-id> --profile contract --as app
contract-cli contract cooperation link get <contract-id> --profile contract --as app
contract-cli contract cooperation record get <contract-id> --profile contract --as app
```

硬约束：

- 当前三条命令都仅支持 `--as app`
- 三条命令都以 `<contract-id>` 作为 path 参数
- `--user-id-type` / `--user-id` 按共享规则透传
- 分享或协商记录为空时，只要响应 `code=0`，命令仍可视为可用

## 2. `contract share get`

接口：`GET /open-apis/contract/v1/contracts/{contract_id}/share_records`

读取路径：`data[]`

| 字段 | 类型 | 业务含义 |
| --- | --- | --- |
| `business_id` | `string` | 业务 id。 |
| `business_type` | `string` | 业务类型，如 `CONTRACT_APPLICATION`。 |
| `from_employee_id` | `string` | 发起分享员工 id。 |
| `from_employee_name` | `string` | 发起分享员工名称。 |
| `to_employee_id` | `string` | 接收分享员工 id。 |
| `to_employee_name` | `string` | 接收分享员工名称。 |
| `start_time` | `string` | 开始时间，格式 `yyyy-MM-dd`。 |
| `has_expiration` | `boolean` | 是否有过期时间配置。 |
| `end_time` | `string` | 结束时间，格式 `yyyy-MM-dd`。 |
| `enable` | `boolean` | 是否有效。 |

## 3. `contract cooperation link get`

接口：`GET /open-apis/contract/v1/contracts/{contract_id}/cooperation_link`

当前本地命令只透传响应。使用时优先查看：

- `data.cooperation_link`
- `data.contract_id`
- `data.contract_number`
- `data.contract_name`

若后端返回字段名与上述路径不同，以 `--raw` 输出为准。

## 4. `contract cooperation record get`

接口：`GET /open-apis/contract/v1/contracts/{contract_id}/cooperation_record_info`

读取路径：`data`

| 字段 | 类型 | 业务含义 |
| --- | --- | --- |
| `contract_id` | `string` | 合同 id。 |
| `contract_name` | `string` | 合同名称。 |
| `contract_number` | `string` | 合同编号。 |
| `cooperation_id` | `string` | 协商 id。 |
| `cooperation_record_infos` | `array<object>` | 协商操作记录列表。 |

`cooperation_record_infos[]` 常用字段：

| 字段 | 类型 | 业务含义 |
| --- | --- | --- |
| `cooperation_record_type_code` | `integer` | 记录类型编码。 |
| `cooperation_record_type_name` | `string` | 记录类型名称。 |
| `cost_time` | `number|null` | 用时，单位 ms。确认合同、定稿等动作有统计。 |
| `message` | `string|null` | 留言。 |
| `operate_info` | `object` | 操作内容。 |
| `operator_user_info` | `object` | 操作人信息。 |
| `operate_time` | `string` | 操作时间。 |

常见记录类型：

- `1`：发起协商
- `11`：模板修改文本后发起协商
- `12`：审批拒绝后发起协商
- `13`：申请撤回后发起协商
- `14`：审批人发起协商
- `15`：再次协商
- `2`：加入协商者
- `21`：退出协商
- `22`：转办协商
- `23`：移除协商
- `24`：协商参与者查看合同
- `25`：确认合同
- `26`：申请定稿
- `27`：定稿
- `3`：添加附件
- `31`：更新文件
- `32`：文件下载
- `8`：完成协商
- `9`：取消协商
