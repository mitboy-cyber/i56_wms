import { useQuery } from "@tanstack/react-query"
import client from "@/api/client"

export default function InboundBoardPage() {
  const { data: parcels = [] } = useQuery<any[]>({
    queryKey: ["inbound-parcels"],
    queryFn: () => client.get("/admin/api/parcels").then(r => r.data),
    retry: false,
  })
  const { data: stats } = useQuery<any>({
    queryKey: ["inbound-stats"],
    queryFn: () => client.get("/admin/api/dashboard/stats").then(r => r.data),
    retry: false,
  })

  const arr = Array.isArray(parcels) ? parcels : []
  const received = arr.filter((p: any) => p.status === "received").length
  const preDeclared = arr.filter((p: any) => p.status === "pre_declared").length
  const weighed = arr.filter((p: any) => p.status === "weighed").length
  const stored = arr.filter((p: any) => p.status === "stored").length

  const stages = [
    { label: "预申报", value: preDeclared, color: "#9ca3af" },
    { label: "已签收", value: received, color: "#3b82f6" },
    { label: "已称重", value: weighed, color: "#14b8a6" },
    { label: "已上架", value: stored, color: "#22c55e" },
    { label: "异常", value: "0", color: "#ef4444" },
  ]

  const statusLabel: Record<string, string> = {
    stored: "已上架", received: "已签收", weighed: "已称重",
    signed: "已签收", pre_declared: "预申报",
  }
  const statusColor: Record<string, string> = {
    stored: "#16a34a", received: "#2563eb", weighed: "#0d9488",
    signed: "#2563eb", pre_declared: "#6b7280",
  }

  return (
    <div>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
        <h1 style={{ fontSize: 20, fontWeight: "bold" }}>入库看板</h1>
        <span style={{ color: "#6b7280", fontSize: 13 }}>今日入库: {received + preDeclared} 件 | {new Date().toLocaleTimeString("zh-CN")}</span>
      </div>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(140px, 1fr))", gap: 10, marginBottom: 20 }}>
        {stages.map(c => (
          <div key={c.label} style={{ background: "white", borderRadius: 8, padding: "16px 20px", border: "1px solid #e5e7eb", borderTop: `3px solid ${c.color}` }}>
            <div style={{ color: "#6b7280", fontSize: 13, marginBottom: 4 }}>{c.label}</div>
            <div style={{ fontSize: 28, fontWeight: "bold" }}>{c.value}</div>
          </div>
        ))}
      </div>
      <div style={{ background: "white", borderRadius: 8, padding: 16, border: "1px solid #e5e7eb" }}>
        <table style={{ width: "100%", fontSize: 14, borderCollapse: "collapse" }}>
          <thead>
            <tr style={{ borderBottom: "1px solid #e5e7eb" }}>
              {["快递单号", "品名", "重量(kg)", "快递", "状态", "时间"].map(h => (
                <th key={h} style={{ padding: "8px 12px", textAlign: "left", color: "#6b7280", fontWeight: 600, fontSize: 13 }}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {arr.map((p: any) => (
              <tr key={p.id} style={{ borderBottom: "1px solid #f3f4f6" }}>
                <td style={{ padding: "8px 12px", fontFamily: "monospace", fontSize: 13 }}>{p.tracking_number}</td>
                <td style={{ padding: "8px 12px" }}>{p.product_name}</td>
                <td style={{ padding: "8px 12px" }}>{p.actual_weight}</td>
                <td style={{ padding: "8px 12px" }}>{p.courier_code || "顺丰速运"}</td>
                <td style={{ padding: "8px 12px" }}>
                  <span style={{ display: "inline-block", padding: "2px 8px", borderRadius: 10, fontSize: 12, fontWeight: 500, background: (statusColor[p.status] || "#6b7280") + "18", color: statusColor[p.status] || "#6b7280" }}>
                    {statusLabel[p.status] || p.status}
                  </span>
                </td>
                <td style={{ padding: "8px 12px", color: "#9ca3af", fontSize: 13 }}>{p.created_at ? new Date(p.created_at).toLocaleDateString("zh-CN") : "-"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
