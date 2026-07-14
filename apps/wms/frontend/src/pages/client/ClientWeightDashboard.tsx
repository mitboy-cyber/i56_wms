import GenericListPage from '@/components/GenericListPage';
import clientApi from '@/api/clientApi';
export default function ClientWeightDashboard() {
  return <GenericListPage title="称重看板" queryKey={['client-weight']} queryFn={() => clientApi.deliveryFees()}
    columns={[{key:'id',label:'#'}]} getRowId={(_:any,i:number)=>String(i)} />;
}
