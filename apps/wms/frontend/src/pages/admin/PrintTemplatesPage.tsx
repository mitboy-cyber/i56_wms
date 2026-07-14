import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PrintTemplatesPage() {
  return (
    <GenericListPage
      title="打印模板"
      queryKey={['admin-print-templates']}
      queryFn={() => client.get('/admin/api/print-templates')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'type', label: 'Type' },
        { key: 'description', label: 'Description' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
