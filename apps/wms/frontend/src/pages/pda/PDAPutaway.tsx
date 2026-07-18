import { useState } from 'react';
import { ScanInput } from './ScanInput';
import { ArrowDownToLine } from 'lucide-react';

export default function PDAPutaway() {
  const [scan, setScan] = useState(''); const [loc, setLoc] = useState(''); const [msg, setMsg] = useState('');
  const submit = async () => {
    try { const r = await fetch('/pda/api/putaway', { method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ scan, location_barcode: loc }) }); const d = await r.json(); setMsg(r.ok ? '上架成功' : (d?.error || '失败')); } catch { setMsg('网络错误'); }
  };
  return (
    <div>
      <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--color-ink)' }}><ArrowDownToLine size={20} style={{ color: 'var(--color-accent)' }} /> 上架</h2>
      <div className="bg-white rounded-xl border p-4 space-y-3 shadow-sm" style={{ borderColor: 'var(--border)' }}>
        <ScanInput value={scan} onChange={setScan} placeholder="扫描包裹条码" />
        <div><label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>货位条码</label>
          <input value={loc} onChange={e => setLoc(e.target.value)} className="w-full px-4 py-3 border rounded-lg outline-none text-center" style={{ borderColor: 'var(--border)' }} placeholder="扫描货位条码" /></div>
        <button onClick={submit} className="w-full py-2.5 text-white rounded-lg font-medium" style={{ background: 'var(--color-accent)' }}>确认上架</button>
        {msg && <p className="text-sm text-center mt-2">{msg}</p>}
      </div>
    </div>
  );
}
