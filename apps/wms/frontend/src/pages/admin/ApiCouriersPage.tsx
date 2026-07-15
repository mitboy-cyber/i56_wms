import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiCouriersPage() {
  return (
    <GenericListPage
      title="快递API"
      queryKey={['admin-api-couriers']}
      queryFn={() => client.get('/admin/api/system/api-couriers')}
      apiBase="/admin/api/system/api-couriers"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'name', label: '名称' },
        { key: 'provider', label: '提供商' },
        { key: 'endpoint', label: '接口地址' },
        { key: 'api_key', label: 'API密钥' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
