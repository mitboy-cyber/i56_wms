import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ServiceTypesPage() {
  return (
    <GenericListPage
      title="服务类型"
      queryKey={['admin-ServiceTypesPage']}
      queryFn={() => client.get('/admin/api/service-types')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
