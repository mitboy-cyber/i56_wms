import { useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import { TrendingUp, DollarSign, Percent, ShoppingCart } from 'lucide-react';

interface ProfitRow { period: string; orders: number; revenue: number; cost: number; profit: number; margin: number; }

export default function ServiceProfitReport() {
  const { data: rows = [] } = useQuery<ProfitRow[]>({
    queryKey: ['admin-service-profit-report'],
    queryFn: () => client.get('/admin/api/report/service-profit').then(r => r.data || []),
  });
  const totals = rows.reduce((a,r) => ({ orders: a.orders+r.orders, revenue: a.revenue+r.revenue, cost: a.cost+r.cost, profit: a.profit+r.profit }), { orders:0, revenue:0, cost:0, profit:0 });
  const margin = totals.revenue > 0 ? (totals.profit/totals.revenue*100).toFixed(2) : '0.00';

  return (
    <div className="space-y-6">
      <h2 className="text-xl font-bold text-gray-800">附加服务盈利报表</h2>
      <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
        <Kpi icon={ShoppingCart} label="服务单数" value={`${totals.orders}单`} color="blue" />
        <Kpi icon={DollarSign} label="销售额" value={`¥${totals.revenue.toFixed(2)}`} color="green" />
        <Kpi icon={DollarSign} label="成本" value={`¥${totals.cost.toFixed(2)}`} color="rose" />
        <Kpi icon={TrendingUp} label="毛利" value={`¥${totals.profit.toFixed(2)}`} color="amber" />
        <Kpi icon={Percent} label="毛利率" value={`${margin}%`} color="purple" />
      </div>
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b"><tr>
            <th className="px-4 py-3 text-left font-medium text-gray-600">周期</th>
            <th className="px-4 py-3 text-right font-medium text-gray-600">服务单数</th>
            <th className="px-4 py-3 text-right font-medium text-gray-600">销售额</th>
            <th className="px-4 py-3 text-right font-medium text-gray-600">成本</th>
            <th className="px-4 py-3 text-right font-medium text-gray-600">毛利</th>
            <th className="px-4 py-3 text-right font-medium text-gray-600">毛利率</th>
          </tr></thead>
          <tbody className="divide-y">
            {rows.map((r, i) => (
              <tr key={i} className="hover:bg-gray-50">
                <td className="px-4 py-3">{r.period}</td>
                <td className="px-4 py-3 text-right">{r.orders}</td>
                <td className="px-4 py-3 text-right">¥{r.revenue.toFixed(2)}</td>
                <td className="px-4 py-3 text-right">¥{r.cost.toFixed(2)}</td>
                <td className="px-4 py-3 text-right text-green-600 font-medium">¥{r.profit.toFixed(2)}</td>
                <td className="px-4 py-3 text-right">{(r.margin*100).toFixed(2)}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function Kpi({ icon: Icon, label, value, color }: any) {
  const colors: Record<string,string> = { blue:'bg-blue-50 text-blue-700 border-blue-200', green:'bg-green-50 text-green-700 border-green-200', rose:'bg-rose-50 text-rose-700 border-rose-200', amber:'bg-amber-50 text-amber-700 border-amber-200', purple:'bg-purple-50 text-purple-700 border-purple-200' };
  return <div className={`rounded-lg border p-4 ${colors[color]}`}><div className="flex items-center gap-2 mb-1"><Icon size={18}/><span className="text-xs font-medium opacity-70">{label}</span></div><div className="text-xl font-bold">{value}</div></div>;
}
