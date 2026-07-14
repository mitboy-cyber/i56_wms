import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientOrders() {
  return <GenericListPage title="集运订单" queryKey={['client-orders']} queryFn={() => clientApi.orders()}
    columns={[{key:'order_no',label:'订单号'},{key:'recipient_name',label:'收件人'},{key:'parcel_count',label:'包裹数'},{key:'total_price',label:'金额',render:(v:unknown)=>`¥${Number(v).toFixed(2)}`},{key:'status',label:'状态'}]}
    getRowId={(r:Record<string,unknown>)=>r.order_no as string} />;
}
