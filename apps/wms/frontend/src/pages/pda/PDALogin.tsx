import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { usePDAAuth } from '@/stores/pdaAuth';
export default function PDALogin() {
  const [code, setCode] = useState('OP001');
  const [pin, setPin] = useState('1234');
  const [error, setError] = useState('');
  const { login } = usePDAAuth();
  const nav = useNavigate();
  const submit = async (e: React.FormEvent) => { e.preventDefault(); setError('');
    const ok = await login(code, pin);
    if (ok) nav('/pda/dashboard'); else setError('工号或PIN错误');
  };
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-cyan-100 p-4">
      <div className="bg-white rounded-2xl shadow-xl p-8 w-full max-w-sm">
        <h1 className="text-2xl font-bold text-center text-blue-600 mb-6">I56 PDA</h1>
        <form onSubmit={submit} className="space-y-4">
          <input value={code} onChange={e=>setCode(e.target.value)} placeholder="工号" className="w-full px-4 py-3 text-lg border rounded-xl outline-none focus:ring-2 focus:ring-blue-500 text-center" />
          <input type="password" value={pin} onChange={e=>setPin(e.target.value)} placeholder="PIN码" className="w-full px-4 py-3 text-lg border rounded-xl outline-none focus:ring-2 focus:ring-blue-500 text-center" />
          {error && <p className="text-red-500 text-sm text-center">{error}</p>}
          <button type="submit" className="w-full py-3 bg-blue-600 text-white rounded-xl text-lg font-medium">登录</button>
        </form>
      </div>
    </div>
  );
}
