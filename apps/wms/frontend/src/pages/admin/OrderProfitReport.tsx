import { useEffect, useRef } from 'react';
import { useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import { BarChart3, TrendingUp, DollarSign, ShoppingCart, Percent, TrendingDown } from 'lucide-react';

interface ProfitRow {
  period: string; orders: number; revenue: number; cost: number; profit: number; margin: number;
}

export default function OrderProfitReport() {
  const { data: rows = [] } = useQuery<ProfitRow[]>({
    queryKey: ['admin-order-profit-report'],
    queryFn: () => client.get('/admin/api/report/order-profit').then(r => r.data || []),
  });

  const totals = rows.reduce((acc, r) => ({
    orders: acc.orders + r.orders,
    revenue: acc.revenue + r.revenue,
    cost: acc.cost + r.cost,
    profit: acc.profit + r.profit,
  }), { orders: 0, revenue: 0, cost: 0, profit: 0 });

  const marginPct = totals.revenue > 0 ? (totals.profit / totals.revenue * 100).toFixed(2) : '0.00';

  return (
    <div className="space-y-6">
      <h2 className="text-xl font-bold text-gray-800">集运订单盈利报表</h2>

      {/* ── KPI 卡片 (BFT56-aligned) ── */}
      <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
        <KpiCard icon={ShoppingCart} label="订单数" value={`${totals.orders}单`} color="blue" />
        <KpiCard icon={DollarSign} label="销售额" value={`¥${totals.revenue.toLocaleString()}`} color="green" />
        <KpiCard icon={TrendingDown} label="成本" value={`¥${totals.cost.toLocaleString()}`} color="rose" />
        <KpiCard icon={TrendingUp} label="毛利" value={`¥${totals.profit.toLocaleString()}`} color="amber" />
        <KpiCard icon={Percent} label="毛利率" value={`${marginPct}%`} color="purple" />
      </div>

      {/* ── 趋势图 (简易SVG) ── */}
      <MiniChart rows={rows} />

      {/* ── 数据表格 ── */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-600">周期</th>
              <th className="px-4 py-3 text-right font-medium text-gray-600">订单数</th>
              <th className="px-4 py-3 text-right font-medium text-gray-600">销售额</th>
              <th className="px-4 py-3 text-right font-medium text-gray-600">成本</th>
              <th className="px-4 py-3 text-right font-medium text-gray-600">毛利</th>
              <th className="px-4 py-3 text-right font-medium text-gray-600">毛利率</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {rows.map((r, i) => (
              <tr key={i} className="hover:bg-gray-50">
                <td className="px-4 py-3">{r.period}</td>
                <td className="px-4 py-3 text-right">{r.orders}</td>
                <td className="px-4 py-3 text-right">¥{r.revenue.toLocaleString()}</td>
                <td className="px-4 py-3 text-right">¥{r.cost.toLocaleString()}</td>
                <td className="px-4 py-3 text-right text-green-600 font-medium">¥{r.profit.toLocaleString()}</td>
                <td className="px-4 py-3 text-right">{(r.margin * 100).toFixed(2)}%</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function KpiCard({ icon: Icon, label, value, color }: { icon: any; label: string; value: string; color: string }) {
  const colors: Record<string, string> = {
    blue: 'bg-blue-50 text-blue-700 border-blue-200',
    green: 'bg-green-50 text-green-700 border-green-200',
    rose: 'bg-rose-50 text-rose-700 border-rose-200',
    amber: 'bg-amber-50 text-amber-700 border-amber-200',
    purple: 'bg-purple-50 text-purple-700 border-purple-200',
  };
  return (
    <div className={`rounded-lg border p-4 ${colors[color] || ''}`}>
      <div className="flex items-center gap-2 mb-1">
        <Icon size={18} />
        <span className="text-xs font-medium opacity-70">{label}</span>
      </div>
      <div className="text-xl font-bold">{value}</div>
    </div>
  );
}

function MiniChart({ rows }: { rows: ProfitRow[] }) {
  if (rows.length === 0) return null;
  const maxVal = Math.max(...rows.map(r => Math.max(r.revenue, r.cost)));
  const h = 120; const w = 600;
  const pad = 40; const bw = (w - pad * 2) / (rows.length || 1);

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h3 className="text-sm font-medium text-gray-600 mb-3">销售额 vs 成本 趋势</h3>
      <svg viewBox={`0 0 ${w} ${h + 20}`} className="w-full h-36">
        {/* Revenue bars */}
        {rows.map((r, i) => {
          const bh = (r.revenue / maxVal) * h;
          return <rect key={`rev-${i}`} x={pad + i * bw + 2} y={h - bh} width={bw / 3} height={bh} fill="#22c55e" rx="2" opacity="0.8" />;
        })}
        {/* Cost bars */}
        {rows.map((r, i) => {
          const bh = (r.cost / maxVal) * h;
          return <rect key={`cost-${i}`} x={pad + i * bw + 2 + bw / 3} y={h - bh} width={bw / 3} height={bh} fill="#f43f5e" rx="2" opacity="0.8" />;
        })}
        {/* Labels */}
        {rows.map((r, i) => (
          <text key={`lbl-${i}`} x={pad + i * bw + bw / 3} y={h + 16} textAnchor="middle" className="text-[8px]" fill="#9ca3af">{r.period.slice(5)}</text>
        ))}
        <line x1={pad} y1={h} x2={w - pad} y2={h} stroke="#e5e7eb" />
      </svg>
      <div className="flex justify-center gap-4 mt-2 text-xs text-gray-500">
        <span className="flex items-center gap-1"><span className="w-3 h-3 bg-green-500 rounded inline-block" /> 销售额</span>
        <span className="flex items-center gap-1"><span className="w-3 h-3 bg-rose-500 rounded inline-block" /> 成本</span>
      </div>
    </div>
  );
}
