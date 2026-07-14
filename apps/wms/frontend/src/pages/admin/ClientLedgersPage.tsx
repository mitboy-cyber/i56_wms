import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ClientLedgersPage() {
  return (
    <GenericListPage
      title="客户账本"
      queryKey={['admin-ClientLedgersPage']}
      queryFn={() => client.get('/admin/api/client-ledgers')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'balance', label: 'Balance' },
        { key: 'total_charged', label: 'Total Charged' },
        { key: 'period', label: 'Period' },
        { key: 'status', label: 'Status' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
