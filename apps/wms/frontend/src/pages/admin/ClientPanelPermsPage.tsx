import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';
import React from 'react';

const LEVEL_CN: Record<string, string> = { enterprise: '🏢 企业版', pro: '💼 专业版', basic: '📦 基础版' };
const STATUS_CN: Record<string, string> = { active: '✅ 生效中', expired: '❌ 已过期', suspended: '⏸️ 已暂停' };

export default function ClientPanelPermsPage() {
  return (
    <GenericListPage title="客户端权限" queryKey={['admin-ClientPanelPerms']}
      queryFn={() => client.get('/admin/api/client-panel-perms')}
      apiBase="/admin/api/client-panel-perms"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'client_name', label: '客户' },
        { key: 'module', label: '模块' },
        { key: 'menu_name', label: '菜单' },
        { key: 'level', label: '等级', render: (v: any) => <>{(LEVEL_CN as any)[v] || v}</> },
        { key: 'status', label: '状态', render: (v: any) => <>{(STATUS_CN as any)[v] || v}</> },
        { key: 'can_view', label: '查看', render: (v: any) => <>{v ? '✅' : '❌'}</> },
        { key: 'can_create', label: '新建', render: (v: any) => <>{v ? '✅' : '❌'}</> },
        { key: 'can_edit', label: '编辑', render: (v: any) => <>{v ? '✅' : '❌'}</> },
        { key: 'can_delete', label: '删除', render: (v: any) => <>{v ? '✅' : '❌'}</> },
        { key: 'can_export', label: '导出', render: (v: any) => <>{v ? '✅' : '❌'}</> },
        { key: 'expires_at', label: '过期时间' },
        { key: 'remarks', label: '备注' },
      ]} getRowId={(r: any) => String(r.id)} />
  );
}
