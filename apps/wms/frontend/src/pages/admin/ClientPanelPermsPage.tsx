import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';
import React from 'react';

const LEVEL_MAP: Record<string, string> = { enterprise: '🏢 企业版', pro: '💼 专业版', basic: '📦 基础版' };
const STATUS_MAP: Record<string, string> = { active: '✅ 生效中', expired: '❌ 已过期', suspended: '⏸️ 已暂停', pending: '🕐 待审批' };

export default function ClientPanelPermsPage() {
  return (
    <GenericListPage title="客户端权限" queryKey={['admin-ClientPanelPerms']}
      queryFn={() => client.get('/admin/api/client-panel-perms').then(r => r.data)}
      apiBase="/admin/api/client-panel-perms"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'client_name', label: '客户' },
        { key: 'module', label: '模块', render: (v: unknown) => <span className="text-blue-600 font-medium">{String(v)}</span> },
        { key: 'menu_name', label: '菜单' },
        { key: 'level', label: '等级', render: (v: unknown) => <>{LEVEL_MAP[String(v)] || String(v)}</> },
        { key: 'status', label: '状态', render: (v: unknown) => <span className={String(v)==='active'?'text-green-600 font-semibold':'text-red-500'}>{STATUS_MAP[String(v)]||String(v)}</span> },
        { key: 'can_view', label: '查看', render: (v: unknown) => <>{v ? '✅' : '❌'}</> },
        { key: 'can_create', label: '新建', render: (v: unknown) => <>{v ? '✅' : '❌'}</> },
        { key: 'can_edit', label: '编辑', render: (v: unknown) => <>{v ? '✅' : '❌'}</> },
        { key: 'can_delete', label: '删除', render: (v: unknown) => <>{v ? '✅' : '❌'}</> },
        { key: 'can_export', label: '导出', render: (v: unknown) => <>{v ? '✅' : '❌'}</> },
        { key: 'expires_at', label: '过期时间', render: (v: unknown) => <span className="text-xs text-gray-500">{v ? new Date(String(v)).toLocaleDateString('zh-CN') : '—'}</span> },
        { key: 'remarks', label: '备注' },
      ]} getRowId={(r: any) => String(r.id)} />
  );
}
