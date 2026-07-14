import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CustomerAddressesPage() {
  return (
    <GenericListPage
      title="客户地址"
      queryKey={['admin-customer-addresses']}
      queryFn={() => client.get('/admin/api/customer-addresses')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'client_name', label: '客户' },
        { key: 'address', label: '地址' },
        { key: 'city', label: '城市' },
        { key: 'country', label: '国家' },
        { key: 'is_default', label: '默认' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
