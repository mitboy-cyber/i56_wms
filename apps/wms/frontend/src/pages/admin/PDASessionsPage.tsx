import { useEffect, useState } from "react"
import client from "@/api/client"

interface PDASession {
  id: string; warehouse: string; operator: string; device: string
  status: string; page: string; area: string; location: string
  login_at: string; last_beat: string
}

const statusColors: Record<string, { bg: string; color: string; label: string }> = {
  scanning: { bg: "#dcfce7", color: "#16a34a", label: "扫描中" },
  online: { bg: "#dbeafe", color: "#2563eb", label: "在线" },
  idle: { bg: "#fef3c7", color: "#d97706", label: "待机" },
  error: { bg: "#fef2f2", color: "#ef4444", label: "异常" },
}

export default function PDASessionsPage() {
  const [sessions, setSessions] = useState<PDASession[]>([])
  const [polling, setPolling] = useState(true)
  const [lastRefresh, setLastRefresh] = useState("")

  const fetchSessions = async () => {
    try {
      const r = await client.get("/admin/api/pda-sessions")
      const d = r.data
      const arr = Array.isArray(d) ? d : (Array.isArray(d?.sessions) ? d.sessions : (Array.isArray(d?.data) ? d.data : []))
      setSessions(arr)
      setLastRefresh(new Date().toLocaleTimeString("zh-CN"))
    } catch {}
  }

  useEffect(() => {
    fetchSessions()
    if (!polling) return
    const interval = setInterval(fetchSessions, 3000)
    return () => clearInterval(interval)
  }, [polling])

  return (
    <div>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
        <h1 style={{ fontSize: 20, fontWeight: "bold" }}>PDA 在线会话</h1>
        <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
          <button onClick={() => setPolling(!polling)} style={{ padding: "4px 12px", borderRadius: 6, border: "1px solid #d1d5db", background: polling ? "#dcfce7" : "#fef2f2", color: polling ? "#16a34a" : "#ef4444", fontSize: 12, cursor: "pointer" }}>
            {polling ? "🔄 实时" : "⏸ 暂停"}
          </button>
          <span style={{ fontSize: 12, color: "#9ca3af" }}>刷新: {lastRefresh || "-"}</span>
        </div>
      </div>

      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(310px, 1fr))", gap: 12 }}>
        {sessions.map(s => {
          const sc = statusColors[s.status] || statusColors.online
          return (
            <div key={s.id} style={{ background: "white", borderRadius: 10, border: "1px solid #e5e7eb", padding: 16 }}>
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", marginBottom: 8 }}>
                <div>
                  <div style={{ fontSize: 15, fontWeight: 600, color: "#1f2937" }}>{s.operator}</div>
                  <div style={{ fontSize: 11, color: "#9ca3af", fontFamily: "monospace" }}>{s.id}</div>
                </div>
                <span style={{ fontSize: 11, padding: "2px 10px", borderRadius: 10, background: sc.bg, color: sc.color, fontWeight: 600 }}>{sc.label}</span>
              </div>
              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 6, fontSize: 12 }}>
                <div><span style={{ color: "#9ca3af" }}>仓库: </span>{s.warehouse || "-"}</div>
                <div><span style={{ color: "#9ca3af" }}>设备: </span><span style={{ fontFamily: "monospace" }}>{s.device || "-"}</span></div>
                <div><span style={{ color: "#9ca3af" }}>页面: </span>{s.page || "-"}</div>
                <div><span style={{ color: "#9ca3af" }}>区域: </span>{s.area || "-"}</div>
                <div style={{ gridColumn: "span 2" }}>
                  <span style={{ color: "#9ca3af" }}>货位: </span>
                  <span style={{ fontFamily: "monospace", background: "#f3f4f6", padding: "1px 6px", borderRadius: 4 }}>{s.location || "-"}</span>
                </div>
                <div><span style={{ color: "#9ca3af" }}>登录: </span>{s.login_at?.slice?.(11, 19) || "-"}</div>
                <div><span style={{ color: "#9ca3af" }}>心跳: </span>{s.last_beat?.slice?.(11, 19) || "-"}</div>
              </div>
            </div>
          )
        })}
      </div>
      {sessions.length === 0 && <p style={{ color: "#9ca3af", textAlign: "center", marginTop: 40 }}>暂无在线 PDA 设备</p>}
    </div>
  )
}
