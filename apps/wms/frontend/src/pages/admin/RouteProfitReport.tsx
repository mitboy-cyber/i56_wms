import { useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import { TrendingUp, DollarSign, Percent, Route } from 'lucide-react';

interface RouteRow { route: string; orders: number; revenue: number; cost: number; profit: number; margin: number; }

export default function RouteProfitReport() {
  const { data: rows = [] } = useQuery<RouteRow[]>({
    queryKey: ['admin-route-profit-report'],
    queryFn: () => client.get('/admin/api/report/route-profit').then(r => r.data || []),
  });
  const totals = rows.reduce((a,r) => ({ orders: a.orders+r.orders, rev: a.rev+r.revenue, cost: a.cost+r.cost, profit: a.profit+r.profit }), { orders:0, rev:0, cost:0, profit:0 });
  const margin = totals.rev > 0 ? (totals.profit/totals.rev*100).toFixed(2) : '0.00';

  return (
    <div className="space-y-6">
      <h2 className="text-xl font-bold text-gray-800">路线盈利汇总</h2>
      <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
        <Kpi icon={Route} label="活跃线路" value={`${rows.length}条`} color="blue" />
        <Kpi icon={DollarSign} label="销售总额" value={`¥${totals.rev.toLocaleString()}`} color="green" />
        <Kpi icon={DollarSign} label="成本总额" value={`¥${totals.cost.toLocaleString()}`} color="rose" />
        <Kpi icon={TrendingUp} label="毛利总额" value={`¥${totals.profit.toLocaleString()}`} color="amber" />
        <Kpi icon={Percent} label="毛利率" value={`${margin}%`} color="purple" />
      </div>
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b"><tr>
            <th className="px-4 py-3 text-left">线路</th><th className="px-4 py-3 text-right">订单数</th>
            <th className="px-4 py-3 text-right">销售额</th><th className="px-4 py-3 text-right">成本</th><th className="px-4 py-3 text-right">毛利</th><th className="px-4 py-3 text-right">毛利率</th>
          </tr></thead>
          <tbody className="divide-y">
            {rows.map((r, i) => (
              <tr key={i} className="hover:bg-gray-50"><td className="px-4 py-3 font-medium">{r.route}</td><td className="px-4 py-3 text-right">{r.orders}</td>
                <td className="px-4 py-3 text-right">¥{r.revenue.toLocaleString()}</td><td className="px-4 py-3 text-right">¥{r.cost.toLocaleString()}</td>
                <td className="px-4 py-3 text-right text-green-600 font-medium">¥{r.profit.toLocaleString()}</td><td className="px-4 py-3 text-right">{(r.margin*100).toFixed(2)}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
function Kpi({ icon: Icon, label, value, color }: any) {
  const c: Record<string,string> = { blue:'bg-blue-50 text-blue-700 border-blue-200', green:'bg-green-50 text-green-700 border-green-200', rose:'bg-rose-50 text-rose-700 border-rose-200', amber:'bg-amber-50 text-amber-700 border-amber-200', purple:'bg-purple-50 text-purple-700 border-purple-200' };
  return <div className={`rounded-lg border p-4 ${c[color]}`}><div className="flex items-center gap-2 mb-1"><Icon size={18}/><span className="text-xs font-medium opacity-70">{label}</span></div><div className="text-xl font-bold">{value}</div></div>;
}
