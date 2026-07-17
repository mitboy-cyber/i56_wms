1|import { useQuery } from '@tanstack/react-query';
2|import client from '@/api/client';
3|import { Package, ShoppingCart, Users, Warehouse, TrendingUp, DollarSign, AlertTriangle, Truck } from 'lucide-react';
4|
5|const STATUS_CN: Record<string, string> = {
6|  pending_picking: '待拣货', pending_packing: '待打包', pending_loading: '待装柜',
7|  in_transit: '运输中', customs_clearance: '清关中', completed: '已完成',
8|  shipped: '已发货', loaded: '已装柜', cancelled: '已取消',
9|  stored: '已上架', packed: '已打包', weighed: '已称重', received: '已签收',
10|  pre_declared: '预申报', delivered: '已送达',
11|};
12|
13|export function DashboardPage() {
14|  const { data: stats } = useQuery<any>({ queryKey: ['dashboard-stats'], queryFn: () => client.get('/admin/api/dashboard/stats') });
15|  const { data: orders = [] } = useQuery<any[]>({ queryKey: ['dashboard-orders'], queryFn: () => client.get('/admin/api/orders') });
16|  const { data: parcels = [] } = useQuery<any[]>({ queryKey: ['dashboard-parcels'], queryFn: () => client.get('/admin/api/parcels') });
17|
18|  const totalRevenue = stats?.total_revenue || 0;
19|  const todayOrders = orders.filter((o: any) => o.status === 'in_transit' || o.status === 'pending_picking').length;
20|  const pendingParcels = parcels.filter((p: any) => p.status === 'stored' || p.status === 'received').length;
21|
22|  const cards = [
23|    { icon: ShoppingCart, label: '总订单', value: stats?.total_orders || 0, color: 'blue' },
24|    { icon: Package, label: '包裹数', value: stats?.total_parcels || 0, color: 'indigo' },
25|    { icon: Users, label: '客户数', value: stats?.total_clients || 0, color: 'green' },
26|    { icon: Warehouse, label: '待处理包裹', value: stats?.pending_parcels || 0, color: 'amber' },
27|    { icon: Truck, label: '进行中订单', value: stats?.active_orders || 0, color: 'purple' },
28|    { icon: DollarSign, label: '本月营收', value: `¥${totalRevenue.toLocaleString()}`, color: 'teal' },
29|    { icon: TrendingUp, label: '承运商', value: stats?.total_carriers || 0, color: 'rose' },
30|    { icon: AlertTriangle, label: '快递渠道', value: stats?.total_couriers || 0, color: 'red' },
31|  ];
32|
33|  return (
34|    <div className="space-y-6">
35|      <div className="flex items-center justify-between">
36|        <h1 className="text-xl font-bold text-gray-800">仪表盘</h1>
37|        <span className="text-sm text-gray-500">最后更新: {new Date().toLocaleTimeString('zh-CN')}</span>
38|      </div>
39|
40|      {/* ── KPI cards ── */}
41|      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
42|        {cards.map((c) => {
43|          const colors: Record<string,string> = {
44|            blue:'bg-blue-50 border-blue-200', indigo:'bg-indigo-50 border-indigo-200',
45|            green:'bg-green-50 border-green-200', amber:'bg-amber-50 border-amber-200',
46|            purple:'bg-purple-50 border-purple-200', teal:'bg-teal-50 border-teal-200',
47|            rose:'bg-rose-50 border-rose-200', red:'bg-red-50 border-red-200',
48|          };
49|          return (
50|            <div key={c.label} className={`rounded-lg border p-4 ${colors[c.color]}`}>
51|              <div className="flex items-center gap-2 mb-2">
52|                <c.icon size={20} className="opacity-60" />
53|                <span className="text-sm font-medium opacity-70">{c.label}</span>
54|              </div>
55|              <div className="text-2xl font-bold">{c.value}</div>
56|            </div>
57|          );
58|        })}
59|      </div>
60|
61|      {/* ── Order status breakdown ── */}
62|      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
63|        <div className="bg-white rounded-lg shadow p-4">
64|          <h3 className="text-sm font-semibold text-gray-600 mb-3">订单状态分布</h3>
65|          <div className="space-y-2">
66|            {(() => {
67|              const counts: Record<string, number> = {}; orders.forEach((o: any) => { const s = String(o.status); counts[s] = (counts[s] || 0) + 1; });
68|              const colors: Record<string, string> = { pending_picking: '#eab308', pending_packing: '#f59e0b', pending_loading: '#f97316', in_transit: '#3b82f6', customs_clearance: '#14b8a6', completed: '#22c55e', shipped: '#6366f1', loaded: '#a855f7', cancelled: '#ef4444' };
69|              return Object.entries(counts).map(([k, v]) => (
70|                <div key={k} className="flex items-center gap-3">
71|                  <span className="text-xs text-gray-600 w-20">{STATUS_CN[k] || k}</span>
72|                  <div className="flex-1 bg-gray-100 rounded-full h-4">
73|                    <div className="h-4 rounded-full" style={{ width: `${(v/orders.length)*100}%`, backgroundColor: colors[k] || '#9ca3af' }} />
74|                  </div>
75|                  <span className="text-xs font-medium w-6 text-right">{v}</span>
76|                </div>
77|              ));
78|            })()}
79|          </div>
80|        </div>
81|        <div className="bg-white rounded-lg shadow p-4">
82|          <h3 className="text-sm font-semibold text-gray-600 mb-3">包裹状态分布</h3>
83|          <div className="space-y-2">
84|            {(() => {
85|              const counts: Record<string, number> = {}; parcels.forEach((p: any) => { const s = String(p.status); counts[s] = (counts[s] || 0) + 1; });
86|              const colors: Record<string, string> = { stored: '#22c55e', packed: '#3b82f6', weighed: '#14b8a6', received: '#6366f1', pre_declared: '#9ca3af', delivered: '#22c55e', shipped: '#a855f7' };
87|              return Object.entries(counts).map(([k, v]) => (
88|                <div key={k} className="flex items-center gap-3">
89|                  <span className="text-xs text-gray-600 w-20">{STATUS_CN[k] || k}</span>
90|                  <div className="flex-1 bg-gray-100 rounded-full h-4">
91|                    <div className="h-4 rounded-full" style={{ width: `${(v/parcels.length)*100}%`, backgroundColor: colors[k] || '#9ca3af' }} />
92|                  </div>
93|                  <span className="text-xs font-medium w-6 text-right">{v}</span>
94|                </div>
95|              ));
96|            })()}
97|          </div>
98|        </div>
99|      </div>
100|
101|
      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="text-sm font-semibold text-gray-600 mb-3">最近订单</h3>
        <table className="w-full text-sm">
          <thead className="border-b"><tr>
            <th className="text-left py-2 font-medium text-gray-500">订单号</th>
            <th className="text-left py-2 font-medium text-gray-500">收件人</th>
            <th className="text-right py-2 font-medium text-gray-500">金额</th>
            <th className="text-right py-2 font-medium text-gray-500">状态</th>
          </tr></thead>
          <tbody className="divide-y">
            {orders.slice(0, 5).map((o: any) => (
              <tr key={o.id} className="hover:bg-gray-50">
                <td className="py-2 font-mono text-xs">{o.order_no}</td>
                <td className="py-2">{o.recipient_name}</td>
                <td className="py-2 text-right">¥{Number(o.total_price).toFixed(2)}</td>
                <td className="py-2 text-right">
                  <span className="inline-flex px-2 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700">
                    {STATUS_CN[String(o.status)] || String(o.status).replace(/_/g, ' ')}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
