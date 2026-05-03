## frontend index

**核心规则**

- 只修改必要文件.
- 最小化修改, 保证程序的可运行性优先.
- 修改基于可拓展原则.
- 如无必要, 勿增实体.
- 增加优于删除.
- 有必要才加注释.

**代码规则**

- 严格按照 ts 规范编码.
- 目录结构需要严格按照模块化设计.
- 任何代码名称需要高度语义化.
- 必要时做编译检查.

### 目录结构

```md
./
├── eslint.config.js
├── index.html <- 项目入口
├── nginx.conf
├── package.json <- 项目依赖和相关信息定义
├── public
│   ├── favicon.svg
│   └── icons.svg
├── readme.md <- 当前文件
├── src
│   ├── api
│   │   └── interceptor.ts
│   ├── app <- 关于应用相关的配置和定义
│   │   └── store.ts
│   ├── App.tsx
│   ├── assets
│   ├── components
│   │   ├── custom <- 组件目录
│   │   └── native <- 原生组建的位置(只有当需要导入原生组建修改时使用)
│   ├── hooks
│   ├── models <- ts 类型定义
│   ├── mocks <- mock data 定义
│   ├── index.css
│   ├── main.tsx
│   └── pages <- 带路由页面目录
├── tsconfig.json
└── vite.config.ts
```

### 细节

- 在构建过程中如果一个长指令需要经常使用, 记录到 package.json 中作为快速调用. 这对一些测试很常见.
- 所有的组件和页面文件名称为大驼峰命名; hooks/ 下采用符合 react 钩子的小驼峰对应命名.
- 所有的 ts 类型定义在 models/ 中.
- 所有的 mock 数据定义在 mocks/ 中.
- node 版本切换时, 优先使用 fnm, 否则使用 nvm.
- 安装依赖时, 优先使用 pnpm, 否则使用 npm.

### 流程

**当前项目为一个 nano template forntend 部分**

> 必要步骤

- 使用 react-router 引入路由功能. 内部不允许直接使用 window 的 api 直接进行页面路由跳转.
- 定义标准 axios 客户端. 用于接收符合如下 json 格式的响应:

```json
{
    "code": 1, // 操作状态代码: 0 成功, 1 失败, > 1 特殊业务代码, -1 错误
    // 操作信息: 
    // 基于 code, 可以做如下处理:
    // 0: 可以直接在前端中渲染的安全提示信息, 一般为: xxx操作成功
    // 1: 同0
    // > 1: 同0
    // -1: 直接得到服务器原始错误, 统一为: 系统错误, 请稍后重试...
    "message": "",
    // 响应数据: null | object
    // 如果没有任何字段, 后端一定不会返回空 {} 而是 null.
    "data": {}
 }
```

- 如果有接口文档指定, 需要完整定义接口的响应结构到 models/ 中, 单独与其他 model 区分开.

> 后续步骤

- 参照一个页面的描述, 按照: `模型定义 -> 组件拆分实现 -> page 组装 -> 页面路由 -> api 数据获取` 来实现一个页面的完整处理流程.
- 页面完成后才调用一系列的编译检查来检查代码并尝试构建.