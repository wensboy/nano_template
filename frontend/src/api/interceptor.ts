import axios, {
  AxiosError,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from "axios";

const BASE_URI = "/api/v1"

// 创建 axios 实例
const api = axios.create({
  baseURL: BASE_URI || "/api",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

// 请求拦截器
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 从 sessionStorage 获取 token
    const token = sessionStorage.getItem("token");
    if (token) {
      config.headers.set("Authorization", `Bearer ${token}`);
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  },
);

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  (error: AxiosError) => {
    // 处理 401 错误
    if (error.response?.status === 401) {
      // 删除 sessionStorage 中的 token
      sessionStorage.removeItem("token");
      // todo: 待优化跳转, 先刷新如果没有token, 再跳转到登录页, 不使用window.location.href
      // 跳转到登录页
      // window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);

export default api;
