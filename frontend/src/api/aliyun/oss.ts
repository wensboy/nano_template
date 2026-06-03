import api from "@/api/interceptor";
import {
  trimOptional,
  trimRequired,
  withPagination,
  type PaginationParams,
} from "@/api/common";
import type { Request } from "@/api/client";

export type PresignUploadRequest = {
  object_key: string;
  mime?: string;
  size?: number;
  sender?: string;
};
export type PresignRequest = PresignUploadRequest;

export type PresignDownloadRequest = {
  bucket_prefix?: string;
  object_key: string;
  getter?: string;
};

export type PresignListRequest = {
  getter?: string;
};

export type PresignListParams = PaginationParams & {
  bucket_prefix?: string;
  marker?: string;
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

export type ObjectOwner = {
  id: string;
  display_name: string;
};

export type ObjectProperties = {
  key: string;
  type: string;
  size: number;
  etag: string;
  last_modified: string;
  storage_class: string;
  owner?: ObjectOwner;
  restore_info?: string;
  transition_time?: string;
};

export type CommonPrefix = {
  prefix: string;
};

export type ListObjectsResponse = {
  name: string;
  prefix: string;
  marker: string;
  max_keys: number;
  delimiter: string;
  is_truncated: boolean;
  next_marker: string;
  encoding_type: string;
  contents: ObjectProperties[];
  common_prefixes: CommonPrefix[];
};

// 获取上传预签名地址
export function PresignOssUpload(req: PresignUploadRequest): Request<PresignResponse> {
  return api.post("/native/aliyun/presign/upload", {
    object_key: trimRequired(req.object_key, "object_key"),
    mime: trimOptional(req.mime),
    size: req.size ?? 0,
    sender: req.sender ?? "",
  });
}

// 获取下载预签名地址
export function PresignOssDownload(req: PresignDownloadRequest): Request<PresignResponse> {
  return api.post("/native/aliyun/presign/download", {
    bucket_prefix: trimOptional(req.bucket_prefix),
    object_key: trimRequired(req.object_key, "object_key"),
    getter: trimOptional(req.getter),
  });
}

// 列出 OSS 对象
export function PresignOssList(
  req: PresignListRequest = {},
  params: PresignListParams = {},
): Request<ListObjectsResponse> {
  return api.post(
    "/native/aliyun/presign/list",
    {
      getter: trimOptional(req.getter),
    },
    {
      params: {
        ...withPagination(params),
        bucket_prefix: trimOptional(params.bucket_prefix),
        marker: trimOptional(params.marker),
      },
    },
  );
}

// 获取上传预签名地址
export function PresignOssObject(req: PresignRequest): Request<PresignResponse> {
  return PresignOssUpload(req);
}
