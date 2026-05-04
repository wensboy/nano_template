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
  withCredentials: true, // 开启 cookie 处理
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
      // todo: 如果 sessionStorage中有 token, 直接删除. 如果没有, 尝试发起登出请求通知后端将 cookie 中的token相关的字段清除
      // 定向到登录页面. 
      sessionStorage.removeItem("token");
    }
    return Promise.reject(error);
  },
);

export default api;
