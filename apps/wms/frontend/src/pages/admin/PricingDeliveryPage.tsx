import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PricingDeliveryPage() {
  return (
    <GenericListPage
      title="配送报价"
      queryKey={['admin-pricing-delivery']}
      queryFn={() => client.get('/admin/api/pricing/delivery')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'carrier_name', label: 'Carrier Name' },
        { key: 'route', label: 'Route' },
        { key: 'delivery_zone', label: 'Delivery Zone' },
        { key: 'price', label: 'Price' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
