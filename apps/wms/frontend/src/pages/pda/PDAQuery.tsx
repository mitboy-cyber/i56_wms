import { useState } from 'react';
import { ScanInput } from './ScanInput';
import { Search } from 'lucide-react';

export default function PDAQuery() {
  const [scan, setScan] = useState(''); const [result, setResult] = useState('');
  const submit = async () => {
    try { const r = await fetch('/pda/api/query', { method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ scan }) }); const d = await r.json(); setResult(r.ok ? JSON.stringify(d, null, 2) : (d?.error || '查询失败')); } catch { setResult('网络错误'); }
  };
  return (<div>
    <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--color-ink)' }}><Search size={20} style={{ color: 'var(--color-accent)' }} /> 包裹查询</h2>
    <div className="bg-white rounded-xl border p-4 space-y-3 shadow-sm" style={{ borderColor: 'var(--border)' }}>
      <ScanInput value={scan} onChange={setScan} placeholder="扫描或输入包裹条码" />
      <button onClick={submit} className="w-full py-2.5 text-white rounded-lg font-medium" style={{ background: 'var(--color-accent)' }}>查询</button>
      {result && <pre className="text-xs mt-3 p-3 bg-gray-50 rounded-lg overflow-auto max-h-64 whitespace-pre-wrap">{result}</pre>}
    </div>
  </div>);
}
