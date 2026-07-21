import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function CustomerAddressesPage() {
  return (
    <MinimalListPage title="客户地址" queryKey={['admin-customer-addresses']}
      queryFn={() => client.get('/admin/api/customer-addresses')}
      apiBase="/admin/api/customer-addresses"
      columns={[
        { key: 'id', label: '编号' }, { key: 'recipient_name', label: '收件人' },
        { key: 'phone', label: '电话' }, { key: 'city', label: '城市' },
        { key: 'district', label: '区域' }, { key: 'address', label: '详细地址' },
        { key: 'is_default', label: '默认' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
