import axios from "@/api/client"

function ok<T>(p: Promise<T>) { return p.then(() => {}) }

export const roleApi = {
  list: () => axios.get("/admin/api/roles").then((r) => r.data),
  create: (data: Record<string, string>) => ok(axios.post("/admin/roles/save", new URLSearchParams(data))),
  update: (id: number, data: Record<string, string>) => ok(axios.put(`/admin/roles/${id}`, new URLSearchParams(data))),
  delete: (id: number) => ok(axios.delete(`/admin/roles/${id}`)),
}

export const parcelApi = {
  list: () => axios.get("/admin/api/parcels").then((r) => r.data),
  create: (data: Record<string, string>) => ok(axios.post("/admin/parcels/save", new URLSearchParams(data))),
  update: (id: number, data: Record<string, string>) => ok(axios.put(`/admin/parcels/${id}`, new URLSearchParams(data))),
  delete: (id: number) => ok(axios.delete(`/admin/parcels/${id}`)),
}

export const carrierApi = {
  list: () => axios.get("/admin/api/carriers").then((r) => r.data),
  create: (data: Record<string, string>) => ok(axios.post("/admin/carriers/save", new URLSearchParams(data))),
  update: (id: number, data: Record<string, string>) => ok(axios.put(`/admin/carriers/${id}`, new URLSearchParams(data))),
  delete: (id: number) => ok(axios.delete(`/admin/carriers/${id}`)),
}

export const courierApi = {
  list: () => axios.get("/admin/api/couriers").then((r) => r.data),
  create: (data: Record<string, string>) => ok(axios.post("/admin/couriers/save", new URLSearchParams(data))),
  update: (id: number, data: Record<string, string>) => ok(axios.put(`/admin/couriers/${id}`, new URLSearchParams(data))),
  delete: (id: number) => ok(axios.delete(`/admin/couriers/${id}`)),
}

export const declarantApi = {
  list: () => axios.get("/admin/api/declarants").then((r) => r.data),
  create: (data: Record<string, string>) => ok(axios.post("/admin/declarants/save", new URLSearchParams(data))),
  update: (id: number, data: Record<string, string>) => ok(axios.put(`/admin/declarants/${id}`, new URLSearchParams(data))),
  delete: (id: number) => ok(axios.delete(`/admin/declarants/${id}`)),
}
