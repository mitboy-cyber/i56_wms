import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function BalanceLogsPage() {
  return (
    <GenericListPage
      title="余额日志"
      queryKey={['admin-balance-logs']}
      queryFn={() => client.get('/admin/api/balance-logs')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'client_id', label: 'Client Id' },
        { key: 'change_amount', label: 'Change Amount' },
        { key: 'balance_after', label: 'Balance After' },
        { key: 'remark', label: 'Remark' },
        { key: 'time', label: 'Time' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
