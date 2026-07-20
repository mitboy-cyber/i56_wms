import client from "@/api/client"

export interface Employee {
  id: number
  username: string
  real_name: string
  role_name: string
  role_id: number
  email: string
  phone: string
  is_active: boolean
}

export const employeeApi = {
  list: () => client.get<Employee[]>("/admin/api/employees").then((r) => r.data),
  create: (data: Record<string, string>) => client.post("/admin/api/employees", new URLSearchParams(data)).then(() => {}),
  update: (id: number, data: Record<string, string>) => client.put(`/admin/api/employees/${id}`, new URLSearchParams(data)).then(() => {}),
  delete: (id: number) => client.delete(`/admin/api/employees/${id}`).then(() => {}),
}
