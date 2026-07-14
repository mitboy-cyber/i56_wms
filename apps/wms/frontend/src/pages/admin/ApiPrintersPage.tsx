import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiPrintersPage() {
  return (
    <GenericListPage
      title="打印机API"
      queryKey={['admin-api-printers']}
      queryFn={() => client.get('/admin/api/api-printers')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'type', label: '类型' },
        { key: 'ip_address', label: 'IP地址' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
