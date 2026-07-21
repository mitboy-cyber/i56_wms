import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function ApiCustomsPage() {
  return (
    <MinimalListPage
      title="报关API"
      queryKey={['admin-api-customs']}
      queryFn={() => client.get('/admin/api/system/api-customs')}
      apiBase="/admin/api/system/api-customs"
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
