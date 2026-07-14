import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function OrderProfitReport() {
  return (
    <GenericListPage
      title="订单利润报表"
      queryKey={['admin-order-profit']}
      queryFn={() => client.get('/admin/api/order-profit')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'order_no', label: '订单号' },
        { key: 'client_name', label: '客户' },
        { key: 'revenue', label: '收入' },
        { key: 'cost', label: '成本' },
        { key: 'profit', label: '利润' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
