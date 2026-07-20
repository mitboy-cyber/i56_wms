import axios from "axios"

// JSON API client — overrides the default x-www-form-urlencoded Content-Type
const jsonClient = axios.create({
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
})

export const roleApi = {
  list: () => jsonClient.get("/admin/api/roles").then((r) => r.data),
  create: (data: Record<string, unknown>) => jsonClient.post("/admin/api/roles", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => jsonClient.put(`/admin/api/roles/${id}`, data).then(() => {}),
  delete: (id: number) => jsonClient.delete(`/admin/api/roles/${id}`).then(() => {}),
}

export const parcelApi = {
  list: () => jsonClient.get("/admin/api/parcels").then((r) => r.data),
  create: (data: Record<string, unknown>) => jsonClient.post("/admin/api/parcels", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => jsonClient.put(`/admin/api/parcels/${id}`, data).then(() => {}),
  delete: (id: number) => jsonClient.delete(`/admin/api/parcels/${id}`).then(() => {}),
}

export const carrierApi = {
  list: () => jsonClient.get("/admin/api/carriers").then((r) => r.data),
  create: (data: Record<string, unknown>) => jsonClient.post("/admin/api/carriers", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => jsonClient.put(`/admin/api/carriers/${id}`, data).then(() => {}),
  delete: (id: number) => jsonClient.delete(`/admin/api/carriers/${id}`).then(() => {}),
}

export const courierApi = {
  list: () => jsonClient.get("/admin/api/couriers").then((r) => r.data),
  create: (data: Record<string, unknown>) => jsonClient.post("/admin/api/couriers", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => jsonClient.put(`/admin/api/couriers/${id}`, data).then(() => {}),
  delete: (id: number) => jsonClient.delete(`/admin/api/couriers/${id}`).then(() => {}),
}

export const declarantApi = {
  list: () => jsonClient.get("/admin/api/declarants").then((r) => r.data),
  create: (data: Record<string, unknown>) => jsonClient.post("/admin/api/declarants", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => jsonClient.put(`/admin/api/declarants/${id}`, data).then(() => {}),
  delete: (id: number) => jsonClient.delete(`/admin/api/declarants/${id}`).then(() => {}),
}
