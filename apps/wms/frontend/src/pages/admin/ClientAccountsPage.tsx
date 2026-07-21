import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function ClientAccountsPage() {
  return (
    <MinimalListPage title="客户账户" queryKey={['admin-client-accounts']}
      queryFn={() => client.get('/admin/api/client-accounts')}
      apiBase="/admin/api/client-accounts"
      columns={[
        { key: 'id', label: '编号' }, { key: 'username', label: '用户名' },
        { key: 'real_name', label: '企业名称' }, { key: 'email', label: '邮箱' },
        { key: 'balance', label: '余额' }, { key: 'status', label: '状态' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
