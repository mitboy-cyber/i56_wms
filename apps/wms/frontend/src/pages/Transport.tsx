import { CrudPage } from "@/components/CrudPage"
import { carrierApi, courierApi, declarantApi } from "@/api/others"

interface Transport {
  id: number
  name: string
  code: string
  contact: string
  phone: string
}

export function CarriersPage() {
  return (
    <CrudPage<Transport>
      config={{
        title: "承运商管理",
        columns: [
          { key: "name", label: "名称" },
          { key: "code", label: "编码" },
          { key: "contact", label: "联系人" },
          { key: "phone", label: "电话" },
        ],
        fields: [
          { name: "name", label: "名称" },
          { name: "code", label: "编码" },
          { name: "contact", label: "联系人" },
          { name: "phone", label: "电话" },
        ],
        api: carrierApi,
        getRowId: (r) => r.id,
      }}
    />
  )
}

export function CouriersPage() {
  return (
    <CrudPage<Transport>
      config={{
        title: "快递公司",
        columns: [
          { key: "name", label: "名称" },
          { key: "code", label: "编码" },
          { key: "contact", label: "联系人" },
          { key: "phone", label: "电话" },
        ],
        fields: [
          { name: "name", label: "名称" },
          { name: "code", label: "编码" },
          { name: "contact", label: "联系人" },
          { name: "phone", label: "电话" },
        ],
        api: courierApi,
        getRowId: (r) => r.id,
      }}
    />
  )
}

export function DeclarantsPage() {
  return (
    <CrudPage<Transport>
      config={{
        title: "申报人管理",
        columns: [
          { key: "name", label: "名称" },
          { key: "code", label: "编码" },
          { key: "contact", label: "联系人" },
          { key: "phone", label: "电话" },
        ],
        fields: [
          { name: "name", label: "名称" },
          { name: "code", label: "编码" },
          { name: "contact", label: "联系人" },
          { name: "phone", label: "电话" },
        ],
        api: declarantApi,
        getRowId: (r) => r.id,
      }}
    />
  )
}
