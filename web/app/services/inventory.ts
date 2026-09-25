import {
  api,
  apiResult,
  type ApiResult,
  type CommonApiErrorCode,
} from "~/services/api";

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

export type InventoryErrorCode =
  | CommonApiErrorCode
  | "unauthorized"
  | "invalid_id"
  | "invalid_type"
  | "invalid_quantity"
  | "invalid_warehouse_id"
  | "invalid_variant_id"
  | "warehouse_not_found"
  | "warehouse_not_active"
  | "variant_not_found"
  | "insufficient_stock"
  | "not_found";

export function getInventoryBalances(
  filters?: InventoryFilters,
): Promise<ApiResult<InventoryBalance[], InventoryErrorCode>> {
  return apiResult<InventoryBalance[], InventoryErrorCode>(
    api.get("inventory/balances/", { searchParams: filters }),
  );
}

export function getInventoryBalance(
  warehouseId: string,
  variantId: string,
): Promise<ApiResult<InventoryBalance, InventoryErrorCode>> {
  return apiResult<InventoryBalance, InventoryErrorCode>(
    api.get(
      `inventory/balances/${encodeURIComponent(warehouseId)}/${encodeURIComponent(variantId)}`,
    ),
  );
}

export function adjustInventory(
  input: AdjustmentInput,
): Promise<ApiResult<InventoryAdjustment, InventoryErrorCode>> {
  return apiResult<InventoryAdjustment, InventoryErrorCode>(
    api.post("inventory/adjustments/", { json: input }),
  );
}

export function getInventoryMovements(
  filters?: InventoryFilters,
): Promise<ApiResult<StockMovement[], InventoryErrorCode>> {
  return apiResult<StockMovement[], InventoryErrorCode>(
    api.get("inventory/movements/", { searchParams: filters }),
  );
}

export function getInventoryMovement(
  id: string,
): Promise<ApiResult<StockMovement, InventoryErrorCode>> {
  return apiResult<StockMovement, InventoryErrorCode>(
    api.get(`inventory/movements/${encodeURIComponent(id)}`),
  );
}
