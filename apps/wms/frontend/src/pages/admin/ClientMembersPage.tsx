import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function ClientMembersPage() {
  return (
    <MinimalListPage title="客户会员" queryKey={['admin-ClientMembers']}
      queryFn={() => client.get('/admin/api/client-members')}
      apiBase="/admin/api/client-members"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '会员名称' },
        { key: 'phone', label: '电话' }, { key: 'id_number', label: '证件号' },
        { key: 'platform', label: '平台' }, { key: 'created_at', label: '创建时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
