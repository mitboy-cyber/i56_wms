import { NavLink, Outlet, useNavigate, useLocation } from "react-router-dom"
import { useAuthStore } from "@/stores/auth"
import { useTabStore } from "@/stores/tabs"
import { TabBar } from "@/components/TabBar"
import {
  LayoutDashboard, Package, Users, Warehouse, Truck, Shield,
  Settings, LogOut, ChevronDown, ChevronRight, Box, Ship, Plane,
  Globe, CreditCard, Printer, Bell, BarChart3, Wrench, FileText,
  Database, Bot, Clock, History, Monitor, HardDrive, Radio, AlertTriangle,
  ClipboardList, MapPin, Briefcase, ScrollText, UserCog, ListChecks,
  Activity, Radar, Camera, Workflow,
} from "lucide-react"
import { useState, useEffect } from "react"

interface MenuGroup {
  label: string
  icon: React.ElementType
  children: { label: string; href: string }[]
  defaultOpen?: boolean
}

const menu: MenuGroup[] = [
  {
    label: "首页", icon: LayoutDashboard, defaultOpen: true,
    children: [
      { label: "仪表盘", href: "/admin/dashboard" },
      { label: "仓库看板", href: "/admin/warehouse-board" },
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
    label: "仓库管理", icon: Warehouse, defaultOpen: true,
    children: [
      { label: "包裹列表", href: "/admin/parcels" },
      { label: "附加服务工单", href: "/admin/service-workorders" },
      { label: "附加服务模板", href: "/admin/service-templates" },
      { label: "附加服务类型", href: "/admin/service-types" },
      { label: "PDA 在线会话", href: "/admin/pda-sessions" },
      { label: "集装柜管理", href: "/admin/containers" },
      { label: "仓库列表", href: "/admin/warehouses" },
      { label: "入库看板", href: "/admin/inbound-board" },
      { label: "仓库作业台", href: "/admin/warehouse-console" },
      { label: "员工任务监控", href: "/admin/work-orders" },
      { label: "PDA 工单模板", href: "/admin/pda-workorder-templates" },
      { label: "工单流程管理", href: "/admin/workflow-management" },
      { label: "工单列表", href: "/admin/work-orders" },
      { label: "异常记录", href: "/admin/exception-reports" },
    ],
  },
  {
    label: "财务报表", icon: BarChart3,
    children: [
      { label: "集运订单盈利", href: "/admin/report/order-profit" },
      { label: "附加服务盈利", href: "/admin/report/service-profit" },
      { label: "客户盈利", href: "/admin/report/client-profit" },
      { label: "路线盈利", href: "/admin/report/route-profit" },
    ],
  },
  {
    label: "物流管理", icon: Truck,
    children: [
      { label: "区域组管理", href: "/admin/area-groups" },
      { label: "货物类型", href: "/admin/cargo-types" },
      { label: "承运商列表", href: "/admin/carriers" },
      { label: "快递公司", href: "/admin/couriers" },
      { label: "清关公司", href: "/admin/customs-brokers" },
      { label: "清关点管理", href: "/admin/customs-points" },
      { label: "线路模板", href: "/admin/route-templates" },
      { label: "运输公司", href: "/admin/shipping-providers" },
      { label: "运输方式", href: "/admin/transport-modes" },
      { label: "物流追踪", href: "/admin/logistics-tracking" },
    ],
  },
  {
    label: "客户管理", icon: Users,
    children: [
      { label: "客户收件地址", href: "/admin/customer-addresses" },
      { label: "客户申报人", href: "/admin/customer-declarants" },
      { label: "客户管理", href: "/admin/clients" },
      { label: "客户账号", href: "/admin/client-accounts" },
      { label: "客户会员", href: "/admin/client-members-list" },
      { label: "客户充值", href: "/admin/client-recharges" },
      { label: "余额日志", href: "/admin/balance-logs" },
      { label: "充值记录", href: "/admin/recharge-records" },
      { label: "客户价格", href: "/admin/client-pricing" },
      { label: "月结对账单", href: "/admin/monthly-statements" },
      { label: "客户端权限", href: "/admin/client-panel-perms" },
    ],
  },
  {
    label: "系统", icon: Settings,
    children: [
      { label: "通知管理", href: "/admin/notifications" },
      { label: "打印模板", href: "/admin/print-templates" },
      { label: "角色管理", href: "/admin/roles" },
      { label: "员工管理", href: "/admin/employees" },
      { label: "系统参数", href: "/admin/system/params" },
    ],
  },
  // ── I56 differentiated sections (beyond BFT56) ──
  {
    label: "计费管理", icon: CreditCard,
    children: [
      { label: "线路报价", href: "/admin/pricing/routes" },
      { label: "派送费用", href: "/admin/pricing/delivery" },
      { label: "附加费", href: "/admin/pricing/surcharges" },
      { label: "服务计费", href: "/admin/pricing/services" },
    ],
  },
  {
    label: "API 集成", icon: Globe,
    children: [
      { label: "快递 API", href: "/admin/system/api-couriers" },
      { label: "报关 API", href: "/admin/system/api-customs" },
      { label: "通知 API", href: "/admin/system/api-notifications" },
      { label: "打印 API", href: "/admin/system/api-printers" },
      { label: "存储 API", href: "/admin/system/api-storage" },
      { label: "EZ Way API", href: "/admin/system/api-ezway" },
      { label: "设备网关", href: "/admin/system/api-devices" },
      { label: "报关经纪 API", href: "/admin/system/customs-broker-api" },
      { label: "物流 API", href: "/admin/system/logistics-api" },
    ],
  },
  {
    label: "智能工具", icon: Bot,
    children: [
      { label: "AI 聊天", href: "/admin/system/ai-chat" },
      { label: "AI 设置", href: "/admin/system/ai-settings" },
    ],
  },
  {
    label: "运维管理", icon: Monitor,
    children: [
      { label: "定时任务", href: "/admin/system/scheduler" },
      { label: "审计日志", href: "/admin/system/audit-logs" },
      { label: "系统报表", href: "/admin/system/reports" },
    ],
  },
]

export function DashboardLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const { openTab, setActiveTab, tabs, activeTabId } = useTabStore()
  const [expanded, setExpanded] = useState<Record<string, boolean>>(() => {
    const state: Record<string, boolean> = {}
    menu.forEach(g => { if (g.defaultOpen) state[g.label] = true })
    return state
  })

  // Sync URL changes to tab store
  useEffect(() => {
    // Find menu item matching current path
    for (const g of menu) {
      for (const c of g.children) {
        if (c.href === location.pathname) {
          openTab({ id: c.href, label: c.label, href: c.href })
          return
        }
      }
    }
    // If path doesn't match any menu item, still set active
    if (tabs.some(t => t.href === location.pathname)) {
      setActiveTab(location.pathname)
    }
  }, [location.pathname])

  const toggle = (label: string) => setExpanded((p) => ({ ...p, [label]: !p[label] }))

  const handleLogout = async () => {
    await logout()
    navigate("/admin/login")
  }

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar — Hallmark: light, border-right, emerald accent */}
      <aside className="w-64 flex flex-col shrink-0 border-r border-gray-200" style={{ background: 'var(--sidebar-bg)', color: 'var(--sidebar-fg)' }}>
        <div className="px-5 py-4 border-b" style={{ borderColor: 'var(--sidebar-border)' }}>
          <h1 className="text-lg font-light tracking-tight" style={{ color: 'var(--color-ink)' }}>I56 Framework</h1>
          <p className="text-xs" style={{ color: 'var(--color-neutral)' }}>Admin Console</p>
        </div>
        <nav className="flex-1 overflow-y-auto px-3 py-4 space-y-0.5">
          {menu.map((g) => (
            <div key={g.label}>
              <button
                onClick={() => toggle(g.label)}
                className="w-full flex items-center gap-2 px-3 py-2 text-sm rounded-md transition-colors"
                style={{ color: 'var(--sidebar-fg)' }}
                onMouseEnter={e => (e.currentTarget.style.background = 'var(--sidebar-accent)')}
                onMouseLeave={e => (e.currentTarget.style.background = 'transparent')}
              >
                <g.icon size={16} />
                <span className="flex-1 text-left">{g.label}</span>
                {expanded[g.label] ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
              </button>
              {expanded[g.label] && (
                <div className="ml-5 mt-0.5 space-y-0.5">
                  {g.children.map((c) => (
                    <NavLink
                      key={c.href}
                      to={c.href}
                      className={({ isActive }: { isActive: boolean }) =>
                        `block px-3 py-1.5 text-sm rounded-md transition-colors ${
                          isActive
                            ? 'active-nav'
                            : ''
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
        <div className="px-5 py-3 border-t border-slate-700 flex items-center gap-3 shrink-0">
          <Shield size={16} className="text-slate-400" />
          <span className="text-sm text-slate-300 flex-1 truncate">{user?.real_name || user?.username}</span>
          <button onClick={handleLogout} className="text-slate-400 hover:text-white"><LogOut size={16} /></button>
        </div>
      </aside>

      {/* Main */}
      <div className="flex-1 flex flex-col min-w-0" style={{ background: 'var(--color-paper)' }}>
        <header className="h-14 bg-white border-b flex items-center px-6 shrink-0" style={{ borderColor: 'var(--color-rule)' }}>
          <h2 className="text-sm font-medium" style={{ color: 'var(--color-neutral)' }}>I56 WMS 管理后台</h2>
        </header>
        <TabBar />
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
