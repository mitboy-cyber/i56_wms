import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';

export default function ClientServices() {
  return (
    <GenericListPage title="附加服务" queryKey={['client-services']}
      queryFn={() => clientApi.serviceOrders()}
      apiBase="/client/api/service-orders"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'service_name', label: '服务名称' },
        { key: 'price', label: '价格', render: (v:unknown)=>`¥${Number(v).toFixed(2)}` },
        { key: 'description', label: '描述' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r:any)=>String(r.id)} />
  );
}
