import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function RolesPage() {
  return (
    <MinimalListPage title="角色管理" queryKey={['admin-roles']}
      queryFn={() => client.get('/admin/api/roles')}
      apiBase="/admin/api/roles"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '角色名称' },
        { key: 'description', label: '描述' },
        { key: 'is_active', label: '启用', render: (v: unknown) => <>{v ? '✅' : '❌'}</> },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
