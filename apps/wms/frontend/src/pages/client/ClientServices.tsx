import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientServices() {
  return <GenericListPage title="附加服务" queryKey={['client-services']} queryFn={() => clientApi.serviceOrders()}
    columns={[{key:'id',label:'ID'},{key:'service_type',label:'服务类型'},{key:'status',label:'状态'}]}
    getRowId={(r:Record<string,unknown>)=>String(r.id)} />;
}
