import { useQuery } from "@tanstack/react-query"
import client from "@/api/client"

export default function WarehouseBoardPage() {
  const { data: stats } = useQuery<any>({
    queryKey: ["wh-board-stats"],
    queryFn: () => client.get("/admin/api/dashboard/stats").then(r => r.data),
    retry: false,
  })
  const { data: orders = [] } = useQuery<any[]>({
    queryKey: ["wh-board-orders"],
    queryFn: () => client.get("/admin/api/orders").then(r => r.data),
    retry: false,
  })

  const picking = Array.isArray(orders) ? orders.filter((o: any) => o.status === "pending_picking").length : 0
  const packing = Array.isArray(orders) ? orders.filter((o: any) => o.status === "pending_packing").length : 0
  const loading = Array.isArray(orders) ? orders.filter((o: any) => o.status === "pending_loading" || o.status === "loaded").length : 0

  const cards = [
    { label: "已上架", value: stats?.total_parcels ?? "-", color: "#22c55e", icon: "📦" },
    { label: "待拣货", value: picking, color: "#eab308", icon: "🕐" },
    { label: "待打包", value: packing, color: "#f97316", icon: "📦" },
    { label: "装箱中", value: loading, color: "#8b5cf6", icon: "🚢" },
  ]

  return (
    <div>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 20 }}>
        <h1 style={{ fontSize: 20, fontWeight: "bold" }}>仓库看板</h1>
        <span style={{ color: "#6b7280", fontSize: 13 }}>厦门仓 | {new Date().toLocaleTimeString("zh-CN")}</span>
      </div>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(160px, 1fr))", gap: 12 }}>
        {cards.map((c) => (
          <div key={c.label} style={{ background: "white", borderRadius: 8, padding: "16px 20px", border: "1px solid #e5e7eb", borderTop: `3px solid ${c.color}` }}>
            <div style={{ fontSize: 20, marginBottom: 4 }}>{c.icon}</div>
            <div style={{ color: "#6b7280", fontSize: 13, marginBottom: 4 }}>{c.label}</div>
            <div style={{ fontSize: 28, fontWeight: "bold" }}>{c.value}</div>
          </div>
        ))}
      </div>
    </div>
  )
}
