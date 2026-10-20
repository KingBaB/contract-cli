# 合同交接 Skill 接口白名单清单

来源：`skills/contract-handover/references/api-contract.yaml` 的 recipes；按 HTTP 方法与 Open Web 路径去重，共 9 个端点。前后回查复用详情接口，无需重复放行。

| 方法 | Open Web 白名单路径 | 后端原始路径 | 用途 |
|---|---|---|---|
| GET | `/open-web/web/clm/employee/info` | `/clm/api/employee/info` | 读取当前运行身份的内部 employeeId/tenantId，避免混用开放平台 ID |
| GET | `/open-web/web/clm/permission/function/list` | `/clm/api/permission/function/list` | 读取当前用户功能权限码，不构造管理员或代办身份 |
| GET | `/open-web/web/clm/employee/findEmployeesByKeywords` | `/clm/api/employee/findEmployeesByKeywords` | 按姓名查找合同接收人候选，展示消歧后选择一个当前内部员工 ID |
| GET | `/open-web/web/clm/employee/getEmployeesByIds` | `/clm/api/employee/getEmployeesByIds` | 重新读取唯一接收人的当前租户、状态、部门；确认部门实际取第一项 |
| POST | `/open-web/web/clm/es/search/contractList` | `/clm/api/es/search/contractList` | 复现全部合同/我的合同页搜索；生成待确认固定入口清单 |
| GET | `/open-web/web/clm/contract/info` | `/clm/api/contract/info` | 读取普通处理器资源的当前角色、需求人、部门、状态和关联标识；已确认组成员由父组成员接口分流 |
| GET | `/open-web/web/clm/contractAssembly/info` | `/clm/api/contractAssembly/info` | 读取合同组/模板多合同本体，取得有效子合同总数及组专用回查字段 |
| GET | `/open-web/web/clm/contractAssembly/childrenInfo` | `/clm/api/contractAssembly/childrenInfo` | 按父组权限完整分页读取有效成员的身份、状态、当前角色和需求人；不要求成员部门读取或核验 |
| POST | `/open-web/web/clm/permission/handover` | `/clm/api/permission/handover` | 对用户确认的固定入口资源执行合同交接，每批至多100个入口且同影响组不拆分 |

普通指定合同核心流程使用员工身份、权限、接收人候选与详情、普通合同详情和交接接口；按编号或条件定位时增加搜索接口。合同组场景增加父组详情、成员分页两个接口。

唯一业务写入端点为 `POST /open-web/web/clm/permission/handover`。搜索虽然使用 POST，但契约定义为读取。白名单应按精确方法与路径放行，不据此放开整个 `/open-web/web/**`。

此清单是 Skill 的业务接口需求，不代表网关白名单已配置。全部调用复用 user 身份；鉴权、租户隔离、资源权限仍由网关及后端校验。登录授权接口属于 CLI 自有流程，未在本 Skill 配方中新增。

当前契约仍为 test 开发态。通用接口合并和取消 CLI 版本记录，不等于解除环境限制或完成 blue 业务验收。
