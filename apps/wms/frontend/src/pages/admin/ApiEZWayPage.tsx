import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiEZWayPage() {
  return (
    <GenericListPage
      title="EZWay API"
      queryKey={['admin-api-ezway']}
      queryFn={() => client.get('/admin/api/system/api-ezway')}
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
