import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function ApiStoragePage() {
  return (
    <MinimalListPage
      title="仓储API"
      queryKey={['admin-api-storage']}
      queryFn={() => client.get('/admin/api/system/api-storage')}
      apiBase="/admin/api/system/api-storage"
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
