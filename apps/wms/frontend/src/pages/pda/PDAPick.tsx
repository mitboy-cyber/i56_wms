import { useState } from 'react';

export default function PDAPick() {
  const [orderNo, setOrderNo] = useState('');
  const [msg, setMsg] = useState('');

  const submit = async () => {
    setMsg('');
    try {
      const res = await fetch('/pda/api/pick', {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ order_no: orderNo }),
      });
      const d = await res.json();
      setMsg(res.ok ? '拣货完成' : (d?.error || '拣货失败'));
    } catch { setMsg('网络错误'); }
  };

  return (
    <div>
      <h2 className="text-lg font-bold mb-4" style={{ color: 'var(--color-ink)' }}>🛒 拣货</h2>
      <div className="bg-white rounded-xl border p-4 space-y-3 shadow-sm" style={{ borderColor: 'var(--border)' }}>
        <div>
          <label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>订单号</label>
          <input value={orderNo} onChange={e => setOrderNo(e.target.value)}
            className="w-full px-4 py-3 border rounded-lg outline-none text-center" style={{ borderColor: 'var(--border)' }}
            placeholder="输入或扫描订单号" />
        </div>
        <button onClick={submit} className="w-full py-2.5 text-white rounded-lg font-medium"
          style={{ background: 'var(--color-accent)' }}>确认拣货</button>
        {msg && <p className="text-sm text-center mt-2">{msg}</p>}
      </div>
    </div>
  );
}
