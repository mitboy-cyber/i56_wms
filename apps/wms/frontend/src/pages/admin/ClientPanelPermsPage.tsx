import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ClientPanelPermsPage() {
  return (
    <GenericListPage title="客户端权限" queryKey={['admin-ClientPanelPerms']}
      queryFn={() => client.get('/admin/api/client-panel-perms').then(r => r.data)}
      columns={[
        { key: 'id', label: '编号' }, { key: 'client_id', label: '客户' },
        { key: 'menu_name', label: '菜单名称' },
        { key: 'can_view', label: '查看', render: (v: unknown) => <>{v ? '✅' : '❌'}</> },
        { key: 'can_operate', label: '操作', render: (v: unknown) => <>{v ? '✅' : '❌'}</> },
      ]} getRowId={(_, i) => String(i)} />
  );
}
