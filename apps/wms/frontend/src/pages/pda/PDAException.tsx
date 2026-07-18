import { useState } from 'react';
import { ScanInput } from './ScanInput';
import { AlertTriangle } from 'lucide-react';

export default function PDAException() {
  const [scan, setScan] = useState(''); const [reason, setReason] = useState(''); const [msg, setMsg] = useState('');
  const submit = async () => {
    try { const r = await fetch('/pda/api/exception', { method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ scan, reason }) }); const d = await r.json(); setMsg(r.ok ? '异常已上报' : (d?.error || '失败')); } catch { setMsg('网络错误'); }
  };
  return (<div>
    <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--color-ink)' }}><AlertTriangle size={20} style={{ color: 'var(--destructive)' }} /> 异常上报</h2>
    <div className="bg-white rounded-xl border p-4 space-y-3 shadow-sm" style={{ borderColor: 'var(--border)' }}>
      <ScanInput value={scan} onChange={setScan} placeholder="扫描异常包裹条码" />
      <div><label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>异常原因</label><input value={reason} onChange={e => setReason(e.target.value)} className="w-full px-4 py-3 border rounded-lg outline-none" style={{ borderColor: 'var(--border)' }} placeholder="请描述异常原因" /></div>
      <button onClick={submit} className="w-full py-2.5 text-white rounded-lg font-medium" style={{ background: 'var(--destructive)' }}>上报异常</button>
      {msg && <p className="text-sm text-center mt-2">{msg}</p>}
    </div>
  </div>);
}
