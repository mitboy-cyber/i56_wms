import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientAccountsPage() {
  return (
    <GenericListPage
      title="客户账户"
      queryKey={['admin-client-accounts']}
      queryFn={() => client.get('/admin/api/client-accounts')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'username', label: 'Username' },
        { key: 'real_name', label: 'Real Name' },
        { key: 'email', label: 'Email' },
        { key: 'balance', label: 'Balance' },
        { key: 'status', label: 'Status' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
