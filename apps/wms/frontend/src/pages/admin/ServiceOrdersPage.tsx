import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ServiceOrdersPage() {
  return (
    <GenericListPage title="附加服务订单" queryKey={['admin-service-orders']}
      queryFn={() => client.get('/admin/api/service-orders')}
      apiBase="/admin/api/service-orders"
      columns={[
        { key: 'id', label: '编号' }, { key: 'order_no', label: '订单号' },
        { key: 'service_name', label: '服务项目' }, { key: 'client_name', label: '客户' },
        { key: 'amount', label: '金额' }, { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
