import { createContext, useContext, useState, useCallback, type ReactNode } from "react"

type Lang = "zh" | "en"

const dict: Record<Lang, Record<string, string>> = {
  zh: {
    dashboard: "仪表盘", orders: "订单管理", warehouse: "仓库管理",
    finance: "财务报表", logistics: "物流管理", customers: "客户管理",
    system: "系统", search: "搜索...", logout: "退出",
    add: "+ 添加", edit: "编辑", delete: "删除", save: "保存", cancel: "取消",
    export: "导出 CSV", recharge: "充值", balance: "余额", total: "总计",
  },
  en: {
    dashboard: "Dashboard", orders: "Orders", warehouse: "Warehouse",
    finance: "Finance", logistics: "Logistics", customers: "Customers",
    system: "System", search: "Search...", logout: "Logout",
    add: "+ Add", edit: "Edit", delete: "Delete", save: "Save", cancel: "Cancel",
    export: "Export CSV", recharge: "Recharge", balance: "Balance", total: "Total",
  },
}

const I18nCtx = createContext<{ lang: Lang; t: (key: string) => string; setLang: (l: Lang) => void }>({
  lang: "zh", t: (k) => k, setLang: () => {},
})

export function I18nProvider({ children }: { children: ReactNode }) {
  const [lang, setLang] = useState<Lang>("zh")
  const t = useCallback((key: string) => dict[lang][key] || key, [lang])
  return <I18nCtx.Provider value={{ lang, t, setLang }}>{children}</I18nCtx.Provider>
}

export function useI18n() { return useContext(I18nCtx) }
