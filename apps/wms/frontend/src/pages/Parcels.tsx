import { CrudPage } from "@/components/CrudPage"
import { parcelApi } from "@/api/others"

interface Parcel {
  id: number
  tracking_no: string
  sender: string
  receiver: string
  status: string
}

export function ParcelsPage() {
  return (
    <CrudPage<Parcel>
      config={{
        title: "包裹管理",
        columns: [
          { key: "tracking_no", label: "运单号" },
          { key: "sender", label: "发件人" },
          { key: "receiver", label: "收件人" },
          { key: "status", label: "状态" },
        ],
        fields: [
          { name: "tracking_no", label: "运单号" },
          { name: "sender", label: "发件人" },
          { name: "receiver", label: "收件人" },
        ],
        api: parcelApi,
        getRowId: (r) => r.id,
      }}
    />
  )
}
