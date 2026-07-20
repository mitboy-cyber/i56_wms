import axios from "@/api/client"

export interface Warehouse {
  id: number
  name: string
  code: string
  address: string
  contact: string
  phone: string
  is_active: boolean
  parcel_count: number
}

export const warehouseApi = {
  list: () => axios.get<Warehouse[]>("/admin/api/warehouses").then((r) => r.data),
  create: (data: Record<string, string>) => axios.post("/admin/api/warehouses", new URLSearchParams(data)).then(() => {}),
  update: (id: number, data: Record<string, string>) => axios.put(`/admin/api/warehouses/${id}`, new URLSearchParams(data)).then(() => {}),
  delete: (id: number) => axios.delete(`/admin/api/warehouses/${id}`).then(() => {}),
}
