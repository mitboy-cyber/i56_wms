import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ServiceOrdersPage() {
  return (
    <GenericListPage
      title="服务订单"
      queryKey={['admin-service-orders']}
      queryFn={() => client.get('/admin/api/service-orders')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'order_no', label: '订单号' },
        { key: 'service_type', label: '服务类型' },
        { key: 'client_name', label: '客户' },
        { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
