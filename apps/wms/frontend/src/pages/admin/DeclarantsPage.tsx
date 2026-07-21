import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function DeclarantsPage() {
  return (
    <MinimalListPage title="申报人" queryKey={['admin-declarants']}
      queryFn={() => client.get('/admin/api/customer-declarants')}
      apiBase="/admin/api/customer-declarants"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '申报人/公司名' },
        { key: 'phone', label: '手机号' }, { key: 'id_number', label: '身份证号' },
        { key: 'status', label: '认证状态' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
