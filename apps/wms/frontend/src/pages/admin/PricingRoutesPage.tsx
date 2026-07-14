import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PricingRoutesPage() {
  return (
    <GenericListPage
      title="路线报价"
      queryKey={['admin-pricing-routes']}
      queryFn={() => client.get('/admin/api/pricing/routes')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'type', label: 'Type' },
        { key: 'amount', label: 'Amount' },
        { key: 'applicable_to', label: 'Applicable To' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
