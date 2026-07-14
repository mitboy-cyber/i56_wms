import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PrintersPage() {
  return (
    <GenericListPage
      title="打印机管理"
      queryKey={['admin-PrintersPage']}
      queryFn={() => client.get('/admin/api/printers')}
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
