import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientAccountsPage() {
  return (
    <GenericListPage
      title="客户账户"
      queryKey={['admin-client-accounts']}
      queryFn={() => client.get('/admin/api/client-accounts')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'client_name', label: '客户' },
        { key: 'balance', label: '余额' },
        { key: 'currency', label: '币种' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
