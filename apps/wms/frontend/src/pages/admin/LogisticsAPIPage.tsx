import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function LogisticsAPIPage() {
  return (
    <GenericListPage
      title="物流API"
      queryKey={['admin-LogisticsAPIPage']}
      queryFn={() => client.get('/admin/api/system/logistics-api')}
      apiBase="/admin/api/system/logistics-api"
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
