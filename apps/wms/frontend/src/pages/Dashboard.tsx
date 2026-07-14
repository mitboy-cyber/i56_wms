export function DashboardPage() {
  return (
    <div>
      <h1 className="text-xl font-bold mb-6">仪表盘</h1>
      <div className="grid grid-cols-4 gap-4">
        {[
          { label: "订单数", value: "1,248" },
          { label: "包裹数", value: "856" },
          { label: "客户数", value: "42" },
          { label: "仓库数", value: "3" },
        ].map((s) => (
          <div key={s.label} className="bg-white rounded-lg border p-6">
            <p className="text-sm text-gray-500">{s.label}</p>
            <p className="text-2xl font-bold mt-1">{s.value}</p>
          </div>
        ))}
      </div>
    </div>
  )
}
