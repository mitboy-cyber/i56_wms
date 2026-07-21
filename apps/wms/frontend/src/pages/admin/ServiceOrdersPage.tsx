import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function ServiceOrdersPage() {
  return (
    <MinimalListPage title="附加服务订单" queryKey={['admin-service-orders']}
      queryFn={() => client.get('/admin/api/service-order-records')}
      apiBase="/admin/api/service-order-records"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '服务名称' },
        { key: 'service_type', label: '服务类型' }, { key: 'price', label: '价格' },
        { key: 'description', label: '描述' }, { key: 'status', label: '状态' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
