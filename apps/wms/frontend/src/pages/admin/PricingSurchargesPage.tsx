import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PricingSurchargesPage() {
  return (
    <GenericListPage
      title="附加费报价"
      queryKey={['admin-pricing-surcharges']}
      queryFn={() => client.get('/admin/api/pricing/surcharges')}
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
