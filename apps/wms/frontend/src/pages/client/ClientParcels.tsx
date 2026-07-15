import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';

export default function ClientParcels() {
  return (
    <GenericListPage title="我的包裹" queryKey={['client-parcels']}
      queryFn={() => clientApi.parcels()}
      apiBase="/client/api/parcels"
      columns={[
        { key: 'tracking_number', label: '快递单号' },
        { key: 'product_name', label: '品名' },
        { key: 'actual_weight', label: '重量(kg)' },
        { key: 'courier_code', label: '快递' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r:any)=>String(r.id)} />
  );
}
