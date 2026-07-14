import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation, useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import clientApi from '@/api/clientApi';

export default function ClientOrderNew() {
  const navigate = useNavigate();
  const [form, setForm] = useState({ order_no: '', recipient_name: '', parcel_count: 1, total_price: 0, route_id: 1 });
  const { data: routes } = useQuery({ queryKey: ['client-routes'], queryFn: () => clientApi.routePrices() } as any);
  const routeList: any[] = (routes as any)?.data ?? [];

  const create = useMutation({
    mutationFn: (d: any) => client.post('/admin/api/orders', d),
    onSuccess: () => navigate('/client/orders'),
  } as any);

  return (
    <div>
      <h2 className="text-xl font-bold text-gray-800 mb-6">新建集运订单</h2>
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 max-w-lg">
        <form onSubmit={(e) => { e.preventDefault(); (create as any).mutate(form); }} className="space-y-4">
          <div><label className="block text-sm font-medium mb-1">订单号</label>
            <input required value={form.order_no} onChange={e => setForm({...form, order_no: e.target.value})}
              className="w-full px-3 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-blue-500" /></div>
          <div><label className="block text-sm font-medium mb-1">收件人</label>
            <input required value={form.recipient_name} onChange={e => setForm({...form, recipient_name: e.target.value})}
              className="w-full px-3 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-blue-500" /></div>
          <div><label className="block text-sm font-medium mb-1">包裹数</label>
            <input type="number" value={form.parcel_count} onChange={e => setForm({...form, parcel_count: +e.target.value})}
              className="w-full px-3 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-blue-500" /></div>
          <div><label className="block text-sm font-medium mb-1">预估金额</label>
            <input type="number" step="0.01" value={form.total_price} onChange={e => setForm({...form, total_price: +e.target.value})}
              className="w-full px-3 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-blue-500" /></div>
          <div><label className="block text-sm font-medium mb-1">线路</label>
            <select value={form.route_id} onChange={e => setForm({...form, route_id: +e.target.value})}
              className="w-full px-3 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-blue-500">
              {routeList.map((r: any, i: number) => <option key={i} value={i+1}>{r.route_name || `线路#${i+1}`}</option>)}
            </select></div>
          <div className="flex gap-3 pt-2">
            <button type="button" onClick={() => navigate('/client/orders')} className="px-4 py-2 border rounded-lg text-sm">取消</button>
            <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm">创建订单</button>
          </div>
        </form>
      </div>
    </div>
  );
}
