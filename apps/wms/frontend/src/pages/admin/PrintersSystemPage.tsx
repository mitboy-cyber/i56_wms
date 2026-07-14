import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PrintersSystemPage() {
  return (
    <GenericListPage
      title="系统打印机"
      queryKey={['admin-PrintersSystemPage']}
      queryFn={() => client.get('/admin/api/system/printers')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
