import axios from "@/api/client"

export const roleApi = {
  list: () => axios.get("/admin/api/roles").then((r) => r.data),
  create: (data: Record<string, unknown>) => axios.post("/admin/api/roles", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => axios.put(`/admin/api/roles/${id}`, data).then(() => {}),
  delete: (id: number) => axios.delete(`/admin/api/roles/${id}`).then(() => {}),
}

export const parcelApi = {
  list: () => axios.get("/admin/api/parcels").then((r) => r.data),
  create: (data: Record<string, unknown>) => axios.post("/admin/api/parcels", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => axios.put(`/admin/api/parcels/${id}`, data).then(() => {}),
  delete: (id: number) => axios.delete(`/admin/api/parcels/${id}`).then(() => {}),
}

export const carrierApi = {
  list: () => axios.get("/admin/api/carriers").then((r) => r.data),
  create: (data: Record<string, unknown>) => axios.post("/admin/api/carriers", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => axios.put(`/admin/api/carriers/${id}`, data).then(() => {}),
  delete: (id: number) => axios.delete(`/admin/api/carriers/${id}`).then(() => {}),
}

export const courierApi = {
  list: () => axios.get("/admin/api/couriers").then((r) => r.data),
  create: (data: Record<string, unknown>) => axios.post("/admin/api/couriers", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => axios.put(`/admin/api/couriers/${id}`, data).then(() => {}),
  delete: (id: number) => axios.delete(`/admin/api/couriers/${id}`).then(() => {}),
}

export const declarantApi = {
  list: () => axios.get("/admin/api/declarants").then((r) => r.data),
  create: (data: Record<string, unknown>) => axios.post("/admin/api/declarants", data).then(() => {}),
  update: (id: number, data: Record<string, unknown>) => axios.put(`/admin/api/declarants/${id}`, data).then(() => {}),
  delete: (id: number) => axios.delete(`/admin/api/declarants/${id}`).then(() => {}),
}
