import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import client from '@/api/client';
import Modal from '@/components/Modal';

interface FieldDef {
  key: string;
  label: string;
  type?: 'text' | 'number';
}

interface ColumnDef {
  key: string;
  label: string;
  render?: (value: any, row: any) => React.ReactNode;
}

interface GenericListPageProps {
  title: string;
  queryKey: string[];
  queryFn: () => Promise<any>;
  mutateUrl?: string;
  columns: ColumnDef[];
  fields?: FieldDef[];
  getRowId: (row: any, i: number) => string;
}

export default function GenericListPage({
  title, queryKey, queryFn, mutateUrl, columns, fields, getRowId,
}: GenericListPageProps) {
  const qc = useQueryClient();
  const url = mutateUrl || (queryFn as any)._url || ''; // try to infer from queryFn
  const result = useQuery({ queryKey, queryFn } as any);
  const rows: any[] = (result as any).data?.data ?? [];
  const isLoading = (result as any).isLoading;

  const [modalOpen, setModalOpen] = useState(false);
  const [editItem, setEditItem] = useState<any>(null);
  const [form, setForm] = useState<Record<string, any>>({});

  const addMut = useMutation({
    mutationFn: (body: any) => client.post(url || '', body),
    onSuccess: () => { qc.invalidateQueries({ queryKey }); setModalOpen(false); },
  });

  const openAdd = () => {
    const init: Record<string, any> = {};
    (fields || columns.map(c => ({ key: c.key, label: c.label, type: 'text' }))).forEach(f => { init[f.key] = ''; });
    setForm(init);
    setEditItem(null);
    setModalOpen(true);
  };

  const openEdit = (row: any) => {
    const init: Record<string, any> = {};
    (fields || columns.map(c => ({ key: c.key, label: c.label, type: 'text' }))).forEach(f => { init[f.key] = row[f.key] ?? ''; });
    setForm(init);
    setEditItem(row);
    setModalOpen(true);
  };

  const handleSave = () => { addMut.mutate(form); };

  const fieldList = fields || columns.map(c => ({ key: c.key, label: c.label, type: 'text' as const }));

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-bold text-gray-800">{title}</h2>
        <button onClick={openAdd} className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors">
          + 添加
        </button>
      </div>

      {isLoading ? (
        <div className="text-center py-8 text-gray-400">加载中...</div>
      ) : rows.length === 0 ? (
        <div className="text-center py-8 text-gray-400">暂无数据</div>
      ) : (
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-gray-50 border-b border-gray-200">
                {columns.map((c) => (
                  <th key={c.key} className="text-left px-4 py-3 font-medium text-gray-600">{c.label}</th>
                ))}
                <th className="text-right px-4 py-3 font-medium text-gray-600 w-32">操作</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row: any, i: number) => (
                <tr key={getRowId(row, i)} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                  {columns.map((c) => (
                    <td key={c.key} className="px-4 py-2.5 text-gray-700">
                      {c.render ? c.render(row[c.key], row) : String(row[c.key] ?? '')}
                    </td>
                  ))}
                  <td className="px-4 py-2.5 text-right">
                    <button onClick={() => openEdit(row)} className="text-blue-600 hover:text-blue-800 mr-3 text-xs font-medium">编辑</button>
                    <button className="text-red-500 hover:text-red-700 text-xs font-medium">删除</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <Modal open={modalOpen} title={editItem ? '编辑' : '添加'} onClose={() => setModalOpen(false)}>
        <div className="space-y-4">
          {fieldList.map(f => (
            <div key={f.key}>
              <label className="block text-sm font-medium text-gray-600 mb-1">{f.label}</label>
              <input
                type={f.type || 'text'}
                value={form[f.key] ?? ''}
                onChange={e => setForm({ ...form, [f.key]: f.type === 'number' ? Number(e.target.value) : e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              />
            </div>
          ))}
          <div className="flex justify-end gap-3 pt-2">
            <button onClick={() => setModalOpen(false)} className="px-4 py-2 text-sm text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200">取消</button>
            <button onClick={handleSave} className="px-4 py-2 text-sm text-white bg-blue-600 rounded-lg hover:bg-blue-700">保存</button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
