import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function NotificationsPage() {
  return (
    <GenericListPage
      title="通知管理"
      queryKey={['admin-notifications']}
      queryFn={() => client.get('/admin/api/notifications')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'title', label: 'Title' },
        { key: 'content', label: 'Content' },
        { key: 'channel', label: 'Channel' },
        { key: 'recipient', label: 'Recipient' },
        { key: 'sent', label: 'Sent' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
