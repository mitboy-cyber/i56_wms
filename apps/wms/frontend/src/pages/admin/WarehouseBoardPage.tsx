import { useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import { Package, ShoppingCart, Truck, Boxes, Container, AlertTriangle, Clock, TrendingUp } from 'lucide-react';

export default function WarehouseBoardPage() {
  const { data: stats } = useQuery<any>({ queryKey: ['wh-board-stats'], queryFn: () => client.get('/admin/api/dashboard/stats') });
  const { data: orders = [] } = useQuery<any[]>({ queryKey: ['wh-board-orders'], queryFn: () => client.get('/admin/api/orders') });
  const { data: parcels = [] } = useQuery<any[]>({ queryKey: ['wh-board-parcels'], queryFn: () => client.get('/admin/api/parcels') });

  const inStock = parcels.filter((p: any) => p.status === 'stored').length;
  const picking = orders.filter((o: any) => o.status === 'pending_picking').length;
  const packing = orders.filter((o: any) => o.status === 'pending_packing').length;
  const loading = orders.filter((o: any) => o.status === 'pending_loading' || o.status === 'loaded').length;

  const cards = [
    { icon: Boxes, label: '已上架', value: inStock, color: 'bg-green-50 border-green-200 text-green-700' },
    { icon: Clock, label: '待拣货', value: picking, color: 'bg-yellow-50 border-yellow-200 text-yellow-700' },
    { icon: Package, label: '待打包', value: packing, color: 'bg-orange-50 border-orange-200 text-orange-700' },
    { icon: Container, label: '装箱中', value: loading, color: 'bg-purple-50 border-purple-200 text-purple-700' },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-gray-800">仓库看板</h1>
        <span className="text-sm text-gray-500">厦门仓 | {new Date().toLocaleTimeString('zh-CN')}</span>
      </div>

      {/* ── Warehouse KPIs ── */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {cards.map(c => (
          <div key={c.color} className={`rounded-lg border p-4 ${c.color}`}>
            <div className="flex items-center gap-2 mb-1"><c.icon size={18} className="opacity-60" />
              <span className="text-xs font-medium opacity-70">{c.label}</span>
            </div>
            <div className="text-2xl font-bold">{c.value}</div>
          </div>
        ))}
      </div>

      {/* ── Parcel location grid ── */}
      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="text-sm font-semibold text-gray-600 mb-4">包裹仓位分布</h3>
        <div className="grid grid-cols-5 gap-2">
          {(() => {
            const locations = ['A-01', 'A-02', 'A-03', 'B-01', 'B-02', 'B-03', 'C-01', 'C-02', 'D-01', 'D-02'];
            return locations.map((loc, i) => {
              const count = i < inStock ? 1 : 0; // Simulate location assignment
              return (
                <div key={loc} className={`rounded-lg border-2 p-3 text-center ${count > 0 ? 'border-green-300 bg-green-50' : 'border-gray-200 bg-gray-50'}`}>
                  <div className="text-xs text-gray-500">{loc}</div>
                  <div className={`text-lg font-bold ${count > 0 ? 'text-green-600' : 'text-gray-300'}`}>
                    {count > 0 ? '●' : '○'}
                  </div>
                </div>
              );
            });
          })()}
        </div>
      </div>

      {/* ── Recent parcels ── */}
      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="text-sm font-semibold text-gray-600 mb-3">最近入库包裹</h3>
        <table className="w-full text-sm">
          <thead className="border-b"><tr>
            <th className="text-left py-2">快递单号</th><th className="text-left py-2">品名</th>
            <th className="text-right py-2">重量(kg)</th><th className="text-right py-2">状态</th>
          </tr></thead>
          <tbody className="divide-y">
            {parcels.slice(0, 8).map((p: any) => (
              <tr key={p.id} className="hover:bg-gray-50">
                <td className="py-2 font-mono text-xs">{p.tracking_number}</td>
                <td className="py-2">{p.product_name}</td>
                <td className="py-2 text-right">{p.actual_weight}</td>
                <td className="py-2 text-right">
                  <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${p.status === 'stored' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}`}>
                    {String(p.status)}
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
