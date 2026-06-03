import api from "@/api/interceptor";
import {
  trimOptional,
  trimRequired,
  withPagination,
  type Pagination,
  type PaginationParams,
} from "@/api/common";
import type { Request } from "@/api/client";

export type Role = {
  id: number;
  name: string;
  description: string;
  level: number;
  state: number;
  created_at: string;
  updated_at: string;
};

// 创建角色
export type CreateRoleRequest = {
  name: string;
  description?: string;
  level?: number;
  state?: number;
};
export type CreateRoleResponse = {
  id: number;
};
export function CreateRole(req: CreateRoleRequest): Request<CreateRoleResponse> {
  return api.post("/role", {
    name: trimRequired(req.name, "角色名称"),
    description: trimOptional(req.description),
    level: req.level ?? 0,
    state: req.state ?? 0,
  });
}

// 更新角色
export type UpdateRoleRequest = {
  name?: string;
  description?: string;
  level?: number;
  state?: number;
};
export function UpdateRole(id: number, req: UpdateRoleRequest): Request<null> {
  return api.put(`/role/${id}`, {
    name: trimOptional(req.name),
    description: trimOptional(req.description),
    level: req.level ?? 0,
    state: req.state ?? 0,
  });
}

// 获取指定角色
export function GetRole(id: number): Request<Role> {
  return api.get(`/role/${id}`);
}

// 列出角色
export type ListRolesParams = PaginationParams & {
  keywords?: string;
  level?: number;
  state?: number;
};
export function ListRoles(params: ListRolesParams = {}): Request<Pagination<Role>> {
  return api.get("/role", {
    params: {
      ...withPagination(params),
      keywords: trimOptional(params.keywords),
      level: params.level,
      state: params.state,
    },
  });
}

// 删除角色
export function DeleteRole(id: number): Request<null> {
  return api.delete(`/role/${id}`);
}
