import { useState, useEffect, useCallback } from "react"
import { Table } from "@/components/Table"
import { Modal } from "@/components/Modal"
import { Button } from "@/components/Button"
import { InputField, SelectField } from "@/components/Form"
import { Plus, Pencil, Trash2, Search } from "lucide-react"

interface Field {
  name: string
  label: string
  type?: "text" | "select"
  options?: { value: string; label: string }[]
}

interface CrudConfig<T> {
  title: string
  columns: { key: string; label: string; render?: (row: T) => React.ReactNode }[]
  fields: Field[]
  api: {
    list: () => Promise<T[]>
    create: (data: Record<string, string>) => Promise<void>
    update: (id: number, data: Record<string, string>) => Promise<void>
    delete: (id: number) => Promise<void>
  }
  getRowId: (row: T) => number
}

export function CrudPage<T>({ config }: { config: CrudConfig<T> }) {
  const [rows, setRows] = useState<T[]>([])
  const [loading, setLoading] = useState(true)
  const [modal, setModal] = useState<{ open: boolean; editId?: number; data: Record<string, string> }>({
    open: false,
    data: {},
  })
  const [search, setSearch] = useState("")

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const data = await config.api.list()
      setRows(data)
    } catch (e) {
      console.error(e)
    } finally {
      setLoading(false)
    }
  }, [config.api])

  useEffect(() => { load() }, [load])

  const handleSave = async () => {
    try {
      if (modal.editId) {
        await config.api.update(modal.editId, modal.data)
      } else {
        await config.api.create(modal.data)
      }
      setModal({ open: false, data: {} })
      load()
    } catch (e) {
      console.error(e)
    }
  }

  const handleEdit = (row: T) => {
    const data: Record<string, string> = {}
    config.fields.forEach((f) => {
      data[f.name] = String((row as any)[f.name] ?? "")
    })
    setModal({ open: true, editId: config.getRowId(row), data })
  }

  const handleDelete = async (row: T) => {
    if (!confirm("确认删除？")) return
    try {
      await config.api.delete(config.getRowId(row))
      load()
    } catch (e) {
      console.error(e)
    }
  }

  const filtered = rows.filter((r) =>
    config.columns.some((c) => String((r as Record<string,any>)[c.key] ?? "").toLowerCase().includes(search.toLowerCase()))
  )

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-xl font-bold">{config.title}</h1>
        <Button onClick={() => setModal({ open: true, data: {} })}>
          <Plus size={16} className="inline mr-1" />添加
        </Button>
      </div>

      <div className="mb-4 relative">
        <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          placeholder="搜索..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full max-w-xs pl-9 pr-3 py-2 border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <Table
        columns={config.columns}
        rows={filtered}
        loading={loading}
        actions={(row) => (
          <div className="flex gap-1 justify-end">
            <Button variant="ghost" onClick={() => handleEdit(row)}><Pencil size={14} /></Button>
            <Button variant="ghost" onClick={() => handleDelete(row)}><Trash2 size={14} /></Button>
          </div>
        )}
      />

      <Modal
        open={modal.open}
        onClose={() => setModal({ open: false, data: {} })}
        title={modal.editId ? `编辑${config.title}` : `新增${config.title}`}
        footer={
          <>
            <Button variant="secondary" onClick={() => setModal({ open: false, data: {} })}>取消</Button>
            <Button onClick={handleSave}>保存</Button>
          </>
        }
      >
        {config.fields.map((f) =>
          f.type === "select" && f.options ? (
            <SelectField
              key={f.name}
              label={f.label}
              name={f.name}
              value={modal.data[f.name] || ""}
              options={f.options}
              onChange={(e) => setModal({ ...modal, data: { ...modal.data, [f.name]: e.target.value } })}
            />
          ) : (
            <InputField
              key={f.name}
              label={f.label}
              name={f.name}
              value={modal.data[f.name] || ""}
              onChange={(e) => setModal({ ...modal, data: { ...modal.data, [f.name]: e.target.value } })}
            />
          )
        )}
      </Modal>
    </div>
  )
}
