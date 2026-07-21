import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function ApiEZWayPage() {
  return (
    <MinimalListPage
      title="EZ Way API"
      queryKey={['admin-api-ezway']}
      queryFn={() => client.get('/admin/api/system/api-ezway')}
      apiBase="/admin/api/system/api-ezway"
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
