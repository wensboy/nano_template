import axios, { AxiosError, AxiosResponse } from "axios";

import { store } from "@/app/store";
import { clearUser } from "@/app/store/userSlice";

const BASE_URI = "/api/v1";
export const UNAUTHORIZED_EVENT = "auth:unauthorized";

function handleUnauthorized() {
  store.dispatch(clearUser());
  window.dispatchEvent(new CustomEvent(UNAUTHORIZED_EVENT));
}

// 创建 axios 实例
const api = axios.create({
  baseURL: BASE_URI || "/api/v0",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true, // 开启 cookie 处理
});

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      handleUnauthorized();
    }
    return Promise.reject(error);
  },
);

export default api;
