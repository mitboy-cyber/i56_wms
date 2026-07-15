import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function MonthlyStatementsPage() {
  return (
    <GenericListPage title="月结账单" queryKey={['admin-monthly-statements']}
      queryFn={() => client.get('/admin/api/monthly-statements')}
      apiBase="/admin/api/monthly-statements"
      columns={[
        { key: 'id', label: '编号' }, { key: 'period', label: '账期' },
        { key: 'total', label: '账单金额' }, { key: 'paid_amount', label: '已付金额' },
        { key: 'status', label: '状态' }, { key: 'created_at', label: '创建时间' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
