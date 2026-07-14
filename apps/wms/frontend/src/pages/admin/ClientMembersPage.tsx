import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientMembersPage() {
  return (
    <GenericListPage
      title="客户成员"
      queryKey={['admin-client-members']}
      queryFn={() => client.get('/admin/api/client-members')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'client_id', label: 'Client Id' },
        { key: 'name', label: 'Name' },
        { key: 'phone', label: 'Phone' },
        { key: 'email', label: 'Email' },
        { key: 'role', label: 'Role' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
