import { useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import { Package, ShoppingCart, Users, Warehouse, TrendingUp, DollarSign, AlertTriangle, Truck } from 'lucide-react';

export function DashboardPage() {
  const { data: orders = [] } = useQuery<any[]>({ queryKey: ['dashboard-orders'], queryFn: () => client.get('/admin/api/orders').then(r => r.data) });
  const { data: parcels = [] } = useQuery<any[]>({ queryKey: ['dashboard-parcels'], queryFn: () => client.get('/admin/api/parcels').then(r => r.data) });
  const { data: clients = [] } = useQuery<any[]>({ queryKey: ['dashboard-clients'], queryFn: () => client.get('/admin/api/client-accounts').then(r => r.data) });
  const { data: employees = [] } = useQuery<any[]>({ queryKey: ['dashboard-employees'], queryFn: () => client.get('/admin/api/employees').then(r => r.data) });
  const { data: orderProfit = [] } = useQuery<any[]>({ queryKey: ['dashboard-opr'], queryFn: () => client.get('/admin/api/report/order-profit').then(r => r.data) });
  const { data: carriers = [] } = useQuery<any[]>({ queryKey: ['dashboard-carriers'], queryFn: () => client.get('/admin/api/carriers').then(r => r.data) });
  const { data: couriers = [] } = useQuery<any[]>({ queryKey: ['dashboard-couriers'], queryFn: () => client.get('/admin/api/couriers').then(r => r.data) });

  const totalRevenue = orderProfit.reduce((s: number, r: any) => s + (r.revenue || 0), 0);
  const todayOrders = orders.filter((o: any) => o.status === 'in_transit' || o.status === 'pending_picking').length;
  const pendingParcels = parcels.filter((p: any) => p.status === 'stored' || p.status === 'received').length;

  const cards = [
    { icon: ShoppingCart, label: '总订单', value: orders.length, color: 'blue' },
    { icon: Package, label: '包裹数', value: parcels.length, color: 'indigo' },
    { icon: Users, label: '客户数', value: clients.length, color: 'green' },
    { icon: Warehouse, label: '待处理包裹', value: pendingParcels, color: 'amber' },
    { icon: Truck, label: '进行中订单', value: todayOrders, color: 'purple' },
    { icon: DollarSign, label: '本月营收', value: `¥${totalRevenue.toLocaleString()}`, color: 'teal' },
    { icon: TrendingUp, label: '承运商', value: carriers.length, color: 'rose' },
    { icon: AlertTriangle, label: '快递渠道', value: couriers.length, color: 'red' },
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
