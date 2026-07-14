import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ShippingProvidersPage() {
  return (
    <GenericListPage
      title="承运商管理"
      queryKey={['admin-shipping-providers']}
      queryFn={() => client.get('/admin/api/shipping-providers')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'code', label: '编码' },
        { key: 'contact', label: '联系人' },
        { key: 'phone', label: '电话' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
