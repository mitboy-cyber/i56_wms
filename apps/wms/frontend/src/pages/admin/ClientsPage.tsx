import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function ClientsPage() {
  return (
    <MinimalListPage title="客户管理" queryKey={['admin-clients']}
      queryFn={() => client.get('/admin/api/clients')}
      apiBase="/admin/api/clients"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'code', label: '编码' }, { key: 'type', label: '类型' },
        { key: 'contact', label: '联系人' }, { key: 'phone', label: '电话' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
