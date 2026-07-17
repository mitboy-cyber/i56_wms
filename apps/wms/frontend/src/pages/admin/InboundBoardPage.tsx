1|import { useQuery } from '@tanstack/react-query';
2|import client from '@/api/client';
3|import { Package, CheckCircle, Clock, AlertTriangle, Truck } from 'lucide-react';
4|
5|export default function InboundBoardPage() {
6|  const { data: parcels = [] } = useQuery<any[]>({ queryKey: ['inbound-parcels'], queryFn: () => client.get('/admin/api/parcels') });
7|  const { data: stats } = useQuery<any>({ queryKey: ['inbound-stats'], queryFn: () => client.get('/admin/api/dashboard/stats') });
8|
9|  const received = parcels.filter((p: any) => p.status === 'received').length;
10|  const preDeclared = parcels.filter((p: any) => p.status === 'pre_declared').length;
11|  const weighed = parcels.filter((p: any) => p.status === 'weighed').length;
12|  const stored = parcels.filter((p: any) => p.status === 'stored').length;
13|
14|  return (
15|    <div className="space-y-6">
16|      <div className="flex items-center justify-between">
17|        <h1 className="text-xl font-bold text-gray-800">入库看板</h1>
18|        <span className="text-sm text-gray-500">今日入库: {received + preDeclared} 件 | {new Date().toLocaleTimeString('zh-CN')}</span>
19|      </div>
20|
21|      <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
22|        {[
23|          { icon: Truck, label: '预申报', value: preDeclared, color: 'bg-gray-50 border-gray-200' },
24|          { icon: Package, label: '已签收', value: received, color: 'bg-blue-50 border-blue-200' },
25|          { icon: Clock, label: '已称重', value: weighed, color: 'bg-teal-50 border-teal-200' },
26|          { icon: CheckCircle, label: '已上架', value: stored, color: 'bg-green-50 border-green-200' },
27|          { icon: AlertTriangle, label: '异常', value: '0', color: 'bg-red-50 border-red-200' },
28|        ].map(c => (
29|          <div key={c.label} className={`rounded-lg border p-4 ${c.color}`}>
30|            <div className="flex items-center gap-2 mb-1"><c.icon size={18} className="opacity-60" /><span className="text-xs font-medium opacity-70">{c.label}</span></div>
31|            <div className="text-2xl font-bold">{c.value}</div>
32|          </div>
33|        ))}
34|      </div>
35|
36|      <div className="bg-white rounded-lg shadow overflow-hidden">
37|        <table className="w-full text-sm">
38|          <thead className="bg-gray-50 border-b"><tr>
39|            <th className="px-4 py-3 text-left font-medium text-gray-600">快递单号</th>
40|            <th className="px-4 py-3 text-left font-medium text-gray-600">品名</th>
41|            <th className="px-4 py-3 text-right font-medium text-gray-600">重量(kg)</th>
42|            <th className="px-4 py-3 text-left font-medium text-gray-600">快递</th>
43|            <th className="px-4 py-3 text-right font-medium text-gray-600">状态</th>
44|            <th className="px-4 py-3 text-right font-medium text-gray-600">时间</th>
45|          </tr></thead>
46|          <tbody className="divide-y">
47|            {parcels.map((p: any) => (
48|              <tr key={p.id} className="hover:bg-gray-50">
49|                <td className="px-4 py-3 font-mono text-xs">{p.tracking_number}</td>
50|                <td className="px-4 py-3">{p.product_name}</td>
51|                <td className="px-4 py-3 text-right">{p.actual_weight}</td>
52|                <td className="px-4 py-3">{p.courier_code || '顺丰速运'}</td>
53|                <td className="px-4 py-3 text-right">
54|                  <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${
55|                    p.status === 'stored' ? 'bg-green-100 text-green-700' : p.status === 'received' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-600'
56|                  }`}>{String(p.status)}</span>
57|                </td>
58|                <td className="px-4 py-3 text-right text-xs text-gray-500">{p.created_at ? new Date(p.created_at).toLocaleDateString('zh-CN') : '-'}</td>
59|              </tr>
60|            ))}
61|          </tbody>
62|        </table>
63|      </div>
64|    </div>
65|  );
66|}
67|