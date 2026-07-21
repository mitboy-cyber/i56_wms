import { useAuthStore } from "@/stores/auth"
import { useQuery } from "@tanstack/react-query"
import client from "@/api/client"

// Pure SVG donut chart component — no external deps
function DonutChart({ data, size = 140 }: { data: { label: string; value: number; color: string }[]; size?: number }) {
  const total = data.reduce((s, d) => s + d.value, 0) || 1
  const r = size / 2 - 15
  const center = size / 2
  const circumference = 2 * Math.PI * r
  let offset = 0

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
      {data.map((d) => {
        const percent = d.value / total
        const dash = circumference * percent
        const o = offset; offset += dash
        return <circle key={d.label} cx={center} cy={center} r={r} fill="none" stroke={d.color} strokeWidth={16}
          strokeDasharray={`${dash} ${circumference - dash}`} strokeDashoffset={-o}
          transform={`rotate(-90 ${center} ${center})`} />
      })}
      <text x={center} y={center - 6} textAnchor="middle" fontSize={24} fontWeight="bold" fill="#1f2937">{total}</text>
      <text x={center} y={center + 14} textAnchor="middle" fontSize={11} fill="#6b7280">总数</text>
    </svg>
  )
}

// Simple horizontal bar chart
function BarChart({ data, maxLabelLen = 8 }: { data: { label: string; value: number; color: string }[]; maxLabelLen?: number }) {
  const max = Math.max(...data.map(d => d.value), 1)
  const barH = 22; const gap = 6; const w = 260
  return (
    <svg width={w + 100} height={data.length * (barH + gap) + 10}>
      {data.map((d, i) => {
        const bw = (d.value / max) * w
        const y = i * (barH + gap)
        return <g key={d.label}>
          <text x={0} y={y + 14} fontSize={11} fill="#6b7280" textAnchor="end">{d.label.slice(0, maxLabelLen)}</text>
          <rect x={8} y={y + 3} width={bw} height={barH} rx={3} fill={d.color} opacity={0.85} />
          <text x={bw + 12} y={y + 16} fontSize={12} fontWeight={600} fill="#374151">{d.value}</text>
        </g>
      }, 0)}
    </svg>
  )
}

export function DashboardPage() {
  const { user } = useAuthStore()
  const { data: stats } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: () => client.get('/admin/api/dashboard/stats').then(r => r.data),
    retry: false,
  })
  const { data: orderStatus } = useQuery({
    queryKey: ['dashboard-order-status'],
    queryFn: () => client.get('/admin/api/dashboard/order-status').then(r => r.data),
    retry: false,
  })
  const { data: revenue } = useQuery({
    queryKey: ['dashboard-revenue'],
    queryFn: () => client.get('/admin/api/dashboard/revenue-by-route').then(r => r.data),
    retry: false,
  })

  const cards = [
    { label: '总订单', value: stats?.total_orders ?? '-', color: '#3b82f6', icon: '📋' },
    { label: '包裹数', value: stats?.total_parcels ?? '-', color: '#6366f1', icon: '📦' },
    { label: '客户数', value: stats?.total_clients ?? '-', color: '#22c55e', icon: '👥' },
    { label: '本月营收', value: stats?.total_revenue != null ? `¥${stats.total_revenue.toLocaleString()}` : '-', color: '#14b8a6', icon: '💰' },
  ]

  // Order status donut data
  const statusColors: Record<string, string> = {
    pending_picking: '#f59e0b', pending_packing: '#f97316', pending_loading: '#8b5cf6',
    loaded: '#6366f1', in_transit: '#06b6d4', customs_clearance: '#ec4899',
    shipped: '#14b8a6', delivered: '#22c55e', completed: '#16a34a',
  }
  const statusLabels: Record<string, string> = { pending_picking:'待拣货', pending_packing:'待打包', pending_loading:'待装车', loaded:'已装车', in_transit:'运输中', customs_clearance:'清关中', shipped:'已发货', delivered:'已送达', completed:'已完成' }
  const donutData = Array.isArray(orderStatus) ? orderStatus.map((s: any) => ({ label: statusLabels[s.status] || s.status, value: s.count, color: statusColors[s.status] || '#9ca3af' })) : []

  // Revenue bar data
  const barData = Array.isArray(revenue) ? revenue.map((r: any) => ({ label: r.route_name || r.route, value: r.total_revenue || 0, color: '#6366f1' })) : []

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <h1 style={{ fontSize: 20, fontWeight: 'bold', color: '#1f2937' }}>仪表盘</h1>
        <span style={{ color: '#6b7280', fontSize: 13 }}>已登录: {user?.real_name || user?.username}</span>
      </div>

      {/* KPI Cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(160px, 1fr))', gap: 12, marginBottom: 20 }}>
        {cards.map(c => (
          <div key={c.label} style={{ background: 'white', borderRadius: 8, padding: '16px 20px', border: '1px solid #e5e7eb', borderTop: `3px solid ${c.color}` }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <div style={{ color: '#6b7280', fontSize: 13, marginBottom: 4 }}>{c.label}</div>
                <div style={{ fontSize: 24, fontWeight: 'bold', color: '#1f2937' }}>{c.value}</div>
              </div>
              <span style={{ fontSize: 22 }}>{c.icon}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Charts row */}
      <div style={{ display: 'flex', gap: 20, flexWrap: 'wrap' }}>
        {/* Order Status Donut */}
        <div style={{ background: 'white', borderRadius: 8, padding: 20, border: '1px solid #e5e7eb', flex: '1 1 280px' }}>
          <h2 style={{ fontSize: 15, fontWeight: 600, marginBottom: 12, color: '#374151' }}>订单状态分布</h2>
          <div style={{ display: 'flex', alignItems: 'center', gap: 20 }}>
            {donutData.length > 0 ? <DonutChart data={donutData} size={140} /> : <p style={{ color: '#9ca3af', fontSize: 13 }}>暂无数据</p>}
            <div>
              {donutData.map(d => (
                <div key={d.label} style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4, fontSize: 12 }}>
                  <span style={{ width: 10, height: 10, borderRadius: '50%', background: d.color, display: 'inline-block' }} />
                  <span style={{ color: '#374151' }}>{d.label}</span>
                  <span style={{ color: '#6b7280', fontWeight: 600, marginLeft: 'auto' }}>{d.value}</span>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Revenue by Route */}
        <div style={{ background: 'white', borderRadius: 8, padding: 20, border: '1px solid #e5e7eb', flex: '2 1 360px' }}>
          <h2 style={{ fontSize: 15, fontWeight: 600, marginBottom: 12, color: '#374151' }}>线路营收 TOP</h2>
          {barData.length > 0 ? <BarChart data={barData.slice(0, 8)} /> : <p style={{ color: '#9ca3af', fontSize: 13 }}>暂无数据</p>}
        </div>
      </div>
    </div>
  )
}
