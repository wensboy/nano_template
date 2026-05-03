## backend index

### 核心规则

- 只修改必要文件.
- 最小化修改, 保证程序的可运行性优先.
- 修改基于可拓展原则.
- 如无必要, 勿增实体.
- 增加优于删除.
- 有必要才加注释.

### 项目结构

**基目录: `backend/`**

```md
tips:

1. 没有出现在下方目录中的但是实际存在的文件为**非重要文件**, 可以忽略.
2. 配置在修改时**只需要修改实际程序读取**的配置.

---

│ .env <- 实际程序读取
│ .env.example <- 配置模板
│ build.sh <- 复用指令记录位置
│ dev.sql
│ env.example.yaml <- 配置模板
│ env.yaml <- 实际程序读取
│ go.mod
│ go.sum
│ readme.md <- 当前文件
├───cmd
│ │ main.go <- 程序入口
│ │ server.go
│ └───docs <- swagger api自动生成目录
│ docs.go
│ swagger.json
│ swagger.yamls
├───data <- 数据库目录
├───pkg <- 按照语义搜索
│ ├───config
│ │ config.go
│ │ database.go
│ │ flag.go
│ │ jwt.go
│ │ llm.go
│ │ proxy.go
│ │ server.go
│ │ template.go
│ ├───middleware
│ │ auth.go
│ │ error.go
│ │ llm.go
│ │ proxy.go
│ │ response.go
│ ├───services
│ │ ├───common
│ │ │ handler.go
│ │ │ router.go
│ │ ├───native <- 第三方 api 接入
│ │ └───user_sys
│ │ handler.go
│ │ model.go
│ │ router.go
│ │ service.go
│ └───util
│ io.go
│ log.go
│ pool.go
└───templates <- llm prompt template
system_default.md
user_default.md
```

### 细节

1. 数据处理

```md
router -> handler -> service
```

- router 按照实际需求拆分为 public router 和 private router, 必要时分别注册中间件
- handler 分为大驼峰的 interface 和小驼峰的 struct, handler 接口不能直接传递 cfg, 需要依赖上下文和中间件传递. handler 总是需要编写 swagger 注解.
- service 中使用 gorm 整合数据库数据.

2. 配置优先级

```md
配置加载流程:

1. .env -> env
2. env -> env.yaml
3. env.yaml -> cli

优先级: cli > env(.env) > env.yaml
```

3. 数据库原则

数据存储目录: `data/`

- 使用 c/s 架构的数据库时, 无须处理数据库存储问题.
- 使用 fs 的数据库时, 将 .db 文件定义到存储目录.

### 流程

当用户给出一个功能描述时(例如: 评论, 帖子等), 按照如下过程构建:

- 构建目录和文件: `{domain}_sys/`, `{domain}_sys/model.go`, `{domain}_sys/service.go`, `{domain}_sys/handler.go`, `{domain}_sys/router.go`.
- 实现顺序: 完成 model.go, 定义相关的数据类型 -> 完成 service.go 和 handler.go 中的业务逻辑交互处理 -> 完成 router.go 的路由设置 -> server.go 中挂载路由部分挂载指定路由.
- 检查流程和测试编译.

**参考: pkg/services/user_sys/**