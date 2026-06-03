import api from "@/api/interceptor";
import {
  trimOptional,
  trimRequired,
  withPagination,
  type Pagination,
  type PaginationParams,
} from "@/api/common";
import type { Request } from "@/api/client";

// 更新用户密码
export type UpdatePasswordRequest = {
  old_password: string;
  new_password: string;
};
export function UpdatePassword(req: UpdatePasswordRequest): Request<null> {
  return api.put("/user/password", {
    old_password: trimRequired(req.old_password, "旧密码"),
    new_password: trimRequired(req.new_password, "新密码"),
  });
}

// 注销用户
export function CancelUser(): Request<null> {
  return api.delete("/user");
}

// 创建用户首选项
export type CreateUserProfileRequest = {
  avatar?: string;
  nickname?: string;
  email?: string;
  phone?: string;
  signature?: string;
};
export type UserProfile = {
  id: number;
  user_id: number;
  avatar: string;
  nickname: string;
  email: string;
  phone: string;
  signature: string;
  created_at: string;
  updated_at: string;
};
function toProfilePayload(req: CreateUserProfileRequest) {
  return {
    avatar: trimOptional(req.avatar),
    nickname: trimOptional(req.nickname),
    email: trimOptional(req.email),
    phone: trimOptional(req.phone),
    signature: trimOptional(req.signature),
  };
}
export function CreateUserProfile(req: CreateUserProfileRequest): Request<null> {
  return api.post("/user/profile", toProfilePayload(req));
}

// 获取用户首选项
export function GetUserProfile(): Request<UserProfile> {
  return api.get("/user/profile");
}

// 列出用户首选项
export type ListUserProfilesParams = PaginationParams & {
  keywords?: string;
};
export function ListUserProfiles(params: ListUserProfilesParams = {}): Request<Pagination<UserProfile>> {
  return api.get("/user/profiles", {
    params: {
      ...withPagination(params),
      keywords: trimOptional(params.keywords),
    },
  });
}

// 更新用户首选项
export type UpdateUserProfileRequest = CreateUserProfileRequest;
export function UpdateUserProfile(req: UpdateUserProfileRequest): Request<null> {
  return api.put("/user/profile", toProfilePayload(req));
}

// 删除用户首选项
export function DeleteUserProfile(): Request<null> {
  return api.delete("/user/profile");
}
