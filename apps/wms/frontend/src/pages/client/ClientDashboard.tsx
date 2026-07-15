import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import { PackageSearch, ShoppingCart, Wallet, Clock, CheckCircle, ArrowRightLeft } from 'lucide-react';

export default function ClientDashboard() {
  const { data } = useQuery({ queryKey: ['client-dashboard'], queryFn: () => clientApi.dashboard() });
  const d = data?.data || {};

  const cards = [
    { icon: ShoppingCart, label: '集运订单', value: d?.order_count ?? '-', color: 'bg-blue-50 border-blue-200 text-blue-700' },
    { icon: PackageSearch, label: '入库包裹', value: d?.total_parcels ?? '-', color: 'bg-green-50 border-green-200 text-green-700' },
    { icon: Wallet, label: '账户余额', value: d?.balance !== undefined ? `¥${Number(d.balance).toFixed(2)}` : '-', color: 'bg-orange-50 border-orange-200 text-orange-700' },
    { icon: ArrowRightLeft, label: '本月消费', value: d?.monthly_spent !== undefined ? `¥${Number(d.monthly_spent).toFixed(2)}` : '-', color: 'bg-purple-50 border-purple-200 text-purple-700' },
    { icon: Clock, label: '进行中订单', value: d?.active_orders ?? '-', color: 'bg-yellow-50 border-yellow-200 text-yellow-700' },
    { icon: CheckCircle, label: '可用线路', value: d?.route_count ?? '-', color: 'bg-teal-50 border-teal-200 text-teal-700' },
  ];

  return (
    <div className="space-y-6">
      <h2 className="text-xl font-bold text-gray-800">仪表盘</h2>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {cards.map((c) => (
          <div key={c.label} className={`rounded-lg border p-4 ${c.color}`}>
            <div className="flex items-center gap-2 mb-1">
              <c.icon size={20} className="opacity-60" />
              <span className="text-sm font-medium opacity-70">{c.label}</span>
            </div>
            <div className="text-2xl font-bold">{c.value}</div>
          </div>
        ))}
      </div>

      {/* Quick actions */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        {[
          { label: '创建订单', href: '/client/orders/new', bg: 'bg-blue-600 hover:bg-blue-700' },
          { label: '预申报包裹', href: '/client/parcels', bg: 'bg-green-600 hover:bg-green-700' },
          { label: '管理申报人', href: '/client/declarants', bg: 'bg-purple-600 hover:bg-purple-700' },
          { label: '充值', href: '/client/ledger', bg: 'bg-orange-600 hover:bg-orange-700' },
        ].map(a => (
          <a key={a.label} href={a.href} className={`${a.bg} text-white rounded-lg px-4 py-3 text-sm font-medium text-center transition-colors`}>
            {a.label}
          </a>
        ))}
      </div>
    </div>
  );
}
