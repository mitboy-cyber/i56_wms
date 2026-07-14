import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import clientApi from '@/api/clientApi';
import Modal from '@/components/Modal';

export default function ClientParcels() {
  const qc = useQueryClient();
  const { data } = useQuery({ queryKey: ['client-parcels'], queryFn: () => clientApi.parcels() } as any);
  const [showAdd, setShowAdd] = useState(false);
  const [form, setForm] = useState({ tracking_number: '', product_name: '', courier_code: '' });
  const rows: any[] = (data as any)?.data ?? [];

  const predeclare = useMutation({
    mutationFn: (d: any) => clientApi.predeclare(d),
    onSuccess: () => { (qc as any).invalidateQueries({ queryKey: ['client-parcels'] }); setShowAdd(false); setForm({ tracking_number: '', product_name: '', courier_code: '' }); },
  } as any);

  const columns: any[] = [
    { key: 'tracking_number', label: '快递单号' },
    { key: 'product_name', label: '品名' },
    { key: 'status', label: '状态', render: (v: any) => <span className="px-2 py-0.5 text-xs rounded-full bg-blue-100 text-blue-700">{String(v)}</span> },
    { key: 'actual_weight', label: '重量(kg)', render: (v: any) => v ? Number(v).toFixed(2) : '-' },
  ];

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-bold text-gray-800">包裹预报</h2>
        <button onClick={() => setShowAdd(true)} className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700">新增预报</button>
      </div>
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead><tr className="bg-gray-50 border-b border-gray-200">
            {columns.map(c => <th key={c.key} className="text-left px-4 py-3 font-medium text-gray-600">{c.label}</th>)}
          </tr></thead>
          <tbody>
            {rows.map((row: any, i: number) => (
              <tr key={row.tracking_number || i} className="border-b border-gray-100 hover:bg-gray-50">
                {columns.map(c => (
                  <td key={c.key} className="px-4 py-2.5 text-gray-700">
                    {c.render ? c.render(row[c.key], row) : String(row[c.key] ?? '')}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <Modal open={showAdd} onClose={() => setShowAdd(false)} title="新增包裹预报">
        <form onSubmit={(e) => { e.preventDefault(); (predeclare as any).mutate(form); }} className="space-y-3">
          <div><label className="block text-sm font-medium mb-1">快递单号 *</label>
            <input required value={form.tracking_number} onChange={(e) => setForm({...form, tracking_number: e.target.value})}
              className="w-full px-3 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-blue-500" /></div>
          <div><label className="block text-sm font-medium mb-1">品名</label>
            <input value={form.product_name} onChange={(e) => setForm({...form, product_name: e.target.value})}
              className="w-full px-3 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-blue-500" /></div>
          <div><label className="block text-sm font-medium mb-1">快递公司代码</label>
            <input value={form.courier_code} onChange={(e) => setForm({...form, courier_code: e.target.value})}
              className="w-full px-3 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-blue-500" /></div>
          <div className="flex gap-3 pt-2">
            <button type="button" onClick={() => setShowAdd(false)} className="px-4 py-2 border rounded-lg text-sm">取消</button>
            <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm">提交预报</button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
