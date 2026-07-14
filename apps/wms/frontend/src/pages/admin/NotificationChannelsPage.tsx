import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function NotificationChannelsPage() {
  return (
    <GenericListPage
      title="通知渠道"
      queryKey={['admin-NotificationChannelsPage']}
      queryFn={() => client.get('/admin/api/system/notification-channels')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
