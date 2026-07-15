import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function NotificationChannelsPage() {
  return (
    <GenericListPage title="通知渠道" queryKey={['admin-NotificationChannelsPage']}
      queryFn={() => client.get('/admin/api/system/notification-channels')}
      apiBase="/admin/api/system/notification-channels"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '渠道名称' },
        { key: 'type', label: '类型' }, { key: 'enabled', label: '启用' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
