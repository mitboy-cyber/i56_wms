import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function CustomsBrokerAPIPage() {
  return (
    <GenericListPage
      title="报关API"
      queryKey={['admin-CustomsBrokerAPIPage']}
      queryFn={() => client.get('/admin/api/system/customs-broker-api')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
