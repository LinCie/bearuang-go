import {
  api,
  apiResult,
  type ApiResult,
  type CommonApiErrorCode,
} from "~/services/api";

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

export type ProductErrorCode =
  | CommonApiErrorCode
  | "unauthorized"
  | "invalid_status"
  | "product_in_use"
  | "not_found";

export type ProductVariantErrorCode =
  | CommonApiErrorCode
  | "unauthorized"
  | "invalid_status"
  | "variant_in_use"
  | "not_found";

export function getProducts(): Promise<ApiResult<Product[], ProductErrorCode>> {
  return apiResult<Product[], ProductErrorCode>(api.get("products/"));
}

export function getProduct(id: string): Promise<ApiResult<Product, ProductErrorCode>> {
  return apiResult<Product, ProductErrorCode>(
    api.get(`products/${encodeURIComponent(id)}`),
  );
}

export function createProduct(
  input: ProductInput,
): Promise<ApiResult<Product, ProductErrorCode>> {
  return apiResult<Product, ProductErrorCode>(api.post("products/", { json: input }));
}

export function updateProduct(
  id: string,
  input: ProductInput,
): Promise<ApiResult<Product, ProductErrorCode>> {
  return apiResult<Product, ProductErrorCode>(
    api.put(`products/${encodeURIComponent(id)}`, { json: input }),
  );
}

export function deleteProduct(id: string): Promise<ApiResult<null, ProductErrorCode>> {
  return apiResult<null, ProductErrorCode>(
    api.delete(`products/${encodeURIComponent(id)}`),
  );
}

export function getProductVariants(
  productId: string,
): Promise<ApiResult<ProductVariant[], ProductVariantErrorCode>> {
  return apiResult<ProductVariant[], ProductVariantErrorCode>(
    api.get(`products/${encodeURIComponent(productId)}/variants/`),
  );
}

export function getProductVariant(
  productId: string,
  variantId: string,
): Promise<ApiResult<ProductVariant, ProductVariantErrorCode>> {
  return apiResult<ProductVariant, ProductVariantErrorCode>(
    api.get(
      `products/${encodeURIComponent(productId)}/variants/${encodeURIComponent(variantId)}`,
    ),
  );
}

export function createProductVariant(
  productId: string,
  input: ProductVariantInput,
): Promise<ApiResult<ProductVariant, ProductVariantErrorCode>> {
  return apiResult<ProductVariant, ProductVariantErrorCode>(
    api.post(`products/${encodeURIComponent(productId)}/variants/`, { json: input }),
  );
}

export function updateProductVariant(
  productId: string,
  variantId: string,
  input: ProductVariantInput,
): Promise<ApiResult<ProductVariant, ProductVariantErrorCode>> {
  return apiResult<ProductVariant, ProductVariantErrorCode>(
    api.put(
      `products/${encodeURIComponent(productId)}/variants/${encodeURIComponent(variantId)}`,
      { json: input },
    ),
  );
}

export function deleteProductVariant(
  productId: string,
  variantId: string,
): Promise<ApiResult<null, ProductVariantErrorCode>> {
  return apiResult<null, ProductVariantErrorCode>(
    api.delete(
      `products/${encodeURIComponent(productId)}/variants/${encodeURIComponent(variantId)}`,
    ),
  );
}
