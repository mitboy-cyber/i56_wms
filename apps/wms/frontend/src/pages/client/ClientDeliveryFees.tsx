import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientDeliveryFees() {
  return <GenericListPage title="派送费" queryKey={['client-delivery-fees']} queryFn={() => clientApi.deliveryFees()}
    columns={[{key:'id',label:'#'}]}
    getRowId={(_:unknown,i:number)=>String(i)} />;
}
