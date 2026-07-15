import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function CouriersPage() {
  return (
    <GenericListPage title="快递公司" queryKey={['admin-couriers']}
      queryFn={() => client.get('/admin/api/couriers')}
      apiBase="/admin/api/couriers"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'code', label: '编码' }, { key: 'contact', label: '联系人' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
