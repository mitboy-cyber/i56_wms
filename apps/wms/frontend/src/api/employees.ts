import axios from "@/api/client"

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
  list: () => axios.get<Employee[]>("/admin/api/employees").then((r) => r.data),
  create: (data: Record<string, string>) => axios.post("/admin/employees/save", new URLSearchParams(data)).then(() => {}),
  update: (id: number, data: Record<string, string>) => axios.put(`/admin/employees/${id}`, new URLSearchParams(data)).then(() => {}),
  delete: (id: number) => axios.delete(`/admin/employees/${id}`).then(() => {}),
}
