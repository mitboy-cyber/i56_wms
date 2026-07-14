import { Outlet, Link, useLocation } from 'react-router-dom';
import { Box, PackageSearch, Scale, ArrowDownToLine, ShoppingCart, Package, Truck, AlertCircle, Search, Home } from 'lucide-react';
const items = [
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
  return (
    <div className="min-h-screen bg-gray-50">
      <main className="max-w-lg mx-auto p-4"><Outlet /></main>
      <nav className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 flex justify-around py-2">
        {items.map(i => (
          <Link key={i.to} to={i.to} className={`flex flex-col items-center text-xs px-1 ${loc.pathname.startsWith(i.to)?'text-blue-600':'text-gray-500'}`}>
            <i.icon size={20} /><span>{i.label}</span>
          </Link>
        ))}
      </nav>
    </div>
  );
}
