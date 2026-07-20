import { useAuthStore } from "@/stores/auth"
import { useQuery } from "@tanstack/react-query"
import client from "@/api/client"

export function DashboardPage() {
  const { user } = useAuthStore()
  const { data: stats } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: () => client.get('/admin/api/dashboard/stats').then(r => r.data),
    retry: false,
  })

  const cards = [
    { label: '总订单', value: stats?.total_orders ?? '-', color: '#3b82f6' },
    { label: '包裹数', value: stats?.total_parcels ?? '-', color: '#6366f1' },
    { label: '客户数', value: stats?.total_clients ?? '-', color: '#22c55e' },
    { label: '待处理', value: stats?.pending_parcels ?? '-', color: '#f59e0b' },
    { label: '进行中', value: stats?.active_orders ?? '-', color: '#8b5cf6' },
    { label: '本月营收', value: stats?.total_revenue != null ? `¥${stats.total_revenue.toLocaleString()}` : '-', color: '#14b8a6' },
    { label: '承运商', value: stats?.total_carriers ?? '-', color: '#f43f5e' },
    { label: '快递渠道', value: stats?.total_couriers ?? '-', color: '#ef4444' },
  ]

  return (
    <div>
      <div style={{display:'flex',justifyContent:'space-between',alignItems:'center',marginBottom:20}}>
        <h1 style={{fontSize:20,fontWeight:'bold',color:'#1f2937'}}>仪表盘</h1>
        <span style={{color:'#6b7280',fontSize:13}}>已登录: {user?.real_name || user?.username}</span>
      </div>
      <div style={{display:'grid',gridTemplateColumns:'repeat(auto-fill,minmax(180px,1fr))',gap:12}}>
        {cards.map(c => (
          <div key={c.label} style={{background:'white',borderRadius:8,padding:'16px 20px',border:'1px solid #e5e7eb',borderTop:`3px solid ${c.color}`}}>
            <div style={{color:'#6b7280',fontSize:13,marginBottom:8}}>{c.label}</div>
            <div style={{fontSize:28,fontWeight:'bold',color:'#1f2937'}}>{c.value}</div>
          </div>
        ))}
      </div>
    </div>
  )
}
