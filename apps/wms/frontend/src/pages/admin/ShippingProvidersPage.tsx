import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ShippingProvidersPage() {
  return (
    <GenericListPage title="承运商列表" queryKey={['admin-shipping-providers']}
      queryFn={() => client.get('/admin/api/shipping-providers')}
      apiBase="/admin/api/shipping-providers"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '承运商名称' },
        { key: 'code', label: '编码' }, { key: 'contact', label: '联系人' },
        { key: 'phone', label: '电话' }, { key: 'status', label: '状态' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
