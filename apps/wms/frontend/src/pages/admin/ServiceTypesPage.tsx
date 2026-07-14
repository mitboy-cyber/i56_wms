import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ServiceTypesPage() {
  return (
    <GenericListPage
      title="服务类型"
      queryKey={['admin-ServiceTypesPage']}
      queryFn={() => client.get('/admin/api/service-types')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'code', label: 'Code' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
