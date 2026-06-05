## backend 设计

文档记录后端开发过程中的统一规范, 将 `nano template` 项目后端完善过程中常见的问题总结并优化记录.

---
- [约束](#1)
- [风格](#2)
  - [restful](#2-1)
  - [rpc](#2-2)
  - [database](#2-3)
- [built-in](#3)
  - [common](#3-1)
  - [role](#3-2)
  - [user](#3-3)
  - [chat::ai](#3-4)
- [middleware](#4)
  - [log](#4-1)
  - [parser](#4-2)
---


### <a id="1">约束</a>

使用: 约束应当始终作为核心的规范贯穿在一个完整的任务执行过程当中.

```md
- 单一职责, 保证各模块之间强解耦.
- 最小化实现, 只做必要的处理, 确保可拓展性.
- 如无必要不加实体, 当能够复用时, 抽离组件.
- lsp通过和lint代码检查和编译构建, 不做服务启动.
```

### <a id="2">风格</a>

使用: 风格应当基于已有的代码样例来仿照实现.

#### <a id="2-1">restful</a>

##### 约束
```md
1. uri 应当总是为: /<domain>[/<domain_key>]/<object>[/<object_key][/verb], 在 method 存在 verb 语义时, 可以省略 verb
eg: POST /user/profile -> 创建用户配置 == POST /user/profile[/create]
2. 参数规则如下:
- GET, DELETE, PUT在处理单个资源时, 使用 path query 在 uri 中, 并且如果该标识就是该相关对象的唯一标识, 那么可以省略`domain_` 的前缀.
eg: PUT /user/:user_id/profile/:id == PUT /user/:user_id/profile/:[profile_]id
- GET, DELETE一般没有任何body参数
- 单个处理更新少部分字段时, 使用 PATCH
3. 由于go存在默认零值问题, 参数按照如下解析:
- 网关层 query param 解析为""的一定是未传递的, 该值本质为null
- body 中如果一个参数是可选的, 那么使用 *T
eg: 一般 keywords 是可选的, 在后端结构中定义为: *string
- body 中如果一个参数是零值有效的, 那么使用 *T
eg: 一般 state 是 0 有效的, 在后端结构中定义为: *int
4. 响应规则:
- 对于 data 字段为可选的情况下, data 总是使用中间件的: SuccNoMore 来定义为: {}
- 总是不推荐直接使用 proto 的响应结构来直接作为 data 的值, 即使字段完全一致的情况下, 依然需要按照 gateway 层自己的响应来类型转化. 使用 `Toxxx` 接口或 `omitempty` 绑定参数来处理.
5. 处理器(handler)接口相关命名规则:
- 格式: <verb><object>
eg: CreateProfile, DeleteProfile
```

##### 后端响应规范

```json
{
    "code": 0, // 0 - 成功, 1 - 失败, -1 - 错误
    // code=0: 可以直接使用, 一般为用户友好提示
    // code=1: 可以直接使用, 一般为用户友好提示
    // code=-1: 不能直接使用, 一般为原始错误, 需要前端按照实际业务进行友好提示.
    "message": "...",
    // data 依照后端可以为: null | {} | [], 在当前的标准下只能为: {}, null | [] 不考虑(也不可能出现, 网关层统一处理了)
    "data": {}
}
```

|http状态码|响应中的逻辑代码(code)|
|:---:|:----:|
|200|0|
|200|1|
|其余http状态码|-1|

##### 后端 sse 响应规范


|字段|格式|是否必须|说明|
|:-:|:-:|:-:|:-:|
|data|data: 内容|可选|消息的数据内容，可多行|
|event|event: 事件名|可选|自定义事件类型，默认为 message|
|id|id: 123|可选|消息ID，断线重连时可设置 Last-Event-ID 头|
|retry|	retry: 3000|可选|重连间隔（毫秒），客户端自动生效|

当前项目中的 `data` 序列化规范: 

```json
{
    // 消息索引片段, 总是从 0 开始
    "index": 0,
    "content": "...",
    // 代表 content 内容可以被解析为什么类型, 例如: string, number, boolean, object等
    "type": ""
}
```
最终序列化的结果为: `data: {"index":0,"content":"...","type":"..."}\n\n`

#### <a id="2-2">rpc</a>

##### 约束

```md
1. gateway 只做初步的参数过滤, 实际的参数合法性校验在 rpc 层按照自己的业务实际处理.
2. rpc 的 proto 定义只做最小化定义, 确保完整的可拓展性, 并不推荐不同的 rpc 调用共享: `*Request` 和 `*Response` 的结构定义.
3. rpc 消息字段尽可能使用类似或者一致的类型, 避免过度转换, 除非需要特殊值的判断.
4. rpc 的校验通常使用对应变量的*T形式作为最终容器.
5. rpc 错误消息统一使用 rpc 规范消息.
```

#### <a id="2-3">database</a>

##### 约束

```md
1. 数据库结构映射总是需要显示实现: `TableName() string`, 严格按照指定的 sql 表定义处理.
2. 数据库字段总是包含 gorm 的通用结构.
3. <table_name>_id 总是 ref 到<table_name>.id, 除非指定的 table 存在类型为 `binary(16) - uuid` 的字段.
4. 数据库建模尽可能解耦, 只定义必要的字段和表, 避免过度设计.
5. 对于 update 操作, 通常使用 map 在 service 和 handler 之间传递会比大量参数传递要更好, 对于 query 参数过多的情况也是类似.
6. 涉及到多表关联成功或者失败的操作, 采用事务处理.
7. 数据库映射字段结构理论上针对数据库可空字段应该为 sql 标准包的兼容类型: *Null 类型.
```

### <a id="3">built-in</a>

记录内置服务状态.

#### <a id="3-1">common</a>

```yaml
desc: 通用处理
rpc: true
web: true
check: "passed"
version: "v1.0.0"
last-modified: "2026-06-02 pm"
```

#### <a id="3-2">role</a>

```yaml
desc: 系统角色处理
rpc: true
web: true
check: "passed"
version: "v1.0.0"
last-modified: "2026-06-04 am"
```

#### <a id="3-3">user</a>

```yaml
desc: 用户处理
rpc: true
web: true
check: "passed"
version: "v1.0.0"
last-modified: "2026-06-04 am"
```

#### <a id="3-4">chat::ai</a>

```yaml
desc: ai集成处理
rpc: true
web: true
check: "testing"
version: "v1.0.0"
last-modified: "2026-06-04 am"
```

### <a id="4">middleware</a>

#### <a id="4-1">log</a>

对于一个灵活的通用日志设计思路记录.

```md
日志至少应该可以做到:
- 将输入组织为结构化的可读的信息
- 能够通过灵活模板改变信息结构
- 能够输出到指定的io, 包括但不限于数据库, 消息队列, 日志采集器, 终端, 文件, 网络传输
- 独立存在, 日志器需要创建并且实例是独立的
构建的细节:
- 使用 zap 作为基本日志库加以改造, 最终的调用与其他的日志库无异
- 需要妥善处理 runtime.Caller 问题, 能够灵活控制并总是能够正确在日志时显示调用位置
- 提供一种更加清晰模板集成语法, 能够处理终端日志指定区域文本上色或者做一些更加详细的渲染, 例如下划线, 加粗等.
- 能够支持多种序列化的日志构建, 灵活运用模板定义
```

#### <a id="4-2">parser</a>
