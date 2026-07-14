import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PricingDeliveryPage() {
  return (
    <GenericListPage
      title="配送报价"
      queryKey={['admin-pricing-delivery']}
      queryFn={() => client.get('/admin/api/pricing-delivery')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'zone', label: '区域' },
        { key: 'weight_range', label: '重量区间' },
        { key: 'price', label: '价格' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
