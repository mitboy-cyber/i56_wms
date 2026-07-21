import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function CustomsBrokerAPIPage() {
  return (
    <MinimalListPage
      title="报关API"
      queryKey={['admin-CustomsBrokerAPIPage']}
      queryFn={() => client.get('/admin/api/system/customs-broker-api')}
      apiBase="/admin/api/system/customs-broker-api"
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
