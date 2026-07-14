import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function OrdersPage() {
  return (
    <GenericListPage
      title="订单管理"
      queryKey={['admin-orders']}
      queryFn={() => client.get('/admin/api/orders')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'order_no', label: '订单号' },
        { key: 'client_name', label: '客户' },
        { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
