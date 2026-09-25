import { api, type ApiResponse } from "~/services/api";

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

export async function getWarehouses(): Promise<Warehouse[]> {
  const response = await api.get("warehouses/").json<ApiResponse<Warehouse[]>>();

  return response.data;
}

export async function getWarehouse(id: string): Promise<Warehouse> {
  const response = await api
    .get(`warehouses/${encodeURIComponent(id)}`)
    .json<ApiResponse<Warehouse>>();

  return response.data;
}

export async function createWarehouse(input: WarehouseInput): Promise<Warehouse> {
  const response = await api
    .post("warehouses/", { json: input })
    .json<ApiResponse<Warehouse>>();

  return response.data;
}

export async function updateWarehouse(id: string, input: WarehouseInput): Promise<Warehouse> {
  const response = await api
    .put(`warehouses/${encodeURIComponent(id)}`, { json: input })
    .json<ApiResponse<Warehouse>>();

  return response.data;
}

export async function deleteWarehouse(id: string): Promise<null> {
  const response = await api
    .delete(`warehouses/${encodeURIComponent(id)}`)
    .json<ApiResponse<null>>();

  return response.data;
}
