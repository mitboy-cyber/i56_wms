import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function OrdersPage() {
  return (
    <GenericListPage
      title="订单管理"
      queryKey={['admin-orders']}
      queryFn={() => client.get('/admin/api/orders')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'order_no', label: 'Order No' },
        { key: 'recipient_name', label: 'Recipient Name' },
        { key: 'parcel_count', label: 'Parcel Count' },
        { key: 'total_price', label: 'Total Price' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
