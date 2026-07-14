import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ServiceTemplatesPage() {
  return (
    <GenericListPage
      title="服务模板"
      queryKey={['admin-ServiceTemplatesPage']}
      queryFn={() => client.get('/admin/api/service-templates')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'type', label: 'Type' },
        { key: 'description', label: 'Description' },
        { key: 'fee', label: 'Fee' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
