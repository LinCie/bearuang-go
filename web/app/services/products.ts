import { api, type ApiResponse } from "~/services/api";

export type ProductStatus = "draft" | "active" | "inactive" | "archived";
export type ProductVariantStatus = "draft" | "active" | "inactive" | "archived";

export type Product = {
  id: string;
  category_id: string | null;
  name: string;
  slug: string;
  description: string;
  status: ProductStatus;
  created_at: string;
  updated_at: string;
};

export type ProductInput = {
  category_id: string | null;
  name: string;
  slug: string;
  description: string;
  status: ProductStatus;
};

export type ProductVariant = {
  id: string;
  product_id: string;
  sku: string;
  name: string;
  price: number;
  stock: number;
  unit: string;
  status: ProductVariantStatus;
  created_at: string;
  updated_at: string;
};

export type ProductVariantInput = {
  sku: string;
  name: string;
  price: number;
  unit: string;
  status: ProductVariantStatus;
};

export async function getProducts(): Promise<Product[]> {
  const response = await api.get("products/").json<ApiResponse<Product[]>>();

  return response.data;
}

export async function getProduct(id: string): Promise<Product> {
  const response = await api
    .get(`products/${encodeURIComponent(id)}`)
    .json<ApiResponse<Product>>();

  return response.data;
}

export async function createProduct(input: ProductInput): Promise<Product> {
  const response = await api
    .post("products/", { json: input })
    .json<ApiResponse<Product>>();

  return response.data;
}

export async function updateProduct(id: string, input: ProductInput): Promise<Product> {
  const response = await api
    .put(`products/${encodeURIComponent(id)}`, { json: input })
    .json<ApiResponse<Product>>();

  return response.data;
}

export async function deleteProduct(id: string): Promise<null> {
  const response = await api
    .delete(`products/${encodeURIComponent(id)}`)
    .json<ApiResponse<null>>();

  return response.data;
}

export async function getProductVariants(productId: string): Promise<ProductVariant[]> {
  const response = await api
    .get(`products/${encodeURIComponent(productId)}/variants/`)
    .json<ApiResponse<ProductVariant[]>>();

  return response.data;
}

export async function getProductVariant(
  productId: string,
  variantId: string,
): Promise<ProductVariant> {
  const response = await api
    .get(
      `products/${encodeURIComponent(productId)}/variants/${encodeURIComponent(variantId)}`,
    )
    .json<ApiResponse<ProductVariant>>();

  return response.data;
}

export async function createProductVariant(
  productId: string,
  input: ProductVariantInput,
): Promise<ProductVariant> {
  const response = await api
    .post(`products/${encodeURIComponent(productId)}/variants/`, { json: input })
    .json<ApiResponse<ProductVariant>>();

  return response.data;
}

export async function updateProductVariant(
  productId: string,
  variantId: string,
  input: ProductVariantInput,
): Promise<ProductVariant> {
  const response = await api
    .put(
      `products/${encodeURIComponent(productId)}/variants/${encodeURIComponent(variantId)}`,
      { json: input },
    )
    .json<ApiResponse<ProductVariant>>();

  return response.data;
}

export async function deleteProductVariant(
  productId: string,
  variantId: string,
): Promise<null> {
  const response = await api
    .delete(
      `products/${encodeURIComponent(productId)}/variants/${encodeURIComponent(variantId)}`,
    )
    .json<ApiResponse<null>>();

  return response.data;
}
