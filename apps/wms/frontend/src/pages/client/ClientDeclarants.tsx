import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';

export default function ClientDeclarants() {
  return (
    <GenericListPage title="申报人管理" queryKey={['client-declarants']}
      queryFn={() => clientApi.declarants()}
      apiBase="/client/api/declarants"
      columns={[
        { key: 'name', label: '姓名' },
        { key: 'id_number', label: '身份证号' },
        { key: 'phone', label: '电话' },
        { key: 'auth_status', label: '认证状态' },
      ]}
      getRowId={(r:any)=>String(r.id)} />
  );
}
