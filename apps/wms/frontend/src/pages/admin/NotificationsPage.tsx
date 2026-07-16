import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function NotificationsPage() {
  return (
    <GenericListPage title="通知管理" queryKey={['admin-notifications']}
      queryFn={() => client.get('/admin/api/notifications')}
      apiBase="/admin/api/notifications"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'title', label: '标题' },
        { key: 'type', label: '类型' },
        { key: 'priority', label: '优先级', render: (v: unknown) => <span className={String(v) === '紧急' ? 'text-red-600 font-medium' : 'text-gray-600'}>{String(v)}</span> },
        { key: 'scope', label: '发送范围' },
        { key: 'content', label: '内容', render: (v: unknown) => <span className="text-xs text-gray-500 max-w-xs truncate block">{String(v)}</span> },
        { key: 'channel', label: '渠道' },
        { key: 'sender_name', label: '发送人' },
        { key: 'sent', label: '已发送', render: (v: unknown) => <span>{v ? '✅' : '⏳'}</span> },
        { key: 'created_at', label: '发送时间', render: (v: unknown) => <span className="text-xs">{v ? new Date(String(v)).toLocaleString('zh-CN') : '—'}</span> },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
