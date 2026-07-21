import { useState } from "react"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import client from "@/api/client"

interface Column {
  key: string
  label: string
  render?: (value: unknown, row: Record<string, unknown>) => React.ReactNode
}

interface Field {
  key: string
  label: string
  type?: "text" | "number" | "select"
  options?: { label: string; value: string }[]
}

interface Props {
  title: string
  queryKey: string[]
  queryFn: () => Promise<any>
  apiBase: string
  columns: Column[]
  fields?: Field[]
  searchable?: boolean
  getRowId?: (row: any, index: number) => string
  enableCreate?: boolean
  enableEdit?: boolean
  enableDelete?: boolean
  enableSearch?: boolean
  enableExport?: boolean
  filters?: any
  onFilterChange?: any
  activeFilters?: any
  renderActions?: any
}

// CSV export helper
function exportCSV(filename: string, columns: Column[], data: any[]) {
  const header = columns.map(c => c.label).join(",")
  const rows = data.map(row => columns.map(c => {
    const v = String(row[c.key] ?? "").replace(/"/g, '""').replace(/\n/g, " ")
    return `"${v}"`
  }).join(","))
  const csv = [header, ...rows].join("\n")
  const blob = new Blob(["\uFEFF" + csv], { type: "text/csv;charset=utf-8" })
  const url = URL.createObjectURL(blob)
  const a = document.createElement("a"); a.href = url; a.download = filename; a.click()
  URL.revokeObjectURL(url)
}

export default function MinimalListPage({ title, queryKey, queryFn, apiBase, columns, fields = [], searchable = true, enableExport = true }: Props) {
  const qc = useQueryClient()
  const [search, setSearch] = useState("")
  const [modalOpen, setModalOpen] = useState(false)
  const [editItem, setEditItem] = useState<Record<string, unknown> | null>(null)
  const [formData, setFormData] = useState<Record<string, unknown>>({})
  const [page, setPage] = useState(1)
  const perPage = 20

  const { data: rawData = [], isLoading, error } = useQuery<any[]>({ queryKey: [...queryKey, apiBase], queryFn: async () => {
    const r = await queryFn()
    return Array.isArray(r?.data) ? r.data : (Array.isArray(r) ? r : [])
  }, retry: false })

  const filtered = search ? (Array.isArray(rawData) ? rawData.filter(row =>
    columns.some(c => String(row[c.key] ?? "").toLowerCase().includes(search.toLowerCase()))
  ) : []) : (Array.isArray(rawData) ? rawData : [])

  const totalPages = Math.ceil(filtered.length / perPage)
  const paged = filtered.slice((page - 1) * perPage, page * perPage)

  const createMut = useMutation({
    mutationFn: (body: Record<string, unknown>) => client.post(apiBase, body),
    onSuccess: () => { qc.invalidateQueries({ queryKey: [...queryKey, apiBase] }); closeModal() },
  })
  const updateMut = useMutation({
    mutationFn: ({ id, body }: { id: number; body: Record<string, unknown> }) => client.put(`${apiBase}/${id}`, body),
    onSuccess: () => { qc.invalidateQueries({ queryKey: [...queryKey, apiBase] }); closeModal() },
  })
  const deleteMut = useMutation({
    mutationFn: (id: number) => client.delete(`${apiBase}/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: [...queryKey, apiBase] }),
  })

  function closeModal() { setModalOpen(false); setEditItem(null); setFormData({}) }
  function openCreate() { setEditItem(null); setFormData({}); setModalOpen(true) }
  function openEdit(row: Record<string, unknown>) { setEditItem(row); setFormData({...row}); setModalOpen(true) }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (editItem?.id != null) {
      updateMut.mutate({ id: editItem.id as number, body: formData })
    } else {
      createMut.mutate(formData)
    }
  }

  const thStyle: React.CSSProperties = { padding: "8px 12px", textAlign: "left", borderBottom: "1px solid #e5e7eb", color: "#6b7280", fontWeight: 600, fontSize: 13 }
  const tdStyle: React.CSSProperties = { padding: "8px 12px", borderBottom: "1px solid #f3f4f6", fontSize: 14 }
  const btnStyle: React.CSSProperties = { padding: "4px 12px", borderRadius: 6, border: "1px solid #d1d5db", background: "white", cursor: "pointer", fontSize: 13, marginRight: 4 }

  return (
    <div>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
        <h1 style={{ fontSize: 20, fontWeight: "bold" }}>{title}</h1>
        <div style={{ display: "flex", gap: 8 }}>
          {enableExport && filtered.length > 0 && (
            <button onClick={() => exportCSV(`${title}-${new Date().toISOString().slice(0,10)}.csv`, columns, filtered)}
              style={{ padding: "8px 16px", borderRadius: 6, border: "1px solid #6366f1", background: "white", color: "#6366f1", cursor: "pointer", fontSize: 14 }}>
              📥 导出 CSV
            </button>
          )}
          <button onClick={openCreate} style={{ background: "#10b981", color: "white", border: "none", borderRadius: 6, padding: "8px 16px", cursor: "pointer", fontSize: 14 }}>
            + 添加
          </button>
        </div>
      </div>

      {searchable && (
        <div style={{ marginBottom: 12 }}>
          <input type="text" placeholder="搜索..." value={search} onChange={e => { setSearch(e.target.value); setPage(1) }}
            style={{ padding: "8px 12px", border: "1px solid #d1d5db", borderRadius: 6, width: 240, fontSize: 14, outline: "none" }} />
        </div>
      )}

      {isLoading && <p style={{ color: "#9ca3af" }}>加载中...</p>}
      {error && <p style={{ color: "#ef4444" }}>数据加载失败，请刷新重试</p>}
      {!isLoading && filtered.length === 0 && <p style={{ color: "#9ca3af" }}>暂无数据</p>}

      {filtered.length > 0 && (
        <>
          <div style={{ background: "white", borderRadius: 8, border: "1px solid #e5e7eb", overflow: "auto" }}>
            <table style={{ width: "100%", borderCollapse: "collapse" }}>
              <thead>
                <tr>
                  <th style={{ ...thStyle, width: 40 }}>#</th>
                  {columns.map(c => <th key={c.key} style={thStyle}>{c.label}</th>)}
                  <th style={{ ...thStyle, width: 120 }}>操作</th>
                </tr>
              </thead>
              <tbody>
                {paged.map((row: any, idx: number) => (
                  <tr key={row.id ?? idx} style={{ background: idx % 2 === 0 ? "#fafafa" : "white" }}>
                    <td style={tdStyle}>{(page - 1) * perPage + idx + 1}</td>
                    {columns.map(c => (
                      <td key={c.key} style={tdStyle}>
                        {c.render ? c.render(row[c.key], row) : (row[c.key] != null ? String(row[c.key]) : "-")}
                      </td>
                    ))}
                    <td style={tdStyle}>
                      <button onClick={() => openEdit(row)} style={btnStyle}>编辑</button>
                      <button onClick={() => { if (confirm("确认删除?")) deleteMut.mutate(row.id) }} style={{...btnStyle, color: "#ef4444", borderColor: "#fca5a5"}}>删除</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {totalPages > 1 && (
            <div style={{ display: "flex", justifyContent: "center", gap: 8, marginTop: 12 }}>
              <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page <= 1} style={{...btnStyle, opacity: page <= 1 ? .4 : 1}}>◀</button>
              <span style={{ padding: "4px 8px", fontSize: 14, color: "#6b7280" }}>{page} / {totalPages}</span>
              <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page >= totalPages} style={{...btnStyle, opacity: page >= totalPages ? .4 : 1}}>▶</button>
            </div>
          )}
        </>
      )}

      {modalOpen && (
        <div style={{ position: "fixed", inset: 0, background: "rgba(0,0,0,0.4)", display: "flex", alignItems: "center", justifyContent: "center", zIndex: 100 }}>
          <div style={{ background: "white", borderRadius: 12, padding: 24, width: 480, maxHeight: "80vh", overflow: "auto" }}>
            <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
              <h2 style={{ fontSize: 18, fontWeight: "bold" }}>{editItem ? "编辑" : "新增"} {title}</h2>
              <button onClick={closeModal} style={{ background: "none", border: "none", fontSize: 20, cursor: "pointer", color: "#9ca3af" }}>✕</button>
            </div>
            <form onSubmit={handleSubmit}>
              {fields.map(f => (
                <div key={f.key} style={{ marginBottom: 12 }}>
                  <label style={{ display: "block", fontSize: 13, marginBottom: 4, color: "#374151" }}>{f.label}</label>
                  {f.type === "select" && f.options ? (
                    <select value={String(formData[f.key] ?? "")} onChange={e => setFormData(prev => ({ ...prev, [f.key]: e.target.value }))}
                      style={{ width: "100%", padding: "8px 12px", border: "1px solid #d1d5db", borderRadius: 6, fontSize: 14 }}>
                      <option value="">请选择</option>
                      {f.options.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
                    </select>
                  ) : (
                    <input type={f.type === "number" ? "number" : "text"} value={String(formData[f.key] ?? "")} onChange={e => setFormData(prev => ({ ...prev, [f.key]: f.type === "number" ? Number(e.target.value) : e.target.value }))}
                      style={{ width: "100%", padding: "8px 12px", border: "1px solid #d1d5db", borderRadius: 6, fontSize: 14, boxSizing: "border-box" }} />
                  )}
                </div>
              ))}
              <div style={{ display: "flex", gap: 8, justifyContent: "flex-end", marginTop: 16 }}>
                <button type="button" onClick={closeModal} style={{ ...btnStyle, padding: "8px 16px" }}>取消</button>
                <button type="submit" style={{ background: "#10b981", color: "white", border: "none", borderRadius: 6, padding: "8px 16px", cursor: "pointer", fontSize: 14 }}>保存</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
