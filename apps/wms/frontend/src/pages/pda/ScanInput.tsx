import { useState, useRef, useEffect } from 'react';
interface Props { placeholder: string; onSubmit: (scan: string) => void; loading?: boolean; }
export default function ScanInput({ placeholder, onSubmit, loading }: Props) {
  const [scan, setScan] = useState('');
  const ref = useRef<HTMLInputElement>(null);
  useEffect(() => { ref.current?.focus(); }, []);
  return (
    <div className="flex gap-2">
      <input ref={ref} value={scan} onChange={e => setScan(e.target.value)}
        placeholder={placeholder} className="flex-1 px-4 py-3 border rounded-xl text-lg outline-none focus:ring-2 focus:ring-blue-500"
        onKeyDown={e => { if (e.key==='Enter' && scan) { onSubmit(scan); setScan(''); }}} autoFocus />
      <button onClick={() => { if (scan) { onSubmit(scan); setScan(''); }}}
        disabled={loading} className="px-6 py-3 bg-blue-600 text-white rounded-xl font-medium disabled:opacity-50">
        {loading ? '处理中' : '确认'}
      </button>
    </div>
  );
}
