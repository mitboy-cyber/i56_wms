import { useState, useEffect } from 'react';
import { Save, Key, Cpu, Shield, AlertCircle } from 'lucide-react';
import client from '@/api/client';

interface AIModel {
  id: string; name: string; provider: string; enabled: boolean;
}
interface APIPermission {
  id: string; name: string; scope: string; enabled: boolean;
}

export default function AISettingsPage() {
  const [models, setModels] = useState<AIModel[]>([
    { id: 'gpt-4o', name: 'GPT-4o', provider: 'OpenAI', enabled: true },
    { id: 'gpt-4o-mini', name: 'GPT-4o Mini', provider: 'OpenAI', enabled: true },
    { id: 'claude-sonnet-4', name: 'Claude Sonnet 4', provider: 'Anthropic', enabled: false },
    { id: 'deepseek-v4', name: 'DeepSeek V4', provider: 'DeepSeek', enabled: true },
    { id: 'qwen-max', name: 'Qwen Max', provider: 'Alibaba', enabled: false },
  ]);
  const [permissions, setPermissions] = useState<APIPermission[]>([
    { id: 'orders:read', name: '订单查询', scope: 'orders:read', enabled: true },
    { id: 'orders:write', name: '订单修改', scope: 'orders:write', enabled: false },
    { id: 'parcels:read', name: '包裹查询', scope: 'parcels:read', enabled: true },
    { id: 'warehouse:read', name: '仓库数据', scope: 'warehouse:read', enabled: true },
    { id: 'finance:read', name: '财务数据', scope: 'finance:read', enabled: false },
    { id: 'customers:read', name: '客户数据', scope: 'customers:read', enabled: false },
    { id: 'reports:generate', name: '报表生成', scope: 'reports:generate', enabled: true },
    { id: 'system:config', name: '系统配置', scope: 'system:config', enabled: false },
  ]);
  const [apiKey, setApiKey] = useState('sk-••••••••••••••••••••••••');
  const [saved, setSaved] = useState(false);

  const toggleModel = (id: string) => {
    setModels(prev => prev.map(m => m.id === id ? {...m, enabled: !m.enabled} : m));
  };
  const togglePerm = (id: string) => {
    setPermissions(prev => prev.map(p => p.id === id ? {...p, enabled: !p.enabled} : p));
  };

  const handleSave = async () => {
    try {
      await client.post('/admin/api/system/ai-settings', {
        models: models.filter(m => m.enabled).map(m => m.id),
        permissions: permissions.filter(p => p.enabled).map(p => p.scope),
        api_key: apiKey.replace(/•/g, ''),
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch { /* ignore */ }
  };

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-bold text-gray-800">AI 配置中心</h2>
          <p className="text-sm text-gray-500 mt-0.5">模型选择、API 权限与密钥管理</p>
        </div>
        <button onClick={handleSave}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm">
          <Save size={15} /> {saved ? '已保存 ✓' : '保存配置'}
        </button>
      </div>

      {/* API Key */}
      <section className="bg-white rounded-xl border border-gray-200 p-5">
        <div className="flex items-center gap-2 mb-3">
          <Key size={16} className="text-amber-500" />
          <h3 className="font-semibold text-gray-800">API 密钥</h3>
        </div>
        <input type="password" value={apiKey} onChange={e => setApiKey(e.target.value)}
          className="w-full px-4 py-2.5 rounded-lg border border-gray-300 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500" />
        <p className="text-xs text-gray-400 mt-1.5">用于调用 AI 模型的 API 密钥，加密存储</p>
      </section>

      {/* Models */}
      <section className="bg-white rounded-xl border border-gray-200 p-5">
        <div className="flex items-center gap-2 mb-4">
          <Cpu size={16} className="text-blue-500" />
          <h3 className="font-semibold text-gray-800">可用模型</h3>
        </div>
        <div className="space-y-2">
          {models.map(m => (
            <div key={m.id} className="flex items-center justify-between py-2 px-3 rounded-lg hover:bg-gray-50">
              <div className="flex items-center gap-3">
                <button onClick={() => toggleModel(m.id)}
                  className={`w-10 h-5 rounded-full transition-colors relative ${m.enabled ? 'bg-blue-600' : 'bg-gray-300'}`}>
                  <span className={`absolute top-0.5 w-4 h-4 rounded-full bg-white transition-transform ${m.enabled ? 'left-5' : 'left-0.5'}`} />
                </button>
                <div>
                  <p className="text-sm font-medium text-gray-800">{m.name}</p>
                  <p className="text-xs text-gray-400">{m.provider}</p>
                </div>
              </div>
              <span className={`text-xs px-2 py-0.5 rounded ${m.enabled ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-400'}`}>
                {m.enabled ? '已启用' : '已禁用'}
              </span>
            </div>
          ))}
        </div>
      </section>

      {/* API Permissions */}
      <section className="bg-white rounded-xl border border-gray-200 p-5">
        <div className="flex items-center gap-2 mb-4">
          <Shield size={16} className="text-purple-500" />
          <h3 className="font-semibold text-gray-800">API 权限范围</h3>
        </div>
        <div className="space-y-2">
          {permissions.map(p => (
            <div key={p.id} className="flex items-center justify-between py-2 px-3 rounded-lg hover:bg-gray-50">
              <div className="flex items-center gap-3">
                <button onClick={() => togglePerm(p.id)}
                  className={`w-10 h-5 rounded-full transition-colors relative ${p.enabled ? 'bg-purple-600' : 'bg-gray-300'}`}>
                  <span className={`absolute top-0.5 w-4 h-4 rounded-full bg-white transition-transform ${p.enabled ? 'left-5' : 'left-0.5'}`} />
                </button>
                <div>
                  <p className="text-sm font-medium text-gray-800">{p.name}</p>
                  <code className="text-xs text-gray-400 bg-gray-100 px-1 rounded">{p.scope}</code>
                </div>
              </div>
            </div>
          ))}
        </div>
        <div className="mt-4 flex items-start gap-2 bg-amber-50 rounded-lg p-3 border border-amber-200">
          <AlertCircle size={14} className="text-amber-600 mt-0.5 shrink-0" />
          <p className="text-xs text-amber-700">
            权限范围控制 AI 助手可以访问的数据。建议仅开启必要的权限，保障数据安全。
          </p>
        </div>
      </section>
    </div>
  );
}
