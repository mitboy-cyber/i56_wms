import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function LogisticsAPIPage() {
  return (
    <GenericListPage
      title="物流API"
      queryKey={['admin-LogisticsAPIPage']}
      queryFn={() => client.get('/admin/api/system/logistics-api')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
