import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PrintTemplatesPage() {
  return (
    <GenericListPage
      title="打印模板"
      queryKey={['admin-print-templates']}
      queryFn={() => client.get('/admin/api/print-templates')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'type', label: '类型' },
        { key: 'description', label: '描述' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
