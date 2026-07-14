import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import pdaApi from '@/api/pdaApi';
import ScanInput from './ScanInput';
export default function PDAWeigh() {
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState('');
  const qc = useQueryClient();
  const mut = useMutation({
    mutationFn: (scan: string) => (pdaApi as any).weigh({ scan } as any),
    onSuccess: (res: any) => { setResult(res.data); setError(''); qc.invalidateQueries({ queryKey: ['pda-dashboard'] } as any); },
    onError: (err: any) => { setError(err.response?.data?.error || err.message); setResult(null); },
  } as any);
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">称重</h2>
      <ScanInput placeholder="称重扫码..." onSubmit={(s) => (mut as any).mutate(s)} loading={mut.isPending as boolean} />
      {error && <div className="mt-3 p-3 bg-red-50 border border-red-200 rounded-xl text-red-600 text-sm">{error}</div>}
      {result && <div className="mt-3 p-3 bg-green-50 border border-green-200 rounded-xl text-sm"><pre className="text-gray-700 overflow-auto max-h-60">{JSON.stringify(result, null, 2)}</pre></div>}
    </div>
  );
}
