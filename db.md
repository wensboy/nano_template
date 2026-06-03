## Nano Template Database Defination

最小化的数据库定义, 以 markdown 表格的形式记录. 对应的表 sql 定义按照实际数据库参照该文档定义即可.

---

### 数据库参数

|表名|字段|值|说明|
|:-:|:-:|:-:|:-:|
|全局|CHARSET|utf8mb4|字符集|
|全局|COLLATE|utf8mb4_general_ci|排序规则|
|全局|ENGINE|InnoDB|存储引擎|

---

### 通用表字段

通用表字段代表所有表在没有特殊说明下一定包含但不显式重复定义的字段, 主要在涉及到: orm 框架, 自定义通用处理 时使用.

*当前项目中包含:*
- gorm

#### gorm

|字段|说明|类型|NULL|DEFAULT|约束|COMMENT|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|id|主键ID|INT UNSIGNED|no|AUTO_INCREMENT|PRIMARY KEY|主键ID|
|created_at|创建时间|DATETIME(3)|yes|NULL||创建时间|
|updated_at|更新时间|DATETIME(3)|yes|NULL||更新时间|
|deleted_at|逻辑删除时间|DATETIME(3)|yes|NULL||逻辑删除时间|

---

### 表定义

> users - 用户表

|字段|说明|类型|NULL|DEFAULT|约束|COMMENT|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|username|用户名|VARCHAR(255)|no|||用户名, 逻辑唯一|
|password|用户密码|VARCHAR(255)|no|||用户密码|

---

> user_profiles - 用户首选项
>
> 1. 用户首选项不会自动关联创建

|字段|说明|类型|NULL|DEFAULT|约束|COMMENT|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|user_id|关联用户ID|INT UNSIGNED|no|||关联: 用户ID|
|avatar|头像地址|VARCHAR(255)|yes|NULL||头像地址|
|nickname|昵称|VARCHAR(255)|yes|NULL||昵称|
|email|邮箱|VARCHAR(255)|yes|NULL||邮箱|
|phone|手机号|VARCHAR(255)|yes|NULL||手机号|
|signature|个性签名|TEXT|yes|NULL||个性签名|

---

> roles - 角色表
>
> 1. 角色表可以直接依赖初始化sql填充数据, 同时支持后续添加
> 2. 角色权限的优先按照: level小优先
> 3. 推荐level和id不要一致, 理由为避免误认

|字段|说明|类型|NULL|DEFAULT|约束|COMMENT|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|name|角色名称|VARCHAR(191)|no|||角色名称|
|level|角色级别|INT UNSIGNED|no|||角色级别|
|description|角色描述|VARCHAR(191)|yes|NULL||角色描述|
|state|状态|INT UNSIGNED|no|||状态 - 0代表未启用, >=1代表其他启用状态|

---

> user_roles - 用户角色关联表
>
> 1. 角色级别锁定用户, 无角色用户禁止进入系统执行后续操作

|字段|说明|类型|NULL|DEFAULT|约束|COMMENT|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|user_id|关联用户ID|INT UNSIGNED|no|||关联用户ID|
|role_id|关联角色ID|INT UNSIGNED|no|||关联角色ID|
|description|关联描述|VARCHAR(191)|yes|NULL||关联描述 -> 备注|
|state|状态|INT UNSIGNED|no|||状态 - 0代表未启用, >=1代表其他启用状态|

---

> ai_sessions - AI会话表

|字段|说明|类型|NULL|DEFAULT|约束|COMMENT|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|user_id|关联用户ID|INT UNSIGNED|no|||关联用户ID|
|session_id|会话标识|BINARY(16)|no|||会话标识|
|title|会话标题|VARCHAR(255)|no|||会话标题|

---

> ai_messages - AI消息表

|字段|说明|类型|NULL|DEFAULT|约束|COMMENT|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|session_id|关联会话ID|INT UNSIGNED|no|||关联会话ID|
|user_message|用户消息|TEXT|no|||用户消息|
|ai_message|AI消息|TEXT|no|||ai消息|
|ai_model|AI模型|VARCHAR(255)|no|||ai模型|

---

> chat_profiles - 聊天首选项表
>
> 1. 不记录和llm供应和连接配置, 只记录llm参数和偏好

|字段|说明|类型|NULL|DEFAULT|约束|COMMENT|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|user_id|关联用户ID|INT UNSIGNED|no|||关联用户ID|
|stream|流式输出|TINYINT(1)|no|||是否开启流式输出|
|think|思考模式|TINYINT(1)|no|||是否开启思考|
|temperature|温度|DECIMAL(10,2)|no|||温度|
|render|渲染结果|TINYINT(1)|no|||是否渲染结果|
|render_rate|渲染速率|INT|no|||渲染速率(打字机渲染速率)|
|render_jit|即时渲染|TINYINT(1)|no|||是否启动即时渲染: 0-关闭, 1-开启|
|max_token|最大token|BIGINT|no|||最大token限制|

---

> chat_appendixs - 聊天附件表

|字段|说明|类型|NULL|DEFAULT|约束|COMMENT|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|session_id|关联会话ID|INT UNSIGNED|no|||关联会话ID|
|mime|MIME类型|VARCHAR(255)|no|||mime类型|
|name|附件名称|VARCHAR(255)|no|||附件名称|
|link|访问链接|VARCHAR(511)|no|||完整访问链接|
|size|附件大小|BIGINT|no|||附件大小|
|source|附件来源|TINYINT(8)|no|||附件源: 1-用户附件, 2-ai附件|

---
