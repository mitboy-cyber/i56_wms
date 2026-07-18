import { useState } from 'react';
import { ScanInput } from './ScanInput';
import { Weight } from 'lucide-react';

export default function PDAWeigh() {
  const [scan, setScan] = useState('');
  const [weight, setWeight] = useState('');
  const [msg, setMsg] = useState('');
  const submit = async () => {
    setMsg(''); if (!scan || !weight) { setMsg('请填写完整信息'); return; }
    try { const r = await fetch('/pda/api/weigh', { method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ scan, weight: +weight }) }); const d = await r.json(); setMsg(r.ok ? '称重成功' : (d?.error || '失败')); } catch { setMsg('网络错误'); }
  };
  return (
    <div>
      <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--color-ink)' }}><Weight size={20} style={{ color: 'var(--color-accent)' }} /> 称重复核</h2>
      <div className="bg-white rounded-xl border p-4 space-y-3 shadow-sm" style={{ borderColor: 'var(--border)' }}>
        <ScanInput value={scan} onChange={setScan} placeholder="扫描包裹条码" />
        <div><label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>实际重量(kg)</label>
          <input type="number" step="0.01" value={weight} onChange={e => setWeight(e.target.value)} className="w-full px-4 py-3 border rounded-lg outline-none text-center" style={{ borderColor: 'var(--border)' }} /></div>
        <button onClick={submit} className="w-full py-2.5 text-white rounded-lg font-medium" style={{ background: 'var(--color-accent)' }}>确认称重</button>
        {msg && <p className="text-sm text-center mt-2">{msg}</p>}
      </div>
    </div>
  );
}
