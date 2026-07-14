import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CustomerAddressesPage() {
  return (
    <GenericListPage
      title="客户地址"
      queryKey={['admin-customer-addresses']}
      queryFn={() => client.get('/admin/api/customer-addresses')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'recipient_name', label: 'Recipient Name' },
        { key: 'phone', label: 'Phone' },
        { key: 'city', label: 'City' },
        { key: 'address', label: 'Address' },
        { key: 'is_default', label: 'Is Default' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
