import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ShippingProvidersPage() {
  return (
    <GenericListPage
      title="承运商管理"
      queryKey={['admin-shipping-providers']}
      queryFn={() => client.get('/admin/api/shipping-providers')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'code', label: 'Code' },
        { key: 'contact', label: 'Contact' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
