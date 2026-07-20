import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom';
import { useClientAuth } from '@/stores/clientAuth';
import { Package, PackageSearch, ShoppingCart, Wallet, Users, MapPin, Warehouse, Truck, Wrench, Settings } from 'lucide-react';

const navItems = [
  { to: '/client/dashboard', label: '仪表盘', icon: Package },
  { to: '/client/parcels', label: '包裹预报', icon: PackageSearch },
  { to: '/client/orders', label: '集运订单', icon: ShoppingCart },
  { to: '/client/ledger', label: '财务明细', icon: Wallet },
  { to: '/client/declarants', label: '申报人', icon: Users },
  { to: '/client/members', label: '会员管理', icon: Users },
  { to: '/client/addresses', label: '收件地址', icon: MapPin },
  { to: '/client/warehouses', label: '仓库', icon: Warehouse },
  { to: '/client/couriers', label: '快递公司', icon: Truck },
  { to: '/client/services', label: '附加服务', icon: Wrench },
  { to: '/client/credentials', label: 'API凭证', icon: Settings },
];

export default function ClientLayout() {
  const { client, logout } = useClientAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/client/login');
  };

  const isActive = (path: string) => location.pathname.startsWith(path);

  return (
    <div className="flex min-h-screen" style={{ background: 'var(--background)' }}>
      <aside className="w-56 flex flex-col" style={{ background: 'var(--sidebar-bg)', borderRight: '1px solid var(--border)' }}>
        <div className="p-4" style={{ borderBottom: '1px solid var(--border)' }}>
          <h1 className="text-lg font-bold" style={{ color: 'var(--color-accent)' }}>I56 Client</h1>
          <p className="text-xs" style={{ color: 'var(--color-neutral)' }}>{client}</p>
        </div>
        <nav className="flex-1 p-2 space-y-1 overflow-y-auto">
          {navItems.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              className="flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors"
              style={{
                background: isActive(item.to) ? 'var(--accent-bg, color-mix(in oklch, var(--color-accent), transparent 90%))' : 'transparent',
                color: isActive(item.to) ? 'var(--color-accent)' : 'var(--color-muted)',
                fontWeight: isActive(item.to) ? 500 : 400,
              }}
            >
              <item.icon className="w-4 h-4" />
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="p-3" style={{ borderTop: '1px solid var(--border)' }}>
          <button
            onClick={handleLogout}
            className="w-full text-left px-3 py-2 text-sm rounded-lg transition-colors"
            style={{ color: 'var(--destructive)' }}
          >
            退出登录
          </button>
        </div>
      </aside>
      <main className="flex-1 p-6 overflow-auto">
        <Outlet />
      </main>
    </div>
  );
}
