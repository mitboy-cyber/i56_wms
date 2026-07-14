import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function CustomsBrokerAPIPage() {
  return (
    <GenericListPage
      title="报关API"
      queryKey={['admin-CustomsBrokerAPIPage']}
      queryFn={() => client.get('/admin/api/system/customs-broker-api')}
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
