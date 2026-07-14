import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function NotificationsPage() {
  return (
    <GenericListPage
      title="通知管理"
      queryKey={['admin-notifications']}
      queryFn={() => client.get('/admin/api/notifications')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'title', label: '标题' },
        { key: 'type', label: '类型' },
        { key: 'recipient', label: '接收人' },
        { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
