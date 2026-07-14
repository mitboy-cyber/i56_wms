import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function MonthlyStatementsPage() {
  return (
    <GenericListPage
      title="月结账单"
      queryKey={['admin-monthly-statements']}
      queryFn={() => client.get('/admin/api/monthly-statements')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'client_name', label: '客户' },
        { key: 'period', label: '账期' },
        { key: 'total_amount', label: '总金额' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
