import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ServiceWorkordersPage() {
  return (
    <GenericListPage
      title="服务工单"
      queryKey={['admin-ServiceWorkordersPage']}
      queryFn={() => client.get('/admin/api/service-workorders')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
