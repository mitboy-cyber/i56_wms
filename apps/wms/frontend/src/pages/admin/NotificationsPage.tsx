import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function NotificationsPage() {
  return (
    <GenericListPage
      title="通知管理"
      queryKey={['admin-notifications']}
      queryFn={() => client.get('/admin/api/notifications')}
      apiBase="/admin/api/notifications"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'title', label: '标题' },
        { key: 'content', label: '内容' },
        { key: 'channel', label: '渠道' },
        { key: 'recipient', label: '接收人' },
        { key: 'sent', label: '已发送' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
