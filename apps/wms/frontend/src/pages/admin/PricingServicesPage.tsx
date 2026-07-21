import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function PricingServicesPage() {
  return (
    <MinimalListPage title="服务报价" queryKey={['admin-pricing-services']}
      queryFn={() => client.get('/admin/api/pricing/services')}
      apiBase="/admin/api/pricing/services"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'type', label: '类型' }, { key: 'price', label: '价格' },
        { key: 'description', label: '描述' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
