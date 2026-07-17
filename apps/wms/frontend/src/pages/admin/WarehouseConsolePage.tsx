1|import { useQuery } from '@tanstack/react-query';
2|import client from '@/api/client';
3|import { Package, Wrench, CheckCircle, Clock, AlertTriangle, BarChart3 } from 'lucide-react';
4|
5|export default function WarehouseConsolePage() {
6|  const { data: stats } = useQuery<any>({ queryKey: ['wh-console-stats'], queryFn: () => client.get('/admin/api/dashboard/stats') });
7|  const { data: orders = [] } = useQuery<any[]>({ queryKey: ['wh-console-orders'], queryFn: () => client.get('/admin/api/orders') });
8|  const { data: parcels = [] } = useQuery<any[]>({ queryKey: ['wh-console-parcels'], queryFn: () => client.get('/admin/api/parcels') });
9|
10|  const pendingOps = orders.filter((o: any) => o.status !== 'completed' && o.status !== 'shipped').length;
11|  const completeOps = orders.filter((o: any) => o.status === 'completed' || o.status === 'shipped').length;
12|
13|  // Hourly rate simulation
14|  const hourlyRate = Math.round(parcels.length / Math.max(1, Math.ceil((Date.now() - new Date('2026-07-08').getTime()) / 3600000)));
15|
16|  return (
17|    <div className="space-y-6">
18|      <div className="flex items-center justify-between">
19|        <h1 className="text-xl font-bold text-gray-800">仓库作业台</h1>
20|        <div className="flex items-center gap-2">
21|          <span className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
22|          <span className="text-sm text-green-700 font-medium">系统运行中</span>
23|        </div>
24|      </div>
25|
26|      <div className="grid grid-cols-2 md:grid-cols-6 gap-4">
27|        {[
28|          { icon: Package, label: '总包裹', value: stats?.total_parcels || 0, color: 'bg-blue-50 border-blue-200' },
29|          { icon: Clock, label: '进行中', value: pendingOps, color: 'bg-yellow-50 border-yellow-200' },
30|          { icon: CheckCircle, label: '已完成', value: completeOps, color: 'bg-green-50 border-green-200' },
31|          { icon: Wrench, label: '待处理', value: stats?.pending_parcels || 0, color: 'bg-orange-50 border-orange-200' },
32|          { icon: AlertTriangle, label: '异常', value: '0', color: 'bg-red-50 border-red-200' },
33|          { icon: BarChart3, label: '时均处理', value: `${hourlyRate}件`, color: 'bg-indigo-50 border-indigo-200' },
34|        ].map(c => (
35|          <div key={c.label} className={`rounded-lg border p-3 ${c.color}`}>
36|            <div className="flex items-center gap-1.5 mb-1"><c.icon size={16} className="opacity-60" /><span className="text-xs font-medium opacity-70">{c.label}</span></div>
37|            <div className="text-xl font-bold">{c.value}</div>
38|          </div>
39|        ))}
40|      </div>
41|
42|      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
43|        {/* Active orders queue */}
44|        <div className="bg-white rounded-lg shadow p-4">
45|          <h3 className="text-sm font-semibold text-gray-600 mb-3">待处理工单队列</h3>
46|          <div className="space-y-2">
47|            {pendingOps > 0 ? orders.filter((o: any) => o.status !== 'completed' && o.status !== 'shipped').slice(0, 5).map((o: any) => (
48|              <div key={o.id} className="flex items-center justify-between p-2 rounded bg-gray-50 hover:bg-gray-100">
49|                <div><span className="font-mono text-xs text-gray-600">{o.order_no}</span><span className="mx-2 text-gray-300">|</span><span className="text-sm">{o.recipient_name}</span></div>
50|                <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${
51|                  o.status === 'pending_picking' ? 'bg-yellow-100 text-yellow-700' :
52|                  o.status === 'pending_packing' ? 'bg-orange-100 text-orange-700' :
53|                  'bg-blue-100 text-blue-700'
54|                }`}>{String(o.status).replace(/_/g, ' ')}</span>
55|              </div>
56|            )) : <p className="text-sm text-gray-400 py-4 text-center">暂无待处理工单</p>}
57|          </div>
58|        </div>
59|
60|        {/* Recent completions */}
61|        <div className="bg-white rounded-lg shadow p-4">
62|          <h3 className="text-sm font-semibold text-gray-600 mb-3">最近完成</h3>
63|          <div className="space-y-2">
64|            {orders.filter((o: any) => o.status === 'completed' || o.status === 'shipped').slice(0, 5).map((o: any) => (
65|              <div key={o.id} className="flex items-center justify-between p-2 rounded bg-green-50">
66|                <div><span className="font-mono text-xs text-gray-600">{o.order_no}</span><span className="mx-2 text-gray-300">|</span><span className="text-sm">{o.recipient_name}</span></div>
67|                <span className="text-xs text-green-700 font-medium">✓ 已完成</span>
68|              </div>
69|            ))}
70|            {orders.filter((o: any) => o.status === 'completed' || o.status === 'shipped').length === 0 && (
71|              <p className="text-sm text-gray-400 py-4 text-center">暂无完成记录</p>
72|            )}
73|          </div>
74|        </div>
75|      </div>
76|    </div>
77|  );
78|}
79|