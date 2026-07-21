import { useState } from "react"
import { useQuery } from "@tanstack/react-query"
import client from "@/api/client"

export default function NotificationsPage() {
  const { data: notes = [] } = useQuery<any[]>({
    queryKey: ["notifications"],
    queryFn: () => client.get("/admin/api/notifications").then(r => r.data),
    retry: false,
  })

  return (
    <div>
      <h1 style={{ fontSize: 20, fontWeight: "bold", marginBottom: 16 }}>通知管理</h1>
      <div style={{ display: "flex", gap: 12, marginBottom: 16 }}>
        <div style={{ background: "white", borderRadius: 8, border: "1px solid #e5e7eb", padding: 16, flex: 1 }}>
          <div style={{ color: "#6b7280", fontSize: 13 }}>总通知</div>
          <div style={{ fontSize: 28, fontWeight: "bold", color: "#1f2937" }}>{notes.length}</div>
        </div>
        <div style={{ background: "white", borderRadius: 8, border: "1px solid #e5e7eb", padding: 16, flex: 1 }}>
          <div style={{ color: "#6b7280", fontSize: 13 }}>已发送</div>
          <div style={{ fontSize: 28, fontWeight: "bold", color: "#16a34a" }}>{notes.filter((n:any) => n.sent).length}</div>
        </div>
        <div style={{ background: "white", borderRadius: 8, border: "1px solid #e5e7eb", padding: 16, flex: 1 }}>
          <div style={{ color: "#6b7280", fontSize: 13 }}>待发送</div>
          <div style={{ fontSize: 28, fontWeight: "bold", color: "#f59e0b" }}>{notes.filter((n:any) => !n.sent).length}</div>
        </div>
      </div>

      <div style={{ background: "white", borderRadius: 8, border: "1px solid #e5e7eb", overflow: "auto" }}>
        <table style={{ width: "100%", borderCollapse: "collapse" }}>
          <thead><tr>
            {["标题","类型","优先级","范围","渠道","发送人","状态","时间"].map(h=>(
              <th key={h} style={{padding:"10px 12px",textAlign:"left",borderBottom:"1px solid #e5e7eb",color:"#6b7280",fontWeight:600,fontSize:13}}>{h}</th>
            ))}
          </tr></thead>
          <tbody>
            {notes.map((n:any,i:number)=>(
              <tr key={n.id||i} style={{background:i%2===0?"#fafafa":"white"}}>
                <td style={td}>{n.title}</td>
                <td style={td}><span style={{fontSize:12,padding:"2px 8px",borderRadius:10,background:n.type==='公告'?'#dbeafe':n.type==='系统通知'?'#f3f4f6':'#fef3c7',color:'#374151',fontWeight:600}}>{n.type}</span></td>
                <td style={td}><span style={{color:n.priority==='紧急'?'#ef4444':'#6b7280',fontWeight:n.priority==='紧急'?600:400}}>{n.priority}</span></td>
                <td style={td}><span style={{fontSize:12}}>{n.scope}</span></td>
                <td style={td}>{n.channel}</td>
                <td style={td}>{n.sender}</td>
                <td style={td}><span style={{fontSize:12,padding:"2px 8px",borderRadius:10,background:n.sent?'#dcfce7':'#fef3c7',color:n.sent?'#16a34a':'#d97706',fontWeight:600}}>{n.sent?'已发送':'待发送'}</span></td>
                <td style={{...td,fontSize:12,color:"#9ca3af"}}>{n.send_time?.slice?.(0,16)||'-'}</td>
              </tr>
            ))}
            {notes.length===0&&<tr><td colSpan={8} style={{...td,textAlign:"center",color:"#9ca3af"}}>暂无通知</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  )
}
const td: React.CSSProperties = { padding: "10px 12px", borderBottom: "1px solid #f3f4f6", fontSize: 13 }
