import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientWarehouses() {
  return <GenericListPage title="仓库" queryKey={['client-warehouses']} queryFn={() => clientApi.warehouses()}
    columns={[{key:'name',label:'名称'},{key:'code',label:'编码'},{key:'address',label:'地址'}]}
    getRowId={(r:Record<string,unknown>)=>String(r.id)} />;
}
