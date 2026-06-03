import api from "@/api/interceptor";
import { trimRequired } from "@/api/common";
import type { Request } from "@/api/client";

// ping gateway
export function Ping(): Request<null> {
  return api.get("/ping");
}

// 拿到服务信息
export type InspectResponse = {
  version: string;
  author: string;
  description: string;
};
export function Inspect(): Request<InspectResponse> {
  return api.get("/inspect");
}

// 拿到指定模板
export type TemplateFrontMatter = {
  role: string;
};
export type GetTemplateResponse = {
  id: string;
  frontmatter: TemplateFrontMatter;
  content: string;
};
export function GetTemplate(templateId: string): Request<GetTemplateResponse> {
  return api.get(`/template/${encodeURIComponent(trimRequired(templateId, "模板 ID"))}`);
}
