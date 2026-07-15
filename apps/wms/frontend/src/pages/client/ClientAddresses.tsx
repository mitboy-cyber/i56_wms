import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';

export default function ClientAddresses() {
  return (
    <GenericListPage title="地址管理" queryKey={['client-addresses']}
      queryFn={() => clientApi.addresses()}
      apiBase="/client/api/addresses"
      columns={[
        { key: 'recipient_name', label: '收件人' },
        { key: 'phone', label: '电话' },
        { key: 'city', label: '城市' },
        { key: 'district', label: '区域' },
        { key: 'address', label: '详细地址' },
      ]}
      getRowId={(r:any)=>String(r.id)} />
  );
}
