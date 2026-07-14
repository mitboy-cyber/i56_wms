import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function AIChatPage() {
  return (
    <GenericListPage
      title="AI对话"
      queryKey={['admin-ai-chat']}
      queryFn={() => client.get('/admin/api/ai-chat')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'session_id', label: '会话ID' },
        { key: 'user', label: '用户' },
        { key: 'message', label: '消息' },
        { key: 'created_at', label: '时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
