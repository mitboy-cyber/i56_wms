import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CustomerDeclarantsPage() {
  return (
    <GenericListPage
      title="客户报关员"
      queryKey={['admin-customer-declarants']}
      queryFn={() => client.get('/admin/api/customer-declarants')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'client_id', label: 'Client Id' },
        { key: 'name', label: 'Name' },
        { key: 'id_number', label: 'Id Number' },
        { key: 'phone', label: 'Phone' },
        { key: 'country', label: 'Country' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
