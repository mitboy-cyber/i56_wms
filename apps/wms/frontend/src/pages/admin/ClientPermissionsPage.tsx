import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ClientPermissionsPage() {
  return (
    <GenericListPage
      title="客户权限"
      queryKey={['admin-ClientPermissionsPage']}
      queryFn={() => client.get('/admin/api/client-permissions')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
