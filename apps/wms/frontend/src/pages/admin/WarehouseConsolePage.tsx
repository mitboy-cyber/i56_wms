export default function WarehouseConsolePage() {
  return (
    <div>
      <h1 style={{fontSize:20,fontWeight:'bold',marginBottom:16}}>仓库作业台</h1>
      <div style={{display:'grid',gridTemplateColumns:'repeat(auto-fill,minmax(200px,1fr))',gap:12}}>
        {[
          {label:'待收货',count:'3',color:'#3b82f6'},
          {label:'待称重',count:'2',color:'#6366f1'},
          {label:'待上架',count:'5',color:'#22c55e'},
          {label:'待拣货',count:'4',color:'#f59e0b'},
          {label:'待打包',count:'2',color:'#f97316'},
          {label:'待装车',count:'1',color:'#8b5cf6'},
          {label:'今日完成',count:'18',color:'#14b8a6'},
          {label:'异常',count:'1',color:'#ef4444'},
        ].map(c=>(
          <div key={c.label} style={{background:'white',borderRadius:8,padding:20,border:'1px solid #e5e7eb',borderLeft:`4px solid ${c.color}`}}>
            <div style={{color:'#6b7280',fontSize:13}}>{c.label}</div>
            <div style={{fontSize:32,fontWeight:'bold',color:'#1f2937'}}>{c.count}</div>
          </div>
        ))}
      </div>
    </div>
  )
}
