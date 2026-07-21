import { useQuery } from "@tanstack/react-query"
import client from "@/api/client"

export default function OrderProfitReport() {
  const { data: stats } = useQuery<any>({
    queryKey: ["order-profit-report"],
    queryFn: () => client.get("/admin/api/finance/order-profit").then(r => r.data),
  })

  return (
    <div>
      <h1 style={{ fontSize: 20, fontWeight: "bold", marginBottom: 16 }}>集运订单盈利</h1>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(180px, 1fr))", gap: 12, marginBottom: 24 }}>
        {[
          { label: "总订单", value: stats?.orders ?? "-", color: "#3b82f6" },
          { label: "总营收", value: stats?.total_revenue != null ? `¥${stats.total_revenue.toLocaleString()}` : "-", color: "#22c55e" },
          { label: "均单价", value: stats?.avg_order != null ? `¥${stats.avg_order.toFixed(2)}` : "-", color: "#8b5cf6" },
        ].map(c => (
          <div key={c.label} style={{ background: "white", borderRadius: 8, padding: "16px 20px", border: "1px solid #e5e7eb", borderTop: `3px solid ${c.color}` }}>
            <div style={{ color: "#6b7280", fontSize: 13, marginBottom: 4 }}>{c.label}</div>
            <div style={{ fontSize: 28, fontWeight: "bold" }}>{c.value}</div>
          </div>
        ))}
      </div>
      <p style={{ color: "#6b7280", fontSize: 13 }}>数据来源: 实时订单聚合</p>
    </div>
  )
}
