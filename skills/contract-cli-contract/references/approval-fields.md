# Approval Fields Reference

本页适用于：

- `contract-cli contract approval start <process-instance-id> --input-file approval.json`

字段口径按 CLM 后端 `TaskApprovalDTO`、`ProcessOpenPlatformController`、`ProcessOpenPlatformService` 和 `ProcessCommandType` 校验整理；JSON 字段使用 snake_case。

## 路由

`contract approval start` 路由：`POST /open-apis/contract/v1/process_instances/{process_instance_id}/task_approval`

## 顶层字段

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `assignee_id` | integer | 是 | 审批人飞书 larkId；后端优先按 larkId 查员工 |
| `assignee_employee_id` | integer | 否 | 员工 ID；只有 `assignee_id` 查不到时才兜底使用，但入口仍要求 `assignee_id` 非空 |
| `command_type` | string | 是 | 审批命令名称 |
| `task_instance_id` | string | 见说明 | 非作废命令必须传；作废命令不传时后端会尝试取当前待办第一个任务 |
| `task_comment` | string | 见说明 | 拒绝命令必填 |
| `archive_number` | string | 否 | 归档节点同意时可传；只允许中文、大写英文、数字和 `-_*/【】` |
| `reject_info` | object | 否 | 拒绝配置 |
| `submit_time` | string | 否 | 定制化参数；普通 `task_approval` 入口不会写入 `IdHolder`，`bespoke` 入口才使用 |

## command_type

当前公开入口实际允许这些命令：

| command_type | 后端枚举 | 说明 |
| --- | --- | --- |
| `general` | `AGREE` | 同意 |
| `terminationProcess` | `CANCEL` | 作废整个流程；仅提交节点支持 |
| `rollBackTaskByExpression` | `REJECT_AND_TERMINATE` | 拒绝；`task_comment` 必填 |
| `recover` | `REVOKE` | 撤回到提交节点 |
| `submit` | `RESUBMIT` | 重新提交；仅提交节点支持 |

其他 `ProcessCommandType` 虽然存在于枚举中，但当前 `taskApproval` switch 会返回“不支持的命令类型”。

## reject_info

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `reject_return_code` | integer | 否 | `0` 重新流转，`1` 回到当前节点；只在节点配置允许审批人选择时生效 |
| `no_cooperation` | boolean | 否 | 拒绝后是否不允许协商；在对应节点配置下生效，默认 `true` |

## 示例

同意：

```json
{
  "assignee_id": 1152307041667121520,
  "task_instance_id": "task-instance-001",
  "command_type": "general"
}
```

拒绝：

```json
{
  "assignee_id": 1152307041667121520,
  "task_instance_id": "task-instance-001",
  "command_type": "rollBackTaskByExpression",
  "task_comment": "补充付款材料后重新提交",
  "reject_info": {
    "reject_return_code": 0,
    "no_cooperation": true
  }
}
```

作废：

```json
{
  "assignee_id": 1152307041667121520,
  "command_type": "terminationProcess"
}
```
