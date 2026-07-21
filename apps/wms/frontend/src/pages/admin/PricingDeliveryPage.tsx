import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function PricingDeliveryPage() {
  return (
    <MinimalListPage title="派送报价" queryKey={['admin-pricing-delivery']}
      queryFn={() => client.get('/admin/api/pricing/delivery')}
      apiBase="/admin/api/pricing/delivery"
      columns={[
        { key: 'carrier', label: '承运商' }, { key: 'customs_point', label: '清关点' },
        { key: 'area', label: '派送区域' }, { key: 'delivery_type', label: '配送方式' },
        { key: 'condition', label: '计价条件' }, { key: 'price', label: '价格' },
        { key: 'unit', label: '单位' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
