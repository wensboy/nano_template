export type PaginationParams = {
  page?: number;
  page_size?: number;
};

export type Pagination<TItem> = {
  total: number;
  page: number;
  page_size: number;
  items: TItem[];
};

export function assertNonEmptyValue(value: string, fieldLabel: string) {
  if (!value.trim()) {
    throw new Error(`${fieldLabel}不能为空`);
  }
}

export function trimRequired(value: string, fieldLabel: string) {
  const normalized = value.trim();
  assertNonEmptyValue(normalized, fieldLabel);
  return normalized;
}

export function trimOptional(value?: string) {
  return value?.trim() ?? "";
}

export function withPagination(params: PaginationParams = {}) {
  return {
    page: params.page ?? 1,
    page_size: params.page_size ?? 10,
  };
}
