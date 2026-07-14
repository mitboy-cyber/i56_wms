import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientRoutePrices() {
  return <GenericListPage title="线路报价" queryKey={['client-route-prices']} queryFn={() => clientApi.routePrices()}
    columns={[{key:'route_name',label:'线路'},{key:'transport_type',label:'运输方式'},{key:'base_weight_price',label:'重量单价',render:(v:unknown)=>'¥'+Number(v).toFixed(2)+'/kg'},{key:'base_volume_price',label:'体积单价'}]}
    getRowId={(_:unknown,i:number)=>String(i)} />;
}
