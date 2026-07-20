import axios from "@/api/client"

export interface Client {
  id: number
  name: string
  code: string
  type: string
  contact: string
  phone: string
  email: string
  is_active: boolean
}

export const clientApi = {
  list: () => axios.get<Client[]>("/admin/api/clients").then((r) => r.data),
  create: (data: Record<string, string>) => axios.post("/admin/api/clients", new URLSearchParams(data)).then(() => {}),
  update: (id: number, data: Record<string, string>) => axios.put(`/admin/api/clients/${id}`, new URLSearchParams(data)).then(() => {}),
  delete: (id: number) => axios.delete(`/admin/api/clients/${id}`).then(() => {}),
}
