import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import { PackageSearch, ShoppingCart, Wallet, TrendingUp } from 'lucide-react';

export default function ClientDashboard() {
  const { data } = useQuery({ queryKey: ['client-dashboard'], queryFn: () => clientApi.dashboard() });
  const d = data?.data;

  const cards = [
    { label: '包裹总数', value: d?.total_parcels ?? '-', icon: PackageSearch, color: 'text-blue-600' },
    { label: '订单数', value: d?.order_count ?? '-', icon: ShoppingCart, color: 'text-green-600' },
    { label: '账户余额', value: d ? `¥${(d.balance as number).toFixed(2)}` : '-', icon: Wallet, color: 'text-orange-600' },
    { label: '线路数', value: d?.route_count ?? '-', icon: TrendingUp, color: 'text-purple-600' },
  ];

  return (
    <div>
      <h2 className="text-xl font-bold text-gray-800 mb-6">仪表盘</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {cards.map((c) => (
          <div key={c.label} className="bg-white rounded-xl shadow-sm border border-gray-200 p-5">
            <div className="flex items-center gap-3">
              <c.icon className={`w-8 h-8 ${c.color}`} />
              <div>
                <p className="text-sm text-gray-500">{c.label}</p>
                <p className="text-2xl font-bold text-gray-900">{c.value}</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
