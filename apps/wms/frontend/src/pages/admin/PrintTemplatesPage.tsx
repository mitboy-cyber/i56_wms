import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PrintTemplatesPage() {
  return (
    <GenericListPage title="打印模板" queryKey={['admin-print-templates']}
      queryFn={() => client.get('/admin/api/print-templates')}
      apiBase="/admin/api/print-templates"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '模板名称' },
        { key: 'type', label: '类型' }, { key: 'is_default', label: '默认' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
