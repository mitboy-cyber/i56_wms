import { useState } from "react"
import { useQuery } from "@tanstack/react-query"
import client from "@/api/client"

// Tracking event timeline
function Timeline({ events }: { events: { time: string; location: string; detail: string; error?: boolean }[] }) {
  return (
    <div style={{ position: "relative", paddingLeft: 24 }}>
      <div style={{ position: "absolute", left: 8, top: 4, bottom: 4, width: 2, background: "#e5e7eb" }} />
      {events.map((ev, i) => (
        <div key={i} style={{ position: "relative", marginBottom: i < events.length - 1 ? 16 : 0 }}>
          <div style={{
            position: "absolute", left: -20, top: 4, width: 10, height: 10, borderRadius: "50%",
            background: ev.error ? "#ef4444" : (i === 0 ? "#10b981" : "#6366f1"),
            border: "2px solid white", boxShadow: "0 0 0 1px #d1d5db"
          }} />
          <div style={{ fontSize: 12, color: "#9ca3af", marginBottom: 2 }}>{ev.time}</div>
          <div style={{ fontSize: 14, fontWeight: 600, color: ev.error ? "#ef4444" : "#1f2937" }}>{ev.location}</div>
          <div style={{ fontSize: 13, color: "#6b7280" }}>{ev.detail}</div>
        </div>
      ))}
    </div>
  )
}

export default function LogisticsTrackingPage() {
  const [selected, setSelected] = useState<string | null>(null)
  const { data: tracks = [] } = useQuery<any[]>({
    queryKey: ["logistics-tracking"],
    queryFn: () => client.get("/admin/api/logistics-tracking").then(r => r.data),
    retry: false,
  })

  const selectedTrack = tracks.find((t: any) => (t.tracking_no || t.id) === selected) || (tracks[0])

  return (
    <div>
      <h1 style={{ fontSize: 20, fontWeight: "bold", marginBottom: 16 }}>物流追踪</h1>

      <div style={{ display: "flex", gap: 16 }}>
        {/* Left: tracking list */}
        <div style={{ flex: "1 1 320px", maxWidth: 380 }}>
          {Array.isArray(tracks) && tracks.map((t: any) => {
            const tid = t.tracking_no || t.id
            const isSelected = tid === selected || (!selected && t === tracks[0])
            return (
              <div key={tid} onClick={() => setSelected(tid)}
                style={{
                  background: isSelected ? "#eef2ff" : "white", padding: 12, borderRadius: 8, marginBottom: 8,
                  border: `1px solid ${isSelected ? "#818cf8" : "#e5e7eb"}`, cursor: "pointer",
                  borderLeft: `4px solid ${t.error_count > 0 ? "#ef4444" : "#10b981"}`
                }}>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                  <span style={{ fontFamily: "monospace", fontSize: 13, fontWeight: 600 }}>{t.tracking_no || tid}</span>
                  {t.error_count > 0 && <span style={{ fontSize: 11, color: "#ef4444", background: "#fef2f2", padding: "1px 6px", borderRadius: 10 }}>⚠ {t.error_count}次失败</span>}
                </div>
                <div style={{ fontSize: 12, color: "#6b7280", marginTop: 4 }}>{t.route || t.line}</div>
                <div style={{ fontSize: 13, fontWeight: 500, marginTop: 2 }}>
                  <span style={{
                    display: "inline-block", padding: "2px 8px", borderRadius: 10, fontSize: 11,
                    background: t.status === "已签收" || t.status === "签收" ? "#dcfce7" : t.status === "运输中" ? "#dbeafe" : "#f3f4f6",
                    color: t.status === "已签收" || t.status === "签收" ? "#16a34a" : t.status === "运输中" ? "#2563eb" : "#6b7280"
                  }}>{t.status || "在途"}</span>
                </div>
              </div>
            )
          })}
          {(!Array.isArray(tracks) || tracks.length === 0) && <p style={{ color: "#9ca3af" }}>暂无追踪记录</p>}
        </div>

        {/* Right: detail panel */}
        <div style={{ flex: "2 1 500px", background: "white", borderRadius: 8, border: "1px solid #e5e7eb", padding: 20 }}>
          {selectedTrack ? (
            <>
              <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
                <div>
                  <h2 style={{ fontSize: 16, fontWeight: "bold", fontFamily: "monospace" }}>{selectedTrack.tracking_no || selectedTrack.id}</h2>
                  <p style={{ fontSize: 13, color: "#6b7280" }}>{selectedTrack.route || selectedTrack.line} — {selectedTrack.carrier}</p>
                </div>
                <div style={{ textAlign: "right" }}>
                  <div style={{ fontSize: 24, fontWeight: "bold", color: selectedTrack.status === "已签收" ? "#16a34a" : "#2563eb" }}>
                    {selectedTrack.status || "在途"}
                  </div>
                  {selectedTrack.order_no && <div style={{ fontSize: 12, color: "#9ca3af" }}>关联: {selectedTrack.order_no}</div>}
                </div>
              </div>

              {/* Timeline */}
              <div style={{ marginBottom: 16 }}>
                <h3 style={{ fontSize: 14, fontWeight: 600, marginBottom: 12, color: "#374151" }}>物流时间线</h3>
                <Timeline events={[
                  { time: "2026-07-21 10:30", location: "派送中", detail: selectedTrack.detail || "快递员派送中", error: false },
                  { time: "2026-07-21 08:15", location: "到达派送网点", detail: selectedTrack.segment || "厦门转运中心→本地配送站", error: false },
                  { time: "2026-07-20 22:00", location: "运输中", detail: "已离开枢纽转运中心", error: false },
                  { time: "2026-07-20 14:00", location: "已出库", detail: "包裹已出库装车", error: selectedTrack.error_count > 0 },
                  { time: "2026-07-20 09:00", location: "已揽件", detail: "快递员已收件", error: false },
                ]} />
              </div>

              {/* Detail info */}
              <div style={{ borderTop: "1px solid #e5e7eb", paddingTop: 12 }}>
                <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8, fontSize: 13 }}>
                  <div><span style={{ color: "#9ca3af" }}>承运商: </span>{selectedTrack.carrier || "-"}</div>
                  <div><span style={{ color: "#9ca3af" }}>出发地: </span>{selectedTrack.from || "厦门"}</div>
                  <div><span style={{ color: "#9ca3af" }}>目的地: </span>{selectedTrack.to || "台北"}</div>
                  <div><span style={{ color: "#9ca3af" }}>失败次数: </span>
                    <span style={{ color: selectedTrack.error_count > 0 ? "#ef4444" : "#16a34a", fontWeight: 600 }}>{selectedTrack.error_count || 0}</span>
                  </div>
                  {selectedTrack.order_no && <div><span style={{ color: "#9ca3af" }}>订单号: </span><span style={{ fontFamily: "monospace" }}>{selectedTrack.order_no}</span></div>}
                  <div><span style={{ color: "#9ca3af" }}>更新时间: </span>{selectedTrack.created_at || selectedTrack.updated_at || "-"}</div>
                </div>
              </div>

              {/* Action buttons */}
              <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
                <button style={{ padding: "6px 16px", borderRadius: 6, border: "1px solid #6366f1", color: "#6366f1", background: "white", cursor: "pointer", fontSize: 13 }}>🔄 刷新查询</button>
                <button style={{ padding: "6px 16px", borderRadius: 6, border: "1px solid #d1d5db", color: "#6b7280", background: "white", cursor: "pointer", fontSize: 13 }}>📋 复制单号</button>
              </div>
            </>
          ) : (
            <p style={{ color: "#9ca3af", textAlign: "center", paddingTop: 40 }}>选择一条物流记录查看详情</p>
          )}
        </div>
      </div>
    </div>
  )
}
