import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PDAWorkorderTemplatesPage() {
  return (
    <GenericListPage
      title="PDA工单模板"
      queryKey={['admin-PDAWorkorderTemplatesPage']}
      queryFn={() => client.get('/admin/api/pda-workorder-templates')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'process_type', label: 'Process Type' },
        { key: 'steps', label: 'Steps' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
