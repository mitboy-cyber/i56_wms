import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom';
import { Home, Box, Scale, ArrowDownToLine, ShoppingCart, Package, Truck, AlertCircle, Search, LogOut } from 'lucide-react';
import { usePDAAuth } from '@/stores/pdaAuth';

const tabs = [
  { to: '/pda/dashboard', label: '首页', icon: Home },
  { to: '/pda/receive', label: '收货', icon: Box },
  { to: '/pda/weigh', label: '称重', icon: Scale },
  { to: '/pda/putaway', label: '上架', icon: ArrowDownToLine },
  { to: '/pda/pick', label: '拣货', icon: ShoppingCart },
  { to: '/pda/pack', label: '打包', icon: Package },
  { to: '/pda/load', label: '装车', icon: Truck },
  { to: '/pda/exception', label: '异常', icon: AlertCircle },
  { to: '/pda/query', label: '查询', icon: Search },
];

export default function PDALayout() {
  const loc = useLocation();
  const nav = useNavigate();
  const { logout } = usePDAAuth();

  return (
    <div className="min-h-screen pb-16" style={{ background: 'var(--background)' }}>
      {/* Top bar */}
      <header className="sticky top-0 z-10 bg-white border-b shadow-sm px-4 py-2 flex items-center justify-between"
        style={{ borderColor: 'var(--border)' }}>
        <span className="text-sm font-semibold" style={{ color: 'var(--color-accent)' }}>I56 PDA</span>
        <button onClick={() => { logout(); nav('/pda/login'); }}
          className="p-1 rounded-md hover:bg-gray-100 transition-colors"
          style={{ color: 'var(--color-neutral)' }}>
          <LogOut size={18} />
        </button>
      </header>

      {/* Content */}
      <main className="max-w-lg mx-auto px-4 py-4">
        <Outlet />
      </main>

      {/* Bottom tab bar */}
      <nav className="fixed bottom-0 left-0 right-0 bg-white border-t flex justify-around py-1 shadow-lg"
        style={{ borderColor: 'var(--border)' }}>
        {tabs.map(t => {
          const active = loc.pathname.startsWith(t.to);
          return (
            <Link key={t.to} to={t.to}
              className={`flex flex-col items-center text-[10px] px-1 py-1 transition-colors`}
              style={{ color: active ? 'var(--color-accent)' : 'var(--color-neutral)' }}>
              <t.icon size={18} />
              <span className="mt-0.5">{t.label}</span>
            </Link>
          );
        })}
      </nav>
    </div>
  );
}
