import { NavLink, Outlet, useNavigate } from "react-router-dom"
import { useAuthStore } from "@/stores/auth"
import {
  LayoutDashboard, Package, Users, Warehouse, Truck, Shield,
  Settings, LogOut, ChevronDown, ChevronRight, Box, Ship, Plane,
  Globe, CreditCard, Printer, Bell, BarChart3, Wrench, FileText,
  Database, Bot, Clock, History, Monitor, HardDrive, Radio, AlertTriangle,
  ClipboardList, MapPin, Briefcase, ScrollText, UserCog, ListChecks,
  Activity, Radar, Camera, Workflow,
} from "lucide-react"
import { useState } from "react"

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
      { label: "入库看板", href: "/admin/inbound-board" },
      { label: "仓库控制台", href: "/admin/warehouse-console" },
      { label: "任务监控", href: "/admin/task-monitor" },
    ],
  },
  {
    label: "订单管理 OMS", icon: Package, defaultOpen: true,
    children: [
      { label: "集运订单", href: "/admin/orders" },
      { label: "附加服务订单", href: "/admin/service-orders" },
      { label: "服务工单", href: "/admin/service-workorders" },
      { label: "服务模板", href: "/admin/service-templates" },
      { label: "服务类型", href: "/admin/service-types" },
    ],
  },
  {
    label: "仓库管理 WMS", icon: Warehouse,
    children: [
      { label: "仓库列表", href: "/admin/warehouses" },
      { label: "包裹列表", href: "/admin/parcels" },
      { label: "工单管理", href: "/admin/work-orders" },
      { label: "PDA工单模板", href: "/admin/pda-workorder-templates" },
      { label: "PDA会话", href: "/admin/pda-sessions" },
      { label: "异常列表", href: "/admin/exceptions" },
      { label: "异常报告", href: "/admin/exception-reports" },
      { label: "AI异常", href: "/admin/ai-exceptions" },
      { label: "工作流管理", href: "/admin/workflow-management" },
    ],
  },
  {
    label: "运输管理 TMS", icon: Truck,
    children: [
      { label: "线路模板", href: "/admin/route-templates" },
      { label: "承运商", href: "/admin/carriers" },
      { label: "快递公司", href: "/admin/couriers" },
      { label: "运输方式", href: "/admin/transport-modes" },
      { label: "区域分组", href: "/admin/area-groups" },
      { label: "货物类型", href: "/admin/cargo-types" },
      { label: "报关行", href: "/admin/customs-brokers" },
      { label: "海关口岸", href: "/admin/customs-points" },
      { label: "集装箱装货", href: "/admin/container-loadings" },
      { label: "物流追踪", href: "/admin/logistics-tracking" },
      { label: "承运商管理", href: "/admin/shipping-providers" },
    ],
  },
  {
    label: "客户管理 CRM", icon: Users,
    children: [
      { label: "客户列表", href: "/admin/clients" },
      { label: "客户账号", href: "/admin/client-accounts" },
      { label: "客户成员", href: "/admin/client-members" },
      { label: "客户充值", href: "/admin/client-recharges" },
      { label: "客户账本", href: "/admin/client-ledgers" },
      { label: "收支明细", href: "/admin/balance-logs" },
      { label: "客户定价", href: "/admin/client-pricing" },
      { label: "客户权限", href: "/admin/client-permissions" },
      { label: "申报人", href: "/admin/declarants" },
      { label: "客户申报人", href: "/admin/customer-declarants" },
      { label: "收件地址", href: "/admin/customer-addresses" },
      { label: "月度对账单", href: "/admin/monthly-statements" },
    ],
  },
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
    label: "财务分析", icon: BarChart3,
    children: [
      { label: "订单利润", href: "/admin/report/order-profit" },
      { label: "线路利润", href: "/admin/report/route-profit" },
      { label: "客户利润", href: "/admin/report/client-profit" },
      { label: "服务利润", href: "/admin/report/service-profit" },
    ],
  },
  {
    label: "组织权限", icon: Shield,
    children: [
      { label: "员工管理", href: "/admin/employees" },
      { label: "角色管理", href: "/admin/roles" },
    ],
  },
  {
    label: "系统管理", icon: Settings,
    children: [
      { label: "系统参数", href: "/admin/system/params" },
      { label: "系统设置", href: "/admin/system/settings" },
      { label: "品牌设置", href: "/admin/system/brand" },
      { label: "通知管理", href: "/admin/notifications" },
      { label: "通知渠道", href: "/admin/system/notification-channels" },
      { label: "打印模板", href: "/admin/print-templates" },
      { label: "打印机", href: "/admin/printers" },
      { label: "存储配置", href: "/admin/storage" },
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
  const { user, logout } = useAuthStore()
  const [expanded, setExpanded] = useState<Record<string, boolean>>(() => {
    const state: Record<string, boolean> = {}
    menu.forEach(g => { if (g.defaultOpen) state[g.label] = true })
    return state
  })

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
        <nav className="flex-1 overflow-y-auto px-3 py-4 space-y-0.5">
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
                <div className="ml-5 mt-0.5 space-y-0.5">
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
        <div className="px-5 py-3 border-t border-slate-700 flex items-center gap-3 shrink-0">
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
