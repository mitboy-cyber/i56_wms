import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';

export default function ClientRoutePrices() {
  return (
    <GenericListPage title="路线价格" queryKey={['client-routes']}
      queryFn={() => clientApi.routePrices()}
      apiBase="/client/api/route-prices"
      columns={[
        { key: 'route_name', label: '线路' },
        { key: 'transport_type', label: '运输方式' },
        { key: 'base_weight_price', label: '重量单价/kg', render: (v:unknown)=>`¥${Number(v).toFixed(2)}` },
        { key: 'base_volume_price', label: '体积单价/kg', render: (v:unknown)=>`¥${Number(v).toFixed(2)}` },
      ]}
      getRowId={(r:any,i:number)=>String(r.id||i)} />
  );
}
