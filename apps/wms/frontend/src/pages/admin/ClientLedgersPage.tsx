import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ClientLedgersPage() {
  return (
    <GenericListPage
      title="客户账本"
      queryKey={['admin-ClientLedgersPage']}
      queryFn={() => client.get('/admin/api/client-ledgers')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
