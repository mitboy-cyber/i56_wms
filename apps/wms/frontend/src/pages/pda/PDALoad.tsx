import { useState } from 'react';
import { Truck } from 'lucide-react';

export default function PDALoad() {
  const [container, setContainer] = useState(''); const [orderNo, setOrderNo] = useState(''); const [msg, setMsg] = useState('');
  const submit = async () => {
    try { const r = await fetch('/pda/api/load', { method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ container_no: container, order_no: orderNo }) }); const d = await r.json(); setMsg(r.ok ? '装车完成' : (d?.error || '失败')); } catch { setMsg('网络错误'); }
  };
  return (<div>
    <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--color-ink)' }}><Truck size={20} style={{ color: 'var(--color-accent)' }} /> 装车</h2>
    <div className="bg-white rounded-xl border p-4 space-y-3 shadow-sm" style={{ borderColor: 'var(--border)' }}>
      <div><label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>柜号</label><input value={container} onChange={e => setContainer(e.target.value)} className="w-full px-4 py-3 border rounded-lg outline-none text-center" style={{ borderColor: 'var(--border)' }} placeholder="输入柜号" /></div>
      <div><label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>订单号</label><input value={orderNo} onChange={e => setOrderNo(e.target.value)} className="w-full px-4 py-3 border rounded-lg outline-none text-center" style={{ borderColor: 'var(--border)' }} placeholder="输入订单号" /></div>
      <button onClick={submit} className="w-full py-2.5 text-white rounded-lg font-medium" style={{ background: 'var(--color-accent)' }}>确认装车</button>
      {msg && <p className="text-sm text-center mt-2">{msg}</p>}
    </div>
  </div>);
}
