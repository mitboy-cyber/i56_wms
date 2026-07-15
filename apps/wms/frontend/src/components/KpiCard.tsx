import type { KpiCardData } from '@/types';
import { TrendingUp, TrendingDown } from 'lucide-react';
import { cn } from '@/lib/utils';

interface KpiCardProps {
  data: KpiCardData;
  onClick?: () => void;
}

const colorMap: Record<string, { bg: string; icon: string; text: string }> = {
  blue: { bg: 'bg-blue-50', icon: 'text-blue-600', text: 'text-blue-700' },
  green: { bg: 'bg-green-50', icon: 'text-green-600', text: 'text-green-700' },
  amber: { bg: 'bg-amber-50', icon: 'text-amber-600', text: 'text-amber-700' },
  red: { bg: 'bg-red-50', icon: 'text-red-600', text: 'text-red-700' },
  purple: { bg: 'bg-purple-50', icon: 'text-purple-600', text: 'text-purple-700' },
};

export function KpiCard({ data, onClick }: KpiCardProps) {
  const palette = colorMap[data.color ?? 'blue'] ?? colorMap.blue;

  return (
    <div
      className={cn(
        'bg-white rounded-xl shadow-sm border border-gray-200 p-5',
        'hover:shadow-md transition-shadow',
        onClick && 'cursor-pointer',
      )}
      onClick={onClick}
    >
      <div className="flex items-start justify-between">
        <div>
          <p className="text-sm font-medium text-gray-500">{data.label}</p>
          <p className={cn('text-2xl font-bold mt-1', palette.text)}>
            {typeof data.value === 'number' ? data.value.toLocaleString() : data.value}
          </p>
          {data.change !== undefined && (
            <div className="flex items-center gap-1 mt-1.5">
              {data.change >= 0 ? (
                <TrendingUp size={14} className="text-green-600" />
              ) : (
                <TrendingDown size={14} className="text-red-600" />
              )}
              <span
                className={cn(
                  'text-xs font-medium',
                  data.change >= 0 ? 'text-green-600' : 'text-red-600',
                )}
              >
                {data.change >= 0 ? '+' : ''}
                {data.change}%
              </span>
              {data.changeLabel && (
                <span className="text-xs text-gray-400 ml-1">{data.changeLabel}</span>
              )}
            </div>
          )}
        </div>
        {data.icon && (
          <div className={cn('p-3 rounded-lg', palette.bg)}>
            <span className={palette.icon}>{data.icon}</span>
          </div>
        )}
      </div>
    </div>
  );
}

interface KpiDashboardProps {
  cards: KpiCardData[];
  columns?: 2 | 3 | 4 | 5;
}

export function KpiDashboard({ cards, columns = 4 }: KpiDashboardProps) {
  const gridCols: Record<number, string> = {
    2: 'grid-cols-1 sm:grid-cols-2',
    3: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3',
    4: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-4',
    5: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-5',
  };

  if (cards.length === 0) {
    return (
      <div className="text-center py-8 text-gray-400">暂无KPI数据</div>
    );
  }

  return (
    <div className={cn('grid gap-4', gridCols[columns] ?? gridCols[4])}>
      {cards.map((card, i) => (
        <KpiCard key={i} data={card} />
      ))}
    </div>
  );
}
