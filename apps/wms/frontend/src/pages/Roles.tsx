import { CrudPage } from "@/components/CrudPage"
import { roleApi } from "@/api/others"

interface Role {
  id: number
  name: string
  description: string
  is_active: boolean
}

export function RolesPage() {
  return (
    <CrudPage<Role>
      config={{
        title: "角色管理",
        columns: [
          { key: "name", label: "角色名称" },
          { key: "description", label: "描述" },
        ],
        fields: [
          { name: "name", label: "角色名称" },
          { name: "description", label: "描述" },
        ],
        api: roleApi,
        getRowId: (r) => r.id,
      }}
    />
  )
}
