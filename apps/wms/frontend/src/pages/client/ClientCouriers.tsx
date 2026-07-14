import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientCouriers() {
  return <GenericListPage title="快递公司" queryKey={['client-couriers']} queryFn={() => clientApi.couriers()}
    columns={[{key:'name',label:'名称'},{key:'code',label:'代码'}]}
    getRowId={(r:Record<string,unknown>)=>String(r.code)} />;
}
