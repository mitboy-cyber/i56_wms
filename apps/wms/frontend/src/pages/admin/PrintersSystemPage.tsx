import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PrintersSystemPage() {
  return (
    <GenericListPage
      title="系统打印机"
      queryKey={['admin-PrintersSystemPage']}
      queryFn={() => client.get('/admin/api/system/printers')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'type', label: 'Type' },
        { key: 'ip', label: 'Ip' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
