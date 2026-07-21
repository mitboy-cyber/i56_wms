import client from "@/api/client"

const STATUS_MAP: Record<string, string> = {
  pending: "待处理", in_progress: "进行中", completed: "已完成", cancelled: "已取消"
}

export default function WorkOrdersPage() {
  return (
    <div>
      <h1 style={{fontSize:20,fontWeight:"bold",marginBottom:16}}>员工任务监控</h1>
      <div style={{display:"grid",gridTemplateColumns:"repeat(auto-fill,minmax(220px,1fr))",gap:12}}>
        {[
          {title:"入库任务",desc:"收货→称重→上架",count:5,color:"#3b82f6"},
          {title:"拣货任务",desc:"按订单拣选",count:3,color:"#f59e0b"},
          {title:"打包任务",desc:"复核+打包",count:2,color:"#f97316"},
          {title:"装车任务",desc:"装柜+发运",count:1,color:"#8b5cf6"},
          {title:"异常处理",desc:"破损/丢失/错发",count:2,color:"#ef4444"},
        ].map(t=>(
          <div key={t.title} style={{background:"white",borderRadius:8,padding:20,border:"1px solid #e5e7eb"}}>
            <div style={{display:"flex",justifyContent:"space-between",alignItems:"center",marginBottom:8}}>
              <span style={{fontWeight:600,fontSize:16}}>{t.title}</span>
              <span style={{background:t.color,color:"white",padding:"2px 10px",borderRadius:12,fontSize:13}}>{t.count}</span>
            </div>
            <div style={{color:"#6b7280",fontSize:13}}>{t.desc}</div>
          </div>
        ))}
      </div>
    </div>
  )
}
