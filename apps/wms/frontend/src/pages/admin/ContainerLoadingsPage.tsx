import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ContainerLoadingsPage() {
  return (
    <GenericListPage title="装箱记录" queryKey={['admin-container-loadings']}
      queryFn={() => client.get('/admin/api/container-loadings')}
      apiBase="/admin/api/container-loadings"
      columns={[
        { key: 'id', label: '编号' }, { key: 'container_no', label: '柜号' },
        { key: 'vessel', label: '船名' }, { key: 'port_from', label: '起运港' },
        { key: 'port_to', label: '目的港' }, { key: 'parcel_count', label: '件数' },
        { key: 'loaded_at', label: '装箱时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
