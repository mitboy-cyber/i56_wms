import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientPricingPage() {
  return (
    <GenericListPage
      title="客户报价"
      queryKey={['admin-client-pricing']}
      queryFn={() => client.get('/admin/api/client-pricing')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'client_id', label: 'Client Id' },
        { key: 'route_id', label: 'Route Id' },
        { key: 'price', label: 'Price' },
        { key: 'discount', label: 'Discount' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
