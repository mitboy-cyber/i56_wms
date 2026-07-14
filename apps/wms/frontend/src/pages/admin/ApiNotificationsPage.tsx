import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiNotificationsPage() {
  return (
    <GenericListPage
      title="通知API"
      queryKey={['admin-api-notifications']}
      queryFn={() => client.get('/admin/api/system/api-notifications')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'provider', label: 'Provider' },
        { key: 'endpoint', label: 'Endpoint' },
        { key: 'api_key', label: 'Api Key' },
        { key: 'status', label: 'Status' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
