import api from "@/api/interceptor";
import { Request } from "@/api/client";
import { RegisterFormState } from "@/models/user_sys/user";

function assertNonEmptyValue(value: string, fieldLabel: string) {
  if (!value.trim()) {
    throw new Error(`${fieldLabel}不能为空`);
  }
}

// 登录请求端点
type LoginResponse = {
  token: string;
};
export function UserLogin(username: string, password: string): Request<LoginResponse> {
  const normalizedUsername = username.trim();
  const normalizedPassword = password.trim();

  assertNonEmptyValue(normalizedUsername, "用户名");
  assertNonEmptyValue(normalizedPassword, "密码");

  return api.post("/user/login", {
    username: normalizedUsername,
    password: normalizedPassword,
  });
}

// 注册请求端点
type RegisterResponse = {
  user_id: number;
};
export function UserRegister(form: RegisterFormState): Request<RegisterResponse> {
  const normalizedUsername = form.username.trim();
  const normalizedPassword = form.password.trim();
  const normalizedConfirmPassword = form.confirmPassword.trim();

  assertNonEmptyValue(normalizedUsername, "用户名");
  assertNonEmptyValue(normalizedPassword, "密码");
  assertNonEmptyValue(normalizedConfirmPassword, "确认密码");

  if (normalizedPassword !== normalizedConfirmPassword) {
    throw new Error("两次输入的密码不一致");
  }

  return api.post("/user/register", {
    username: normalizedUsername,
    password: normalizedPassword,
  });
}
