import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function OrderProfitReport() {
  return (
    <GenericListPage
      title="订单利润报表"
      queryKey={['admin-order-profit']}
      queryFn={() => client.get('/admin/api/report/order-profit')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'type', label: 'Type' },
        { key: 'status', label: 'Status' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
