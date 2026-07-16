import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PrintersSystemPage() {
  return (
    <GenericListPage title="系统打印机" queryKey={['admin-PrintersSystemPage']}
      queryFn={() => client.get('/admin/api/printers')}
      apiBase="/admin/api/printers"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'type', label: '类型' }, { key: 'ip', label: 'IP地址' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
