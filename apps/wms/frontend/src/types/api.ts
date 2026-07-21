// Shared API types — mirrors backend/internal/types/api.go

export interface ApiResponse<T = unknown> {
  success: boolean
  data?: T
  error?: string
  fields?: Record<string, string>
}

// ─── Order ────────────────────────────────────────────────────────

export interface Order {
  id: number
  order_no: string
  recipient_name: string
  parcel_count: number
  total_price: number
  total_actual_weight: number
  total_chargeable_weight: number
  tracking_numbers: string
  route_id: number
  warehouse_id: number
  client_id: number
  status: string
  remark: string
  created_at: string
  updated_at: string
}

export interface CreateOrderRequest {
  recipient_name: string    // required, 1-128 chars
  route_id: number          // required, >0
  parcel_count: number      // required, >=1
  total_price?: number
  tracking_numbers?: string
  remark?: string
}

// ─── Parcel ───────────────────────────────────────────────────────

export interface Parcel {
  id: number
  tracking_number: string
  product_name: string
  actual_weight: number
  courier_code: string
  cargo_type: string
  status: string
  warehouse_id: number
  created_at: string
}

export interface CreateParcelRequest {
  tracking_number: string   // required
  product_name: string      // required
  actual_weight?: number
  courier_code?: string
  cargo_type?: string
  warehouse_id?: number
}

// ─── Warehouse ────────────────────────────────────────────────────

export interface Warehouse {
  id: number
  name: string
  code: string
  address: string
  contact: string
  phone: string
}

export interface CreateWarehouseRequest {
  name: string    // required
  code: string    // required
  address: string // required
  contact?: string
  phone?: string
}

// ─── Client ───────────────────────────────────────────────────────

export interface Client {
  id: number
  name: string
  code: string
  contact_name: string
  contact_phone: string
}

export interface CreateClientRequest {
  name: string    // required
  code: string    // required
  contact?: string
  phone?: string
}

// ─── Validation helpers ───────────────────────────────────────────

export interface ValidationErrors {
  [field: string]: string
}

export function validateRequired(value: unknown, fieldName: string): string | null {
  if (value === "" || value === null || value === undefined) {
    return `${fieldName} 为必填项`
  }
  return null
}

export function validateForm(fields: Record<string, [unknown, string]>): ValidationErrors {
  const errors: ValidationErrors = {}
  for (const [key, [value, label]] of Object.entries(fields)) {
    const err = validateRequired(value, label)
    if (err) errors[key] = err
  }
  return errors
}
