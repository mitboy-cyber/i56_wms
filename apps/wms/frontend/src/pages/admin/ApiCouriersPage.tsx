import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiCouriersPage() {
  return (
    <GenericListPage
      title="快递API"
      queryKey={['admin-api-couriers']}
      queryFn={() => client.get('/admin/api/system/api-couriers')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'provider', label: '服务商' },
        { key: 'api_key', label: 'API Key' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
