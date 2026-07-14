import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ContainerLoadingsPage() {
  return (
    <GenericListPage
      title="集装箱装货"
      queryKey={['admin-container-loadings']}
      queryFn={() => client.get('/admin/api/container-loadings')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'container_no', label: '箱号' },
        { key: 'order_no', label: '订单号' },
        { key: 'seal_no', label: '封条号' },
        { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
