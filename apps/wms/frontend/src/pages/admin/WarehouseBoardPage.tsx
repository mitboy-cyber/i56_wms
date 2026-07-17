1|import { useQuery } from '@tanstack/react-query';
2|import client from '@/api/client';
3|import { Package, ShoppingCart, Truck, Boxes, Container, AlertTriangle, Clock, TrendingUp } from 'lucide-react';
4|
5|export default function WarehouseBoardPage() {
6|  const { data: stats } = useQuery<any>({ queryKey: ['wh-board-stats'], queryFn: () => client.get('/admin/api/dashboard/stats') });
7|  const { data: orders = [] } = useQuery<any[]>({ queryKey: ['wh-board-orders'], queryFn: () => client.get('/admin/api/orders') });
8|  const { data: parcels = [] } = useQuery<any[]>({ queryKey: ['wh-board-parcels'], queryFn: () => client.get('/admin/api/parcels') });
9|
10|  const inStock = parcels.filter((p: any) => p.status === 'stored').length;
11|  const picking = orders.filter((o: any) => o.status === 'pending_picking').length;
12|  const packing = orders.filter((o: any) => o.status === 'pending_packing').length;
13|  const loading = orders.filter((o: any) => o.status === 'pending_loading' || o.status === 'loaded').length;
14|
15|  const cards = [
16|    { icon: Boxes, label: '已上架', value: inStock, color: 'bg-green-50 border-green-200 text-green-700' },
17|    { icon: Clock, label: '待拣货', value: picking, color: 'bg-yellow-50 border-yellow-200 text-yellow-700' },
18|    { icon: Package, label: '待打包', value: packing, color: 'bg-orange-50 border-orange-200 text-orange-700' },
19|    { icon: Container, label: '装箱中', value: loading, color: 'bg-purple-50 border-purple-200 text-purple-700' },
20|  ];
21|
22|  return (
23|    <div className="space-y-6">
24|      <div className="flex items-center justify-between">
25|        <h1 className="text-xl font-bold text-gray-800">仓库看板</h1>
26|        <span className="text-sm text-gray-500">厦门仓 | {new Date().toLocaleTimeString('zh-CN')}</span>
27|      </div>
28|
29|      {/* ── Warehouse KPIs ── */}
30|      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
31|        {cards.map(c => (
32|          <div key={c.color} className={`rounded-lg border p-4 ${c.color}`}>
33|            <div className="flex items-center gap-2 mb-1"><c.icon size={18} className="opacity-60" />
34|              <span className="text-xs font-medium opacity-70">{c.label}</span>
35|            </div>
36|            <div className="text-2xl font-bold">{c.value}</div>
37|          </div>
38|        ))}
39|      </div>
40|
41|      {/* ── Parcel location grid ── */}
42|      <div className="bg-white rounded-lg shadow p-4">
43|        <h3 className="text-sm font-semibold text-gray-600 mb-4">包裹仓位分布</h3>
44|        <div className="grid grid-cols-5 gap-2">
45|          {(() => {
46|            const locations = ['A-01', 'A-02', 'A-03', 'B-01', 'B-02', 'B-03', 'C-01', 'C-02', 'D-01', 'D-02'];
47|            return locations.map((loc, i) => {
48|              const count = i < inStock ? 1 : 0; // Simulate location assignment
49|              return (
50|                <div key={loc} className={`rounded-lg border-2 p-3 text-center ${count > 0 ? 'border-green-300 bg-green-50' : 'border-gray-200 bg-gray-50'}`}>
51|                  <div className="text-xs text-gray-500">{loc}</div>
52|                  <div className={`text-lg font-bold ${count > 0 ? 'text-green-600' : 'text-gray-300'}`}>
53|                    {count > 0 ? '●' : '○'}
54|                  </div>
55|                </div>
56|              );
57|            });
58|          })()}
59|        </div>
60|      </div>
61|
62|      {/* ── Recent parcels ── */}
63|      <div className="bg-white rounded-lg shadow p-4">
64|        <h3 className="text-sm font-semibold text-gray-600 mb-3">最近入库包裹</h3>
65|        <table className="w-full text-sm">
66|          <thead className="border-b"><tr>
67|            <th className="text-left py-2">快递单号</th><th className="text-left py-2">品名</th>
68|            <th className="text-right py-2">重量(kg)</th><th className="text-right py-2">状态</th>
69|          </tr></thead>
70|          <tbody className="divide-y">
71|            {parcels.slice(0, 8).map((p: any) => (
72|              <tr key={p.id} className="hover:bg-gray-50">
73|                <td className="py-2 font-mono text-xs">{p.tracking_number}</td>
74|                <td className="py-2">{p.product_name}</td>
75|                <td className="py-2 text-right">{p.actual_weight}</td>
76|                <td className="py-2 text-right">
77|                  <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${p.status === 'stored' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}`}>
78|                    {String(p.status)}
79|                  </span>
80|                </td>
81|              </tr>
82|            ))}
83|          </tbody>
84|        </table>
85|      </div>
86|    </div>
87|  );
88|}
89|