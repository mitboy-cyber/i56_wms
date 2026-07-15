import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiNotificationsPage() {
  return (
    <GenericListPage
      title="通知API"
      queryKey={['admin-api-notifications']}
      queryFn={() => client.get('/admin/api/system/api-notifications')}
      apiBase="/admin/api/system/api-notifications"
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
