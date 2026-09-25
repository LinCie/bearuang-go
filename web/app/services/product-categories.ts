import {
  api,
  apiResult,
  type ApiResult,
  type CommonApiErrorCode,
} from "~/services/api";

export type ProductCategoryStatus = "draft" | "active" | "inactive" | "archived";

export type ProductCategory = {
  id: string;
  parent_id: string | null;
  name: string;
  slug: string;
  description: string;
  status: ProductCategoryStatus;
  created_at: string;
  updated_at: string;
};

export type ProductCategoryInput = {
  parent_id: string | null;
  name: string;
  slug: string;
  description: string;
  status: ProductCategoryStatus;
};

export type ProductCategoryErrorCode =
  | CommonApiErrorCode
  | "unauthorized"
  | "invalid_status"
  | "not_found";

export function getProductCategories(): Promise<
  ApiResult<ProductCategory[], ProductCategoryErrorCode>
> {
  return apiResult<ProductCategory[], ProductCategoryErrorCode>(
    api.get("product-categories/"),
  );
}

export function getProductCategory(
  id: string,
): Promise<ApiResult<ProductCategory, ProductCategoryErrorCode>> {
  return apiResult<ProductCategory, ProductCategoryErrorCode>(
    api.get(`product-categories/${encodeURIComponent(id)}`),
  );
}

export function createProductCategory(
  input: ProductCategoryInput,
): Promise<ApiResult<ProductCategory, ProductCategoryErrorCode>> {
  return apiResult<ProductCategory, ProductCategoryErrorCode>(
    api.post("product-categories/", { json: input }),
  );
}

export function updateProductCategory(
  id: string,
  input: ProductCategoryInput,
): Promise<ApiResult<ProductCategory, ProductCategoryErrorCode>> {
  return apiResult<ProductCategory, ProductCategoryErrorCode>(
    api.put(`product-categories/${encodeURIComponent(id)}`, { json: input }),
  );
}

export function deleteProductCategory(
  id: string,
): Promise<ApiResult<null, ProductCategoryErrorCode>> {
  return apiResult<null, ProductCategoryErrorCode>(
    api.delete(`product-categories/${encodeURIComponent(id)}`),
  );
}
