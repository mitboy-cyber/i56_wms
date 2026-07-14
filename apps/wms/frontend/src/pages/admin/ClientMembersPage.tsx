import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientMembersPage() {
  return (
    <GenericListPage
      title="客户成员"
      queryKey={['admin-client-members']}
      queryFn={() => client.get('/admin/api/client-members')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'client_name', label: '客户' },
        { key: 'username', label: '用户名' },
        { key: 'role', label: '角色' },
        { key: 'phone', label: '电话' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
