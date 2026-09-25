import {
  api,
  apiResult,
  type ApiResult,
  type CommonApiErrorCode,
} from "~/services/api";

export type WarehouseStatus = "active" | "inactive" | "archived";

export type Warehouse = {
  id: string;
  code: string;
  name: string;
  description: string;
  address: string | null;
  status: WarehouseStatus;
  created_at: string;
  updated_at: string;
};

export type WarehouseInput = {
  code: string;
  name: string;
  description: string;
  address: string | null;
  status: WarehouseStatus;
};

export type WarehouseErrorCode =
  | CommonApiErrorCode
  | "unauthorized"
  | "invalid_code"
  | "invalid_name"
  | "invalid_status"
  | "duplicate_code"
  | "warehouse_in_use"
  | "warehouse_has_stock"
  | "not_found";

export function getWarehouses(): Promise<
  ApiResult<Warehouse[], WarehouseErrorCode>
> {
  return apiResult<Warehouse[], WarehouseErrorCode>(api.get("warehouses/"));
}

export function getWarehouse(
  id: string,
): Promise<ApiResult<Warehouse, WarehouseErrorCode>> {
  return apiResult<Warehouse, WarehouseErrorCode>(
    api.get(`warehouses/${encodeURIComponent(id)}`),
  );
}

export function createWarehouse(
  input: WarehouseInput,
): Promise<ApiResult<Warehouse, WarehouseErrorCode>> {
  return apiResult<Warehouse, WarehouseErrorCode>(
    api.post("warehouses/", { json: input }),
  );
}

export function updateWarehouse(
  id: string,
  input: WarehouseInput,
): Promise<ApiResult<Warehouse, WarehouseErrorCode>> {
  return apiResult<Warehouse, WarehouseErrorCode>(
    api.put(`warehouses/${encodeURIComponent(id)}`, { json: input }),
  );
}

export function deleteWarehouse(
  id: string,
): Promise<ApiResult<null, WarehouseErrorCode>> {
  return apiResult<null, WarehouseErrorCode>(
    api.delete(`warehouses/${encodeURIComponent(id)}`),
  );
}
