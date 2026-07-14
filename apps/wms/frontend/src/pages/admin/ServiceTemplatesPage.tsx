import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ServiceTemplatesPage() {
  return (
    <GenericListPage
      title="服务模板"
      queryKey={['admin-ServiceTemplatesPage']}
      queryFn={() => client.get('/admin/api/service-templates')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
