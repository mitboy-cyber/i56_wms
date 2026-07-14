import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function NotificationChannelsPage() {
  return (
    <GenericListPage
      title="通知渠道"
      queryKey={['admin-NotificationChannelsPage']}
      queryFn={() => client.get('/admin/api/system/notification-channels')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'type', label: 'Type' },
        { key: 'config', label: 'Config' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
