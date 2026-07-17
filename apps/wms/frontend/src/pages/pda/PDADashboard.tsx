import { useEffect, useState } from 'react';
import { Box, Scale, ArrowDownToLine, ShoppingCart, Package, Truck } from 'lucide-react';

interface PendingCounts {
  receive: number;
  putaway: number;
  weigh: number;
  pick: number;
  pack: number;
  load: number;
}

const tasks = [
  { key: 'receive', label: '待收货', icon: Box, color: 'var(--color-accent)' },
  { key: 'weigh', label: '待称重', icon: Scale, color: 'oklch(60% 0.15 200)' },
  { key: 'putaway', label: '待上架', icon: ArrowDownToLine, color: 'oklch(55% 0.18 80)' },
  { key: 'pick', label: '待拣货', icon: ShoppingCart, color: 'oklch(60% 0.16 280)' },
  { key: 'pack', label: '待打包', icon: Package, color: 'oklch(58% 0.15 40)' },
  { key: 'load', label: '待装车', icon: Truck, color: 'oklch(52% 0.12 240)' },
];

export default function PDADashboard() {
  const [pending, setPending] = useState<PendingCounts>({ receive: 0, putaway: 0, weigh: 0, pick: 0, pack: 0, load: 0 });

  useEffect(() => {
    fetch('/pda/api/dashboard', { credentials: 'include' })
      .then(r => r.json())
      .then(d => {
        if (d?.data?.pending) setPending(d.data.pending);
        else if (d?.pending) setPending(d.pending);
      })
      .catch(() => {});
  }, []);

  return (
    <div>
      <h2 className="text-lg font-bold mb-4" style={{ color: 'var(--color-ink)' }}>作业台</h2>

      <div className="grid grid-cols-3 gap-3">
        {tasks.map(t => (
          <div key={t.key}
            className="bg-white rounded-xl border p-3 text-center shadow-sm cursor-pointer hover:shadow-md transition-shadow"
            style={{ borderColor: 'var(--border)' }}
            onClick={() => window.location.href = `/pda/${t.key}`}>
            <t.icon size={24} className="mx-auto mb-1" style={{ color: t.color }} />
            <div className="text-2xl font-bold" style={{ color: 'var(--color-ink)' }}>{pending[t.key as keyof PendingCounts]}</div>
            <div className="text-xs mt-1" style={{ color: 'var(--color-neutral)' }}>{t.label}</div>
          </div>
        ))}
      </div>

      <div className="mt-4 bg-white rounded-xl border p-4 shadow-sm" style={{ borderColor: 'var(--border)' }}>
        <h3 className="text-sm font-semibold mb-2" style={{ color: 'var(--color-muted)' }}>快捷操作</h3>
        <div className="grid grid-cols-2 gap-2">
          {tasks.map(t => (
            <a key={t.key} href={`/pda/${t.key}`}
              className="flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors border"
              style={{ borderColor: 'var(--border)', color: 'var(--color-muted)' }}>
              <t.icon size={14} style={{ color: t.color }} />
              {t.label}
            </a>
          ))}
        </div>
      </div>
    </div>
  );
}
