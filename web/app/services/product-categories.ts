import { api, type ApiResponse } from "~/services/api";

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

export async function getProductCategories(): Promise<ProductCategory[]> {
  const response = await api
    .get("product-categories/")
    .json<ApiResponse<ProductCategory[]>>();

  return response.data;
}

export async function getProductCategory(id: string): Promise<ProductCategory> {
  const response = await api
    .get(`product-categories/${encodeURIComponent(id)}`)
    .json<ApiResponse<ProductCategory>>();

  return response.data;
}

export async function createProductCategory(
  input: ProductCategoryInput,
): Promise<ProductCategory> {
  const response = await api
    .post("product-categories/", { json: input })
    .json<ApiResponse<ProductCategory>>();

  return response.data;
}

export async function updateProductCategory(
  id: string,
  input: ProductCategoryInput,
): Promise<ProductCategory> {
  const response = await api
    .put(`product-categories/${encodeURIComponent(id)}`, { json: input })
    .json<ApiResponse<ProductCategory>>();

  return response.data;
}

export async function deleteProductCategory(id: string): Promise<null> {
  const response = await api
    .delete(`product-categories/${encodeURIComponent(id)}`)
    .json<ApiResponse<null>>();

  return response.data;
}
