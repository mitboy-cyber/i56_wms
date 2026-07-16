import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CustomerDeclarantsPage() {
  return (
    <GenericListPage title="客户报关员" queryKey={['admin-customer-declarants']}
      queryFn={() => client.get('/admin/api/customer-declarants')}
      apiBase="/admin/api/customer-declarants"
      columns={[
        { key: 'id', label: '编号' }, { key: 'client_id', label: '客户编号' },
        { key: 'name', label: '姓名' }, { key: 'id_number', label: '证件号' },
        { key: 'phone', label: '电话' }, { key: 'country', label: '国家/地区' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
