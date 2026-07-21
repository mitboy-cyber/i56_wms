import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function BalanceLogsPage() {
  return (
    <MinimalListPage title="余额日志" queryKey={['admin-BalanceLogs']}
      queryFn={() => client.get('/admin/api/balance-logs')}
      apiBase="/admin/api/balance-logs"
      columns={[
        { key: 'id', label: '编号' }, { key: 'client_id', label: '客户' },
        { key: 'type', label: '类型' }, { key: 'amount', label: '金额' },
        { key: 'balance', label: '余额' }, { key: 'note', label: '备注' },
        { key: 'created_at', label: '时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
