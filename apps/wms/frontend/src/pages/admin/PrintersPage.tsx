import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PrintersPage() {
  return (
    <GenericListPage title="打印机管理" queryKey={['admin-printers']}
      queryFn={() => client.get('/admin/api/printers')}
      apiBase="/admin/api/printers"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'type', label: '类型' }, { key: 'ip', label: 'IP地址' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
