import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function CarriersPage() {
  return (
    <GenericListPage title="承运商列表" queryKey={['admin-carriers']}
      queryFn={() => client.get('/admin/api/carriers')}
      apiBase="/admin/api/carriers"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'code', label: '编码' }, { key: 'contact', label: '联系人' },
        { key: 'phone', label: '电话' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
