interface TableProps<T> {
  columns: { key: string; label: string; render?: (row: T) => React.ReactNode }[]
  rows: T[]
  actions?: (row: T) => React.ReactNode
  loading?: boolean
}

export function Table<T>({ columns, rows, actions, loading }: TableProps<T>) {
  if (loading) return <div className="text-center py-12 text-gray-500">加载中...</div>
  if (rows.length === 0) return <div className="text-center py-12 text-gray-500">暂无数据</div>

  return (
    <div className="overflow-x-auto border rounded-lg">
      <table className="w-full text-sm">
        <thead className="bg-gray-50 border-b">
          <tr>
            {columns.map((c) => (
              <th key={c.key} className="px-4 py-3 text-left font-medium text-gray-600">{c.label}</th>
            ))}
            {actions && <th className="px-4 py-3 text-right font-medium text-gray-600">操作</th>}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b last:border-b-0 hover:bg-gray-50">
              {columns.map((c) => (
                <td key={c.key} className="px-4 py-3 text-gray-700">
                  {c.render ? c.render(row) : String((row as Record<string,any>)[c.key] ?? "")}
                </td>
              ))}
              {actions && <td className="px-4 py-3 text-right">{actions(row)}</td>}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
