import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function BalanceLogsPage() {
  return (
    <GenericListPage
      title="余额日志"
      queryKey={['admin-balance-logs']}
      queryFn={() => client.get('/admin/api/balance-logs')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'client_name', label: '客户' },
        { key: 'change_amount', label: '变动金额' },
        { key: 'balance_after', label: '变动后余额' },
        { key: 'type', label: '类型' },
        { key: 'created_at', label: '时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
