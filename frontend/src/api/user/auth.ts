import api from "@/api/interceptor";
import { trimRequired } from "@/api/common";
import type { Request } from "@/api/client";
import type { RegisterFormState } from "@/models/user/user";

// 登录
export type LoginResponse = {
  token: string;
};
export function UserLogin(username: string, password: string, roleId = 0): Request<LoginResponse> {
  return api.post("/user/login", {
    username: trimRequired(username, "用户名"),
    password: trimRequired(password, "密码"),
    role_id: roleId,
  });
}

// 登出
export function UserLogout(): Request<null> {
  return api.get("/user/logout");
}

// 注册
export function UserRegister(form: RegisterFormState, roleId = 0): Request<null> {
  const normalizedUsername = trimRequired(form.username, "用户名");
  const normalizedPassword = trimRequired(form.password, "密码");
  const normalizedConfirmPassword = trimRequired(form.confirmPassword, "确认密码");

  if (normalizedPassword !== normalizedConfirmPassword) {
    throw new Error("两次输入的密码不一致");
  }

  return api.post("/user/register", {
    username: normalizedUsername,
    password: normalizedPassword,
    role_id: roleId,
  });
}

// 登录校验
export type AuthResponse = {
  user_id: number;
  username: string;
  role_id: number;
  level: number;
};
export function AuthUser(): Request<AuthResponse> {
  return api.get("/user/auth");
}
