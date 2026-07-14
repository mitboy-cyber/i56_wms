import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CustomerDeclarantsPage() {
  return (
    <GenericListPage
      title="客户报关员"
      queryKey={['admin-customer-declarants']}
      queryFn={() => client.get('/admin/api/customer-declarants')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'client_name', label: '客户' },
        { key: 'name', label: '姓名' },
        { key: 'id_number', label: '证件号' },
        { key: 'phone', label: '电话' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
