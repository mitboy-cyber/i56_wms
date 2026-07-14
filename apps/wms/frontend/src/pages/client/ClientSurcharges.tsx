import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientSurcharges() {
  return <GenericListPage title="附加费" queryKey={['client-surcharges']} queryFn={() => clientApi.surcharges()}
    columns={[{key:'id',label:'#'}]}
    getRowId={(_:unknown,i:number)=>String(i)} />;
}
