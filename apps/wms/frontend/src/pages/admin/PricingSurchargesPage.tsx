import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function PricingSurchargesPage() {
  return (
    <MinimalListPage title="附加费报价" queryKey={['admin-pricing-surcharges']}
      queryFn={() => client.get('/admin/api/pricing/surcharges')}
      apiBase="/admin/api/pricing/surcharges"
      columns={[
        { key: 'carrier', label: '承运商' }, { key: 'customs_point', label: '清关点' },
        { key: 'surcharge_name', label: '附加费名称' }, { key: 'condition', label: '触发条件' },
        { key: 'rule', label: '计算规则' }, { key: 'price', label: '金额' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
