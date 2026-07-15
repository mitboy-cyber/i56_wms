import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientPricingPage() {
  return (
    <GenericListPage title="客户报价" queryKey={['admin-client-pricing']}
      queryFn={() => client.get('/admin/api/client-pricing')}
      apiBase="/admin/api/client-pricing"
      columns={[
        { key: 'id', label: '编号' }, { key: 'client_name', label: '客户' },
        { key: 'route_name', label: '线路' }, { key: 'price', label: '报价' },
        { key: 'discount', label: '折扣' }, { key: 'status', label: '状态' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
