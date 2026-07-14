import { CrudPage } from "@/components/CrudPage"
import { employeeApi, type Employee } from "@/api/employees"

export function EmployeesPage() {
  return (
    <CrudPage<Employee>
      config={{
        title: "员工管理",
        columns: [
          { key: "real_name", label: "姓名" },
          { key: "username", label: "账号" },
          { key: "role_name", label: "角色" },
          { key: "email", label: "邮箱" },
          { key: "phone", label: "电话" },
        ],
        fields: [
          { name: "username", label: "账号" },
          { name: "password", label: "密码" },
          { name: "real_name", label: "姓名" },
          { name: "role_id", label: "角色", type: "select", options: [] },
          { name: "email", label: "邮箱" },
          { name: "phone", label: "电话" },
        ],
        api: employeeApi,
        getRowId: (r) => r.id,
      }}
    />
  )
}
