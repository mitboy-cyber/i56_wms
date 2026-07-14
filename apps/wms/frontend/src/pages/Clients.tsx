import { CrudPage } from "@/components/CrudPage"
import { clientApi, type Client } from "@/api/clients"

export function ClientsPage() {
  return (
    <CrudPage<Client>
      config={{
        title: "客户管理",
        columns: [
          { key: "name", label: "名称" },
          { key: "code", label: "编码" },
          { key: "type", label: "类型" },
          { key: "contact", label: "联系人" },
          { key: "phone", label: "电话" },
        ],
        fields: [
          { name: "name", label: "客户名称" },
          { name: "code", label: "编码" },
          { name: "type", label: "类型" },
          { name: "contact", label: "联系人" },
          { name: "phone", label: "电话" },
          { name: "email", label: "邮箱" },
        ],
        api: clientApi,
        getRowId: (r) => r.id,
      }}
    />
  )
}
