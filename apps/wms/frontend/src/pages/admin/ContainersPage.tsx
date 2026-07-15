import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ContainersPage() {
  return (
    <GenericListPage title="集装柜管理" queryKey={['admin-Containers']}
      queryFn={() => client.get('/admin/api/containers').then(r => r.data)}
      apiBase="/admin/api/containers"
      columns={[
        { key: 'id', label: '编号' }, { key: 'warehouse', label: '仓库' },
        { key: 'container_no', label: '柜号' }, { key: 'route_name', label: '线路' },
        { key: 'status', label: '状态' }, { key: 'max_weight', label: '限重(kg)' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
