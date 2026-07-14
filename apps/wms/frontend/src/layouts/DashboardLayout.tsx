import { NavLink, Outlet, useNavigate } from "react-router-dom"
import { useAuthStore } from "@/stores/auth"
import {
  LayoutDashboard, Package, Users, Warehouse, Truck, Shield,
  UserCog, Settings, LogOut, ChevronDown, ChevronRight, Box, Ship, Plane,
} from "lucide-react"
import { useState } from "react"

interface MenuGroup {
  label: string
  icon: React.ElementType
  children: { label: string; href: string }[]
}

const menu: MenuGroup[] = [
  {
    label: "首页", icon: LayoutDashboard,
    children: [
      { label: "仪表盘", href: "/admin/dashboard" },
      { label: "仓库看板", href: "/admin/warehouse-dashboard" },
    ],
  },
  {
    label: "订单管理", icon: Package,
    children: [
      { label: "集运订单", href: "/admin/orders" },
      { label: "附加服务订单", href: "/admin/service-orders" },
    ],
  },
  {
    label: "仓库管理", icon: Warehouse,
    children: [
      { label: "仓库列表", href: "/admin/warehouses" },
      { label: "包裹列表", href: "/admin/parcels" },
      { label: "工单列表", href: "/admin/work-orders" },
    ],
  },
  {
    label: "运输管理", icon: Truck,
    children: [
      { label: "承运商", href: "/admin/carriers" },
      { label: "快递公司", href: "/admin/couriers" },
      { label: "运输方式", href: "/admin/transport-modes" },
    ],
  },
  {
    label: "客户管理", icon: Users,
    children: [
      { label: "客户列表", href: "/admin/clients" },
      { label: "客户账号", href: "/admin/client-accounts" },
      { label: "申报人", href: "/admin/declarants" },
    ],
  },
  {
    label: "系统管理", icon: Settings,
    children: [
      { label: "员工管理", href: "/admin/employees" },
      { label: "角色管理", href: "/admin/roles" },
      { label: "系统参数", href: "/admin/settings" },
    ],
  },
]

export function DashboardLayout() {
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()
  const [expanded, setExpanded] = useState<Record<string, boolean>>({})

  const toggle = (label: string) => setExpanded((p) => ({ ...p, [label]: !p[label] }))

  const handleLogout = async () => {
    await logout()
    navigate("/admin/login")
  }

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <aside className="w-64 bg-slate-900 text-white flex flex-col shrink-0">
        <div className="px-5 py-4 border-b border-slate-700">
          <h1 className="text-lg font-bold tracking-tight">I56 Framework</h1>
          <p className="text-xs text-slate-400">Admin Console</p>
        </div>
        <nav className="flex-1 overflow-y-auto px-3 py-4 space-y-1">
          {menu.map((g) => (
            <div key={g.label}>
              <button
                onClick={() => toggle(g.label)}
                className="w-full flex items-center gap-2 px-3 py-2 text-sm text-slate-300 hover:text-white hover:bg-slate-800 rounded-md transition-colors"
              >
                <g.icon size={16} />
                <span className="flex-1 text-left">{g.label}</span>
                {expanded[g.label] ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
              </button>
              {expanded[g.label] && (
                <div className="ml-4 mt-1 space-y-0.5">
                  {g.children.map((c) => (
                    <NavLink
                      key={c.href}
                      to={c.href}
                      className={({ isActive }) =>
                        `block px-3 py-1.5 text-sm rounded-md transition-colors ${
                          isActive ? "bg-blue-600 text-white" : "text-slate-400 hover:text-white hover:bg-slate-800"
                        }`
                      }
                    >
                      {c.label}
                    </NavLink>
                  ))}
                </div>
              )}
            </div>
          ))}
        </nav>
        <div className="px-5 py-3 border-t border-slate-700 flex items-center gap-3">
          <Shield size={16} className="text-slate-400" />
          <span className="text-sm text-slate-300 flex-1 truncate">{user?.real_name || user?.username}</span>
          <button onClick={handleLogout} className="text-slate-400 hover:text-white"><LogOut size={16} /></button>
        </div>
      </aside>

      {/* Main */}
      <div className="flex-1 flex flex-col min-w-0">
        <header className="h-14 bg-white border-b flex items-center px-6 shrink-0">
          <h2 className="text-sm font-medium text-gray-700">I56 WMS 管理后台</h2>
        </header>
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
