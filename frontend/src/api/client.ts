import type { AxiosResponse } from "axios";

import type { ApiResponse } from "@/models/response";

type ApiResponseLike<TData> = ApiResponse<TData> | AxiosResponse<ApiResponse<TData>>;

export type RequestFactory<TData> = () => Promise<ApiResponseLike<TData>>;

export type ResponseInterceptor<TData> = (
  response: ApiResponse<TData>,
) => ApiResponse<TData> | Promise<ApiResponse<TData>>;

function isAxiosResponse<TData>(
  response: ApiResponseLike<TData>,
): response is AxiosResponse<ApiResponse<TData>> {
  return typeof response === "object" && response !== null && "data" in response && "status" in response;
}

function normalizeResponse<TData>(response: ApiResponseLike<TData>): ApiResponse<TData> {
  return isAxiosResponse(response) ? response.data : response;
}

export async function request<TData>(
  requestFactory: RequestFactory<TData>,
  ...interceptors: ResponseInterceptor<TData>[]
): Promise<ApiResponse<TData>> {
  let response = normalizeResponse(await requestFactory());

  for (const interceptor of interceptors) {
    response = await interceptor(response);
  }

  return response;
}

export default request;
