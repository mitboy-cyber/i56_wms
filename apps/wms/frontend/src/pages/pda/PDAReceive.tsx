import { useState } from 'react';
import { ScanInput } from './ScanInput';
import { Package } from 'lucide-react';

export default function PDAReceive() {
  const [scan, setScan] = useState('');
  const [weight, setWeight] = useState('');
  const [length, setLength] = useState('');
  const [width, setWidth] = useState('');
  const [height, setHeight] = useState('');
  const [msg, setMsg] = useState('');

  const submit = async () => {
    setMsg('');
    try {
      const res = await fetch('/pda/api/receive', {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ scan, weight: +weight, length: +length, width: +width, height: +height }),
      });
      const d = await res.json();
      setMsg(res.ok ? (d?.message || '收货成功') : (d?.error || '收货失败'));
    } catch { setMsg('网络错误'); }
  };

  return (
    <div>
      <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--color-ink)' }}>
        <Package size={20} style={{ color: 'var(--color-accent)' }} /> 收货入库
      </h2>
      <div className="bg-white rounded-xl border p-4 space-y-3 shadow-sm" style={{ borderColor: 'var(--border)' }}>
        <ScanInput value={scan} onChange={setScan} placeholder="扫描包裹条码" />
        <div className="grid grid-cols-2 gap-2">
          <div><label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>重量(kg)</label>
            <input type="number" step="0.01" value={weight} onChange={e => setWeight(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg outline-none text-sm" style={{ borderColor: 'var(--border)' }} /></div>
          <div><label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>长(cm)</label>
            <input type="number" step="0.1" value={length} onChange={e => setLength(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg outline-none text-sm" style={{ borderColor: 'var(--border)' }} /></div>
          <div><label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>宽(cm)</label>
            <input type="number" step="0.1" value={width} onChange={e => setWidth(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg outline-none text-sm" style={{ borderColor: 'var(--border)' }} /></div>
          <div><label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>高(cm)</label>
            <input type="number" step="0.1" value={height} onChange={e => setHeight(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg outline-none text-sm" style={{ borderColor: 'var(--border)' }} /></div>
        </div>
        <button onClick={submit} className="w-full py-2.5 text-white rounded-lg font-medium"
          style={{ background: 'var(--color-accent)' }}>确认收货</button>
        {msg && <p className="text-sm text-center" style={{ color: msg.includes('成功') ? 'var(--color-accent)' : 'var(--destructive)' }}>{msg}</p>}
      </div>
    </div>
  );
}
