import { useQuery } from "@tanstack/react-query"
import client from "@/api/client"

export default function WebhooksPage() {
  const { data: hooks = [] } = useQuery<any[]>({
    queryKey: ["webhooks"],
    queryFn: () => client.get("/admin/api/webhooks").then(r => r.data),
    retry: false,
  })

  return (
    <div>
      <h1 style={{ fontSize: 20, fontWeight: "bold", marginBottom: 16 }}>Webhook 投递</h1>

      <div style={{ background: "white", borderRadius: 8, border: "1px solid #e5e7eb", marginBottom: 16, padding: 20 }}>
        <h2 style={{ fontSize: 15, fontWeight: 600, marginBottom: 12, color: "#374151" }}>配置 Webhook 端点</h2>
        <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
          <input placeholder="回调 URL (https://...)" style={{ flex: 1, minWidth: 200, padding: "8px 12px", border: "1px solid #d1d5db", borderRadius: 6, fontSize: 14 }} />
          <select style={{ padding: "8px 12px", border: "1px solid #d1d5db", borderRadius: 6, fontSize: 14 }}>
            <option>所有事件</option><option>order.created</option><option>parcel.received</option><option>parcel.shipped</option>
          </select>
          <button style={{ padding: "8px 20px", background: "#6366f1", color: "white", border: "none", borderRadius: 6, fontSize: 14, fontWeight: 600, cursor: "pointer" }}>添加</button>
        </div>
      </div>

      <div style={{ background: "white", borderRadius: 8, border: "1px solid #e5e7eb", overflow: "auto" }}>
        <table style={{ width: "100%", borderCollapse: "collapse" }}>
          <thead><tr>
            {["客户","事件","回调URL","状态","响应码","时间"].map(h=>(
              <th key={h} style={{padding:"10px 12px",textAlign:"left",borderBottom:"1px solid #e5e7eb",color:"#6b7280",fontWeight:600,fontSize:13}}>{h}</th>
            ))}
          </tr></thead>
          <tbody>
            {hooks.map((h:any,i:number)=>(
              <tr key={h.id||i} style={{background:i%2===0?"#fafafa":"white"}}>
                <td style={td}>{h.client}</td>
                <td style={td}><span style={{fontSize:12,fontFamily:"monospace",background:"#f3f4f6",padding:"2px 6px",borderRadius:4}}>{h.event}</span></td>
                <td style={{...td,fontFamily:"monospace",fontSize:12,maxWidth:200,overflow:"hidden",textOverflow:"ellipsis",whiteSpace:"nowrap"}}>{h.url}</td>
                <td style={td}><span style={{fontSize:12,padding:"2px 8px",borderRadius:10,background:h.status==='success'?'#dcfce7':'#fef2f2',color:h.status==='success'?'#16a34a':'#ef4444',fontWeight:600}}>{h.status==='success'?'成功':'失败'}</span></td>
                <td style={td}><span style={{fontFamily:"monospace",fontSize:12,color:h.response_code>=400?'#ef4444':'#16a34a'}}>{h.response_code}</span></td>
                <td style={{...td,fontSize:12,color:"#9ca3af"}}>{h.created_at?.slice?.(0,19)||'-'}</td>
              </tr>
            ))}
            {hooks.length===0&&<tr><td colSpan={6} style={{...td,textAlign:"center",color:"#9ca3af"}}>暂无投递记录</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  )
}
const td: React.CSSProperties = { padding: "10px 12px", borderBottom: "1px solid #f3f4f6", fontSize: 13 }
