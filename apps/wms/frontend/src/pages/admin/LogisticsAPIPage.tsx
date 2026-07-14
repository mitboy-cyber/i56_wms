import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function LogisticsAPIPage() {
  return (
    <GenericListPage
      title="物流API"
      queryKey={['admin-LogisticsAPIPage']}
      queryFn={() => client.get('/admin/api/system/logistics-api')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'provider', label: 'Provider' },
        { key: 'endpoint', label: 'Endpoint' },
        { key: 'api_key', label: 'Api Key' },
        { key: 'status', label: 'Status' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
