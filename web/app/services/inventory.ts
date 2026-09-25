import { api, type ApiResponse } from "~/services/api";

export type InventoryBalance = {
  warehouse_id: string;
  variant_id: string;
  quantity: number;
  created_at: string;
  updated_at: string;
};

export type StockMovementType = "adjustment_in" | "adjustment_out";

export type StockMovement = {
  id: string;
  warehouse_id: string;
  variant_id: string;
  movement_type: StockMovementType;
  quantity: number;
  balance_after: number;
  note: string;
  created_at: string;
};

export type InventoryAdjustment = {
  balance: InventoryBalance;
  movement: StockMovement;
};

export type InventoryFilters = {
  warehouse_id?: string;
  variant_id?: string;
};

export type AdjustmentType = "adjustment_in" | "adjustment_out";

export type AdjustmentInput = {
  warehouse_id: string;
  variant_id: string;
  type: AdjustmentType;
  quantity: number;
  note: string;
};

export async function getInventoryBalances(
  filters?: InventoryFilters,
): Promise<InventoryBalance[]> {
  const response = await api
    .get("inventory/balances/", { searchParams: filters })
    .json<ApiResponse<InventoryBalance[]>>();

  return response.data;
}

export async function getInventoryBalance(
  warehouseId: string,
  variantId: string,
): Promise<InventoryBalance> {
  const response = await api
    .get(
      `inventory/balances/${encodeURIComponent(warehouseId)}/${encodeURIComponent(variantId)}`,
    )
    .json<ApiResponse<InventoryBalance>>();

  return response.data;
}

export async function adjustInventory(input: AdjustmentInput): Promise<InventoryAdjustment> {
  const response = await api
    .post("inventory/adjustments/", { json: input })
    .json<ApiResponse<InventoryAdjustment>>();

  return response.data;
}

export async function getInventoryMovements(
  filters?: InventoryFilters,
): Promise<StockMovement[]> {
  const response = await api
    .get("inventory/movements/", { searchParams: filters })
    .json<ApiResponse<StockMovement[]>>();

  return response.data;
}

export async function getInventoryMovement(id: string): Promise<StockMovement> {
  const response = await api
    .get(`inventory/movements/${encodeURIComponent(id)}`)
    .json<ApiResponse<StockMovement>>();

  return response.data;
}
