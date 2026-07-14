import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function AIChatPage() {
  return (
    <GenericListPage
      title="AI对话"
      queryKey={['admin-ai-chat']}
      queryFn={() => client.get('/admin/api/system/ai-chat')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'role', label: 'Role' },
        { key: 'content', label: 'Content' },
        { key: 'time', label: 'Time' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
