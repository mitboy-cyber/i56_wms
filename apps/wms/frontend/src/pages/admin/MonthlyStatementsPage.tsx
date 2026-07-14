import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function MonthlyStatementsPage() {
  return (
    <GenericListPage
      title="月结账单"
      queryKey={['admin-monthly-statements']}
      queryFn={() => client.get('/admin/api/monthly-statements')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'client_id', label: 'Client Id' },
        { key: 'period', label: 'Period' },
        { key: 'total', label: 'Total' },
        { key: 'paid_amount', label: 'Paid Amount' },
        { key: 'status', label: 'Status' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
