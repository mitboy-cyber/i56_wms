import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';

export default function ClientMembers() {
  return (
    <GenericListPage title="会员管理" queryKey={['client-members']}
      queryFn={() => clientApi.members()}
      apiBase="/client/api/members"
      columns={[
        { key: 'name', label: '姓名' },
        { key: 'phone', label: '电话' },
        { key: 'email', label: '邮箱' },
        { key: 'member_code', label: '会员编号' },
      ]}
      getRowId={(r:any)=>String(r.id)} />
  );
}
