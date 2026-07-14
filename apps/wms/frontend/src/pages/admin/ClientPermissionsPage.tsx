import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ClientPermissionsPage() {
  return (
    <GenericListPage
      title="客户权限"
      queryKey={['admin-ClientPermissionsPage']}
      queryFn={() => client.get('/admin/api/client-permissions')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'client_id', label: 'Client Id' },
        { key: 'module', label: 'Module' },
        { key: 'can_read', label: 'Can Read' },
        { key: 'can_write', label: 'Can Write' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
