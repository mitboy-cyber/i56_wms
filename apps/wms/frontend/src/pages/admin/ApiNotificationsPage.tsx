import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiNotificationsPage() {
  return (
    <GenericListPage
      title="通知API"
      queryKey={['admin-api-notifications']}
      queryFn={() => client.get('/admin/api/api-notifications')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'channel', label: '渠道' },
        { key: 'config', label: '配置' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
