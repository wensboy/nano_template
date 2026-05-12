import api from "@/api/interceptor";
import type { Request } from "@/api/client";

export type PresignRequest = {
  object_key: string;
  mime: string;
  size?: number;
  sender?: string;
};

export type SignedHeader = {
  key: string;
  value: string;
};

export type PresignResponse = {
  signed_url: string;
  method: string;
  expiration: string;
  signed_headers: SignedHeader[];
};

/**
 * Request a presigned URL for uploading an object to Aliyun OSS.
 * The returned signed_url should be used with an HTTP PUT to upload the file.
 */
export function PresignOssObject(
  req: PresignRequest,
): Request<PresignResponse> {
  if (!req.object_key.trim()) {
    throw new Error("object_key 不能为空");
  }
  if (!req.mime.trim()) {
    throw new Error("mime 不能为空");
  }

  return api.post("/native/aliyun/presign", {
    object_key: req.object_key,
    mime: req.mime,
    size: req.size ?? 0,
    sender: req.sender ?? "",
  });
}
