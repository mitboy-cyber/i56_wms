import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientPricingPage() {
  return (
    <GenericListPage
      title="客户报价"
      queryKey={['admin-client-pricing']}
      queryFn={() => client.get('/admin/api/client-pricing')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'client_name', label: '客户' },
        { key: 'route', label: '路线' },
        { key: 'price', label: '价格' },
        { key: 'effective_date', label: '生效日期' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
