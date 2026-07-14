import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PDASessionsPage() {
  return (
    <GenericListPage
      title="PDA会话"
      queryKey={['admin-PDASessionsPage']}
      queryFn={() => client.get('/admin/api/pda-sessions')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
