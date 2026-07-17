import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { usePDAAuth } from '@/stores/pdaAuth';

export default function PDALogin() {
  const [code, setCode] = useState('OP001');
  const [pin, setPin] = useState('1234');
  const [error, setError] = useState('');
  const { login } = usePDAAuth();
  const nav = useNavigate();

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    const ok = await login(code, pin);
    if (ok) nav('/pda/dashboard');
    else setError('工号或PIN码错误');
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4" style={{ background: 'var(--background)' }}>
      <div className="w-full max-w-sm">
        <div className="text-center mb-6">
          <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl mb-3" style={{ background: 'var(--color-accent-light)' }}>
            <span className="text-2xl font-bold" style={{ color: 'var(--color-accent)' }}>P</span>
          </div>
          <h1 className="text-xl font-bold" style={{ color: 'var(--color-ink)' }}>I56 PDA 手持终端</h1>
          <p className="text-sm mt-1" style={{ color: 'var(--color-neutral)' }}>仓库作业系统</p>
        </div>

        <div className="bg-white rounded-2xl shadow-xl border p-6" style={{ borderColor: 'var(--border)' }}>
          <form onSubmit={submit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1" style={{ color: 'var(--color-muted)' }}>工号</label>
              <input value={code} onChange={e => setCode(e.target.value)}
                className="w-full px-4 py-3 border rounded-xl outline-none text-center focus:ring-2 transition-shadow"
                style={{ borderColor: 'var(--border)', '--tw-ring-color': 'var(--ring)' } as React.CSSProperties}
                placeholder="请输入工号" />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1" style={{ color: 'var(--color-muted)' }}>PIN码</label>
              <input type="password" value={pin} onChange={e => setPin(e.target.value)}
                className="w-full px-4 py-3 border rounded-xl outline-none text-center focus:ring-2 transition-shadow"
                style={{ borderColor: 'var(--border)', '--tw-ring-color': 'var(--ring)' } as React.CSSProperties}
                placeholder="请输入PIN码" />
            </div>
            {error && <p className="text-sm text-center" style={{ color: 'var(--destructive)' }}>{error}</p>}
            <button type="submit" className="w-full py-3 text-white rounded-xl font-medium transition-colors"
              style={{ background: 'var(--color-accent)' }}>
              登录系统
            </button>
          </form>
        </div>

        <p className="text-xs text-center mt-4" style={{ color: 'var(--color-neutral)' }}>
          请使用仓库分配的工号和PIN码登录
        </p>
      </div>
    </div>
  );
}
