import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientMembers() {
  return <GenericListPage title="会员管理" queryKey={['client-members']} queryFn={() => clientApi.members()}
    columns={[{key:'name',label:'名称'}]}
    getRowId={(r:Record<string,unknown>)=>String(r.id)} />;
}
