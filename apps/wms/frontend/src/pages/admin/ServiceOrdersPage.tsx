import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ServiceOrdersPage() {
  return (
    <GenericListPage
      title="服务订单"
      queryKey={['admin-service-orders']}
      queryFn={() => client.get('/admin/api/service-orders')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'order_no', label: 'Order No' },
        { key: 'type', label: 'Type' },
        { key: 'status', label: 'Status' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
