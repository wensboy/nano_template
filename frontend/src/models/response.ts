export interface ApiResponse<TData = unknown> {
  code: number;
  message: string;
  data: TData | null;
}

