import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientDeclarants() {
  return <GenericListPage title="申报人" queryKey={['client-declarants']} queryFn={() => clientApi.declarants()}
    columns={[{key:'name',label:'姓名'},{key:'id_number',label:'证件号'},{key:'type',label:'类型'},{key:'phone',label:'电话'}]}
    getRowId={(r:Record<string,unknown>)=>String(r.id)} />;
}
