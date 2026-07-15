import { useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import { Package, ShoppingCart, Users, Warehouse, TrendingUp, DollarSign, AlertTriangle, Truck } from 'lucide-react';

export function DashboardPage() {
  const { data: stats } = useQuery<any>({ queryKey: ['dashboard-stats'], queryFn: () => client.get('/admin/api/dashboard/stats').then(r => r.data) });
  const { data: orders = [] } = useQuery<any[]>({ queryKey: ['dashboard-orders'], queryFn: () => client.get('/admin/api/orders').then(r => r.data) });
  const { data: parcels = [] } = useQuery<any[]>({ queryKey: ['dashboard-parcels'], queryFn: () => client.get('/admin/api/parcels').then(r => r.data) });

  const totalRevenue = stats?.total_revenue || 0;
  const todayOrders = orders.filter((o: any) => o.status === 'in_transit' || o.status === 'pending_picking').length;
  const pendingParcels = parcels.filter((p: any) => p.status === 'stored' || p.status === 'received').length;

  const cards = [
    { icon: ShoppingCart, label: '总订单', value: stats?.total_orders || 0, color: 'blue' },
    { icon: Package, label: '包裹数', value: stats?.total_parcels || 0, color: 'indigo' },
    { icon: Users, label: '客户数', value: stats?.total_clients || 0, color: 'green' },
    { icon: Warehouse, label: '待处理包裹', value: stats?.pending_parcels || 0, color: 'amber' },
    { icon: Truck, label: '进行中订单', value: stats?.active_orders || 0, color: 'purple' },
    { icon: DollarSign, label: '本月营收', value: `¥${totalRevenue.toLocaleString()}`, color: 'teal' },
    { icon: TrendingUp, label: '承运商', value: stats?.total_carriers || 0, color: 'rose' },
    { icon: AlertTriangle, label: '快递渠道', value: stats?.total_couriers || 0, color: 'red' },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-gray-800">仪表盘</h1>
        <span className="text-sm text-gray-500">最后更新: {new Date().toLocaleTimeString('zh-CN')}</span>
      </div>

      {/* ── KPI cards ── */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {cards.map((c) => {
          const colors: Record<string,string> = {
            blue:'bg-blue-50 border-blue-200', indigo:'bg-indigo-50 border-indigo-200',
            green:'bg-green-50 border-green-200', amber:'bg-amber-50 border-amber-200',
            purple:'bg-purple-50 border-purple-200', teal:'bg-teal-50 border-teal-200',
            rose:'bg-rose-50 border-rose-200', red:'bg-red-50 border-red-200',
          };
          return (
            <div key={c.label} className={`rounded-lg border p-4 ${colors[c.color]}`}>
              <div className="flex items-center gap-2 mb-2">
                <c.icon size={20} className="opacity-60" />
                <span className="text-sm font-medium opacity-70">{c.label}</span>
              </div>
              <div className="text-2xl font-bold">{c.value}</div>
            </div>
          );
        })}
      </div>

      {/* ── Order status breakdown ── */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="bg-white rounded-lg shadow p-4">
          <h3 className="text-sm font-semibold text-gray-600 mb-3">订单状态分布</h3>
          <div className="space-y-2">
            {(() => {
              const counts: Record<string, number> = {}; orders.forEach((o: any) => { const s = String(o.status); counts[s] = (counts[s] || 0) + 1; });
              const statusCN: Record<string, string> = { pending_picking: '待拣货', pending_packing: '待打包', pending_loading: '待装柜', in_transit: '运输中', customs_clearance: '清关中', completed: '已完成', shipped: '已发货', loaded: '已装柜', cancelled: '已取消' };
              const colors: Record<string, string> = { pending_picking: '#eab308', pending_packing: '#f59e0b', pending_loading: '#f97316', in_transit: '#3b82f6', customs_clearance: '#14b8a6', completed: '#22c55e', shipped: '#6366f1', loaded: '#a855f7', cancelled: '#ef4444' };
              return Object.entries(counts).map(([k, v]) => (
                <div key={k} className="flex items-center gap-3">
                  <span className="text-xs text-gray-600 w-20">{statusCN[k] || k.replace(/_/g, ' ')}</span>
                  <div className="flex-1 bg-gray-100 rounded-full h-4">
                    <div className="h-4 rounded-full" style={{ width: `${(v/orders.length)*100}%`, backgroundColor: colors[k] || '#9ca3af' }} />
                  </div>
                  <span className="text-xs font-medium w-6 text-right">{v}</span>
                </div>
              ));
            })()}
          </div>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <h3 className="text-sm font-semibold text-gray-600 mb-3">包裹状态分布</h3>
          <div className="space-y-2">
            {(() => {
              const counts: Record<string, number> = {}; parcels.forEach((p: any) => { const s = String(p.status); counts[s] = (counts[s] || 0) + 1; });
              const statusCN: Record<string, string> = { stored: '已上架', packed: '已打包', weighed: '已称重', received: '已签收', pre_declared: '预申报', delivered: '已送达', shipped: '已发货' };
              const colors: Record<string, string> = { stored: '#22c55e', packed: '#3b82f6', weighed: '#14b8a6', received: '#6366f1', pre_declared: '#9ca3af', delivered: '#22c55e', shipped: '#a855f7' };
              return Object.entries(counts).map(([k, v]) => (
                <div key={k} className="flex items-center gap-3">
                  <span className="text-xs text-gray-600 w-20">{statusCN[k] || k.replace(/_/g, ' ')}</span>
                  <div className="flex-1 bg-gray-100 rounded-full h-4">
                    <div className="h-4 rounded-full" style={{ width: `${(v/parcels.length)*100}%`, backgroundColor: colors[k] || '#9ca3af' }} />
                  </div>
                  <span className="text-xs font-medium w-6 text-right">{v}</span>
                </div>
              ));
            })()}
          </div>
        </div>
      </div>

      {/* ── Recent orders ── */}
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
                    {String(o.status).replace(/_/g, ' ')}
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
