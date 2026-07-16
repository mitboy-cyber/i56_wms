import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import client from '@/api/client';
import { ChevronDown, ChevronRight, Check, X, Save } from 'lucide-react';

// ── Client types ──
const CLIENT_TYPES = ['platform', 'shopee', 'taobao', 'pdd', 'jd', 'douyin'];
const CLIENT_LABELS: Record<string, string> = {
  platform: '平台客户', shopee: '蝦皮商家', taobao: '淘宝商家',
  pdd: '拼多多商家', jd: '京东商家', douyin: '抖音商家',
};

// ── Feature groups ──
const FEATURE_GROUPS = [
  {
    label: '包裹管理',
    features: ['我的包裹', '预报包裹', '认领包裹', 'Excel批量导入预报'],
  },
  {
    label: '订单管理',
    features: ['我的订单', '集运下单', '取消订单', '打包下单', '附加服务订单'],
  },
  {
    label: '客户资料',
    features: ['客户会员', '收件地址', '申报人'],
  },
  {
    label: '财务管理',
    features: ['余额明细', '月结对账单', '路线价格', '承运商派送价', '承运商加收价'],
  },
  {
    label: '系统设置',
    features: ['仓库信息', 'Webhook投递', 'API凭证', '客服工单'],
  },
];

// ── Helpers ──
function getDefaultPermKey(type: string, feature: string) {
  return `perm-${type}-${feature}`;
}

// ── Page component ──
export default function ClientPanelPermsPage() {
  const qc = useQueryClient();
  const [expanded, setExpanded] = useState<Record<string, boolean>>({ platform: true, shopee: true });
  // Checkbox state: { "perm-platform-我的包裹": true }
  const [checks, setChecks] = useState<Record<string, boolean>>({});
  const [dirty, setDirty] = useState(false);

  // Load permissions
  const { data: perms = [], isLoading } = useQuery<any[]>({
    queryKey: ['admin-client-panel-perms'],
    queryFn: () => client.get('/admin/api/client-panel-perms'),
  });

  // Import on load
  useEffect(() => {
    if (perms.length > 0 && !dirty) {
      const init: Record<string, boolean> = {};
      perms.forEach((p: any) => {
        init[getDefaultPermKey(p.client_type || p.level, p.menu_name || p.feature)] = p.enabled ?? p.can_view ?? true;
      });
      // Default: all enabled for empty state
      CLIENT_TYPES.forEach(t => {
        FEATURE_GROUPS.forEach(g => {
          g.features.forEach(f => {
            const k = getDefaultPermKey(t, f);
            if (!(k in init)) init[k] = true;
          });
        });
      });
      setChecks(init);
    }
  });

  // Toggle type
  const toggleType = (t: string) => setExpanded(p => ({ ...p, [t]: !p[t] }));

  // Toggle feature
  const toggleFeature = (type: string, feature: string) => {
    const k = getDefaultPermKey(type, feature);
    setChecks(prev => ({ ...prev, [k]: !prev[k] }));
    setDirty(true);
  };

  // Select/Deselect all in a group
  const selectAll = (type: string, groupFeatures: string[], val: boolean) => {
    setChecks(prev => {
      const next = { ...prev };
      groupFeatures.forEach(f => { next[getDefaultPermKey(type, f)] = val; });
      return next;
    });
    setDirty(true);
  };

  // Save
  const saveMut = useMutation({
    mutationFn: (data: any) => client.post('/admin/api/client-panel-perms/batch', data),
    onSuccess: () => { alert('保存成功'); setDirty(false); qc.invalidateQueries({ queryKey: ['admin-client-panel-perms'] }); },
  });

  const handleSave = () => {
    const payload = CLIENT_TYPES.flatMap(type =>
      FEATURE_GROUPS.flatMap(g =>
        g.features.map(feature => ({
          client_type: type,
          feature_group: g.label,
          feature,
          enabled: checks[getDefaultPermKey(type, feature)] ?? true,
        }))
      )
    );
    saveMut.mutate(payload);
  };

  if (isLoading) return <div className="p-8 text-gray-400">加载中...</div>;

  return (
    <div className="max-w-5xl">
      {/* Header */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
        <p className="text-sm text-blue-800">
          正在配置公司：<strong>嗨购邦集团公司</strong>
        </p>
        <p className="text-xs text-blue-600 mt-1">勾选 = 该客户类型默认可用；全勾 = 全开；取消 = 关闭。</p>
      </div>

      {/* Matrix */}
      {CLIENT_TYPES.map(type => (
        <div key={type} className="mb-6 bg-white border rounded-lg shadow-sm">
          {/* Type header */}
          <button
            onClick={() => toggleType(type)}
            className="w-full flex items-center justify-between px-5 py-3 bg-gray-50 hover:bg-gray-100 rounded-t-lg border-b"
          >
            <h2 className="text-lg font-semibold text-gray-800">{CLIENT_LABELS[type]}</h2>
            {expanded[type] ? <ChevronDown size={20} /> : <ChevronRight size={20} />}
          </button>

          {/* Feature groups */}
          {expanded[type] && (
            <div className="px-5 py-4 space-y-4">
              {FEATURE_GROUPS.map(group => {
                const allChecked = group.features.every(f => checks[getDefaultPermKey(type, f)]);
                const noneChecked = group.features.every(f => !checks[getDefaultPermKey(type, f)]);
                return (
                  <div key={group.label}>
                    <div className="flex items-center gap-3 mb-2">
                      <span className="text-sm font-medium text-gray-600 w-20">{group.label}</span>
                      <button
                        onClick={() => selectAll(type, group.features, !allChecked)}
                        className="text-xs px-2 py-0.5 rounded border text-gray-500 hover:bg-gray-100"
                      >
                        {allChecked ? '取消全选' : '全选'}
                      </button>
                    </div>
                    <div className="flex flex-wrap gap-2 ml-20">
                      {group.features.map(feature => {
                        const checked = checks[getDefaultPermKey(type, feature)] ?? true;
                        return (
                          <label
                            key={feature}
                            className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full border cursor-pointer text-sm transition-colors ${
                              checked
                                ? 'bg-green-50 border-green-300 text-green-800'
                                : 'bg-red-50 border-red-200 text-red-500'
                            }`}
                          >
                            <input
                              type="checkbox"
                              className="sr-only"
                              checked={checked}
                              onChange={() => toggleFeature(type, feature)}
                            />
                            {checked ? <Check size={12} /> : <X size={12} />}
                            {feature}
                          </label>
                        );
                      })}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      ))}

      {/* Save bar */}
      {dirty && (
        <div className="fixed bottom-0 left-64 right-0 bg-white border-t shadow-lg px-6 py-3 flex items-center justify-between z-50">
          <span className="text-sm text-orange-600">有未保存的修改</span>
          <button
            onClick={handleSave}
            disabled={saveMut.isPending}
            className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            <Save size={16} />
            {saveMut.isPending ? '保存中...' : '保存设置'}
          </button>
        </div>
      )}
    </div>
  );
}
