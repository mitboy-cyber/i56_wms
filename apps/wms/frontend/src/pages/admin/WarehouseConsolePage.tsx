import { useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import { Package, Wrench, CheckCircle, Clock, AlertTriangle, BarChart3 } from 'lucide-react';

export default function WarehouseConsolePage() {
  const { data: stats } = useQuery<any>({ queryKey: ['wh-console-stats'], queryFn: () => client.get('/admin/api/dashboard/stats') });
  const { data: orders = [] } = useQuery<any[]>({ queryKey: ['wh-console-orders'], queryFn: () => client.get('/admin/api/orders') });
  const { data: parcels = [] } = useQuery<any[]>({ queryKey: ['wh-console-parcels'], queryFn: () => client.get('/admin/api/parcels') });

  const pendingOps = orders.filter((o: any) => o.status !== 'completed' && o.status !== 'shipped').length;
  const completeOps = orders.filter((o: any) => o.status === 'completed' || o.status === 'shipped').length;

  // Hourly rate simulation
  const hourlyRate = Math.round(parcels.length / Math.max(1, Math.ceil((Date.now() - new Date('2026-07-08').getTime()) / 3600000)));

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-gray-800">仓库作业台</h1>
        <div className="flex items-center gap-2">
          <span className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
          <span className="text-sm text-green-700 font-medium">系统运行中</span>
        </div>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-6 gap-4">
        {[
          { icon: Package, label: '总包裹', value: stats?.total_parcels || 0, color: 'bg-blue-50 border-blue-200' },
          { icon: Clock, label: '进行中', value: pendingOps, color: 'bg-yellow-50 border-yellow-200' },
          { icon: CheckCircle, label: '已完成', value: completeOps, color: 'bg-green-50 border-green-200' },
          { icon: Wrench, label: '待处理', value: stats?.pending_parcels || 0, color: 'bg-orange-50 border-orange-200' },
          { icon: AlertTriangle, label: '异常', value: '0', color: 'bg-red-50 border-red-200' },
          { icon: BarChart3, label: '时均处理', value: `${hourlyRate}件`, color: 'bg-indigo-50 border-indigo-200' },
        ].map(c => (
          <div key={c.label} className={`rounded-lg border p-3 ${c.color}`}>
            <div className="flex items-center gap-1.5 mb-1"><c.icon size={16} className="opacity-60" /><span className="text-xs font-medium opacity-70">{c.label}</span></div>
            <div className="text-xl font-bold">{c.value}</div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Active orders queue */}
        <div className="bg-white rounded-lg shadow p-4">
          <h3 className="text-sm font-semibold text-gray-600 mb-3">待处理工单队列</h3>
          <div className="space-y-2">
            {pendingOps > 0 ? orders.filter((o: any) => o.status !== 'completed' && o.status !== 'shipped').slice(0, 5).map((o: any) => (
              <div key={o.id} className="flex items-center justify-between p-2 rounded bg-gray-50 hover:bg-gray-100">
                <div><span className="font-mono text-xs text-gray-600">{o.order_no}</span><span className="mx-2 text-gray-300">|</span><span className="text-sm">{o.recipient_name}</span></div>
                <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${
                  o.status === 'pending_picking' ? 'bg-yellow-100 text-yellow-700' :
                  o.status === 'pending_packing' ? 'bg-orange-100 text-orange-700' :
                  'bg-blue-100 text-blue-700'
                }`}>{String(o.status).replace(/_/g, ' ')}</span>
              </div>
            )) : <p className="text-sm text-gray-400 py-4 text-center">暂无待处理工单</p>}
          </div>
        </div>

        {/* Recent completions */}
        <div className="bg-white rounded-lg shadow p-4">
          <h3 className="text-sm font-semibold text-gray-600 mb-3">最近完成</h3>
          <div className="space-y-2">
            {orders.filter((o: any) => o.status === 'completed' || o.status === 'shipped').slice(0, 5).map((o: any) => (
              <div key={o.id} className="flex items-center justify-between p-2 rounded bg-green-50">
                <div><span className="font-mono text-xs text-gray-600">{o.order_no}</span><span className="mx-2 text-gray-300">|</span><span className="text-sm">{o.recipient_name}</span></div>
                <span className="text-xs text-green-700 font-medium">✓ 已完成</span>
              </div>
            ))}
            {orders.filter((o: any) => o.status === 'completed' || o.status === 'shipped').length === 0 && (
              <p className="text-sm text-gray-400 py-4 text-center">暂无完成记录</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
