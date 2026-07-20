import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useClientAuth } from '@/stores/clientAuth';

export default function ClientLogin() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login } = useClientAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    const ok = await login(username, password);
    if (ok) navigate('/client/dashboard');
    else setError('账号或密码错误');
  };

  return (
    <div className="min-h-screen flex items-center justify-center" style={{ background: 'var(--background)' }}>
      <div className="bg-white rounded-2xl shadow-xl p-8 w-full max-w-md border" style={{ borderColor: 'var(--border)' }}>
        <div className="text-center mb-6">
          <h1 className="text-2xl font-bold" style={{ color: 'var(--color-accent)' }}>I56 客户中心</h1>
          <p className="text-sm mt-1" style={{ color: 'var(--color-neutral)' }}>客户端管理平台</p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1" style={{ color: 'var(--color-muted)' }}>账号</label>
            <input type="text" value={username} onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-2.5 border rounded-lg outline-none focus:ring-2 transition-shadow"
              style={{ borderColor: 'var(--border)', '--tw-ring-color': 'var(--ring)' } as React.CSSProperties} />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1" style={{ color: 'var(--color-muted)' }}>密码</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2.5 border rounded-lg outline-none focus:ring-2 transition-shadow"
              style={{ borderColor: 'var(--border)', '--tw-ring-color': 'var(--ring)' } as React.CSSProperties} />
          </div>
          {error && <p className="text-red-500 text-sm">{error}</p>}
          <button type="submit" className="w-full py-2.5 text-white rounded-lg font-medium transition-colors"
            style={{ background: 'var(--color-accent)' }}>
            登录
          </button>
        </form>
      </div>
    </div>
  );
}
