import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ServiceTemplatesPage() {
  return (
    <GenericListPage title="附加服务模板" queryKey={['admin-service-templates']}
      queryFn={() => client.get('/admin/api/service-templates')}
      apiBase="/admin/api/service-templates"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '模板名称' },
        { key: 'service_type', label: '服务类型' }, { key: 'price', label: '价格' },
        { key: 'description', label: '描述' }, { key: 'status', label: '状态' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
