import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PricingRoutesPage() {
  return (
    <GenericListPage title="路线报价" queryKey={['admin-pricing-routes']}
      queryFn={() => client.get('/admin/api/pricing/routes')}
      apiBase="/admin/api/pricing/routes"
      columns={[
        { key: 'route_name', label: '路线名称' }, { key: 'transport_type', label: '运输方式' },
        { key: 'cargo_type', label: '货物类型' }, { key: 'tax_type', label: '计税方式' },
        { key: 'first_weight', label: '首重(kg)' }, { key: 'first_weight_price', label: '首重价格' },
        { key: 'continuation_weight', label: '续重(kg)' }, { key: 'continuation_price', label: '续重价格' },
        { key: 'min_weight', label: '最低重量' }, { key: 'max_weight', label: '最高重量' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
