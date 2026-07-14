import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function PricingSurchargesPage() {
  return (
    <GenericListPage
      title="附加费报价"
      queryKey={['admin-pricing-surcharges']}
      queryFn={() => client.get('/admin/api/pricing-surcharges')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'type', label: '类型' },
        { key: 'amount', label: '金额' },
        { key: 'applicable_to', label: '适用范围' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
