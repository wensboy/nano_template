## Store

关于 frontend(针对 react 框架而言) 中的状态管理技术: 上下文提供器, 全局状态管理, 组件内部状态.

---

- [上下文提供器](#context-provider)
- [全局状态管理](#global-state)
- [组件内部状态](#local-state)

---

### <a id="context-provider">上下文提供器</a>

通用的上下文提供器模板:

```js
import { useContext, createContext, useState, useMemo } from "react";

const AppContext = createContext(undefined);

function AppProvider({ children }) {
  const [user, setUser] = useState(null);

  const value = useMemo(()=>{
    user: user,
    setUser: setUser,
  },[user]);

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}

function useAppContext() {
    const context = useContext(AppContext);

    if (context===undefined) {
        throw new Error('useAppContext must be used within an AppProvider');
    }

    return context;
}

export default useAppContext;
```

### <a id="global-state">全局状态管理</a>

一些常见的全局状态管理库: `redux, zustand, pinia`等.

> Redux

redux 的核心原则为: `state -> view -> action`.

一般包含如下几个组件:

**action**

```js
const action = {
    type: "domain/eventName",
    payload: {},
    meta: ...,
    error: ...
};
// 通常会封装为: action creater
const actionCreater = () => {
    return {
        type: "domain/eventName",
        payload: {},
        meta: ...,
        error: ...
    }
}
```

**reducer**

```js
// reducer本质是一个状态计算纯函数
const reducer = (state, action) => {
  // action.type在js中直接依赖 dispatch 来做校验
  switch (action.type) {
  }
};
```

**store**

```js
import { configureStore } from "@reduxjs/toolkit";
const store = configureStore({
    reducer: ...
})
// store 的核心方法:
// 1. dispatch
// 2. getState
export default store;
```

**selecter**

```js
// 一般用作需要拿到重复相同的数据集合的方法
const selecter = (state) => {
  // 返回 state 中指定数据构建的集合
};
```

### <a id="local-state">组件内部状态</a>
