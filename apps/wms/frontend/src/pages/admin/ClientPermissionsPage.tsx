import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function ClientPermissionsPage() {
  return (
    <MinimalListPage title="客户权限" queryKey={['admin-ClientPermissionsPage']}
      queryFn={() => client.get('/admin/api/client-permissions')}
      apiBase="/admin/api/client-permissions"
      columns={[
        { key: 'id', label: '编号' }, { key: 'client_id', label: '客户编号' },
        { key: 'module', label: '模块' }, { key: 'can_read', label: '查看' },
        { key: 'can_write', label: '编辑' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
