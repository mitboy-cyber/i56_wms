import { useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import { Package, CheckCircle, Clock, AlertTriangle, Truck } from 'lucide-react';

export default function InboundBoardPage() {
  const { data: parcels = [] } = useQuery<any[]>({ queryKey: ['inbound-parcels'], queryFn: () => client.get('/admin/api/parcels') });
  const { data: stats } = useQuery<any>({ queryKey: ['inbound-stats'], queryFn: () => client.get('/admin/api/dashboard/stats') });

  const received = parcels.filter((p: any) => p.status === 'received').length;
  const preDeclared = parcels.filter((p: any) => p.status === 'pre_declared').length;
  const weighed = parcels.filter((p: any) => p.status === 'weighed').length;
  const stored = parcels.filter((p: any) => p.status === 'stored').length;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-gray-800">入库看板</h1>
        <span className="text-sm text-gray-500">今日入库: {received + preDeclared} 件 | {new Date().toLocaleTimeString('zh-CN')}</span>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
        {[
          { icon: Truck, label: '预申报', value: preDeclared, color: 'bg-gray-50 border-gray-200' },
          { icon: Package, label: '已签收', value: received, color: 'bg-blue-50 border-blue-200' },
          { icon: Clock, label: '已称重', value: weighed, color: 'bg-teal-50 border-teal-200' },
          { icon: CheckCircle, label: '已上架', value: stored, color: 'bg-green-50 border-green-200' },
          { icon: AlertTriangle, label: '异常', value: '0', color: 'bg-red-50 border-red-200' },
        ].map(c => (
          <div key={c.label} className={`rounded-lg border p-4 ${c.color}`}>
            <div className="flex items-center gap-2 mb-1"><c.icon size={18} className="opacity-60" /><span className="text-xs font-medium opacity-70">{c.label}</span></div>
            <div className="text-2xl font-bold">{c.value}</div>
          </div>
        ))}
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b"><tr>
            <th className="px-4 py-3 text-left font-medium text-gray-600">快递单号</th>
            <th className="px-4 py-3 text-left font-medium text-gray-600">品名</th>
            <th className="px-4 py-3 text-right font-medium text-gray-600">重量(kg)</th>
            <th className="px-4 py-3 text-left font-medium text-gray-600">快递</th>
            <th className="px-4 py-3 text-right font-medium text-gray-600">状态</th>
            <th className="px-4 py-3 text-right font-medium text-gray-600">时间</th>
          </tr></thead>
          <tbody className="divide-y">
            {parcels.map((p: any) => (
              <tr key={p.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-mono text-xs">{p.tracking_number}</td>
                <td className="px-4 py-3">{p.product_name}</td>
                <td className="px-4 py-3 text-right">{p.actual_weight}</td>
                <td className="px-4 py-3">{p.courier_code || '顺丰速运'}</td>
                <td className="px-4 py-3 text-right">
                  <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${
                    p.status === 'stored' ? 'bg-green-100 text-green-700' : p.status === 'received' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-600'
                  }`}>{String(p.status)}</span>
                </td>
                <td className="px-4 py-3 text-right text-xs text-gray-500">{p.created_at ? new Date(p.created_at).toLocaleDateString('zh-CN') : '-'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
