import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PricingRoutesPage() {
  return (
    <GenericListPage
      title="路线报价"
      queryKey={['admin-pricing-routes']}
      queryFn={() => client.get('/admin/api/pricing/routes')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'origin', label: '起运地' },
        { key: 'destination', label: '目的地' },
        { key: 'transport_mode', label: '运输方式' },
        { key: 'base_price', label: '基础价格' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
