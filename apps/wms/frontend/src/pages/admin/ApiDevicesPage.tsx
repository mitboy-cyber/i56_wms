import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiDevicesPage() {
  return (
    <GenericListPage
      title="设备API"
      queryKey={['admin-api-devices']}
      queryFn={() => client.get('/admin/api/system/api-devices')}
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
