import { CrudPage } from "@/components/CrudPage"
import { warehouseApi, type Warehouse } from "@/api/warehouses"

export function WarehousesPage() {
  return (
    <CrudPage<Warehouse>
      config={{
        title: "仓库管理",
        columns: [
          { key: "name", label: "仓库" },
          { key: "code", label: "编码" },
          { key: "address", label: "地址" },
          { key: "contact", label: "联系人" },
          { key: "phone", label: "电话" },
          { key: "parcel_count", label: "包裹数", render: (r) => `${r.parcel_count}件` },
        ],
        fields: [
          { name: "name", label: "仓库名称" },
          { name: "code", label: "编码" },
          { name: "address", label: "地址" },
          { name: "contact", label: "联系人" },
          { name: "phone", label: "电话" },
        ],
        api: warehouseApi,
        getRowId: (r) => r.id,
      }}
    />
  )
}
