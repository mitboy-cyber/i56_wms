import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PDASessionsPage() {
  return (
    <GenericListPage
      title="PDA会话"
      queryKey={['admin-PDASessionsPage']}
      queryFn={() => client.get('/admin/api/pda-sessions')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'operator_id', label: 'Operator Id' },
        { key: 'device', label: 'Device' },
        { key: 'login_at', label: 'Login At' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
