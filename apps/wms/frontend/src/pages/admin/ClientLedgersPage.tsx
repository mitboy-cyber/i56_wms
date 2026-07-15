import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ClientLedgersPage() {
  return (
    <GenericListPage title="客户账本" queryKey={['admin-ClientLedgersPage']}
      queryFn={() => client.get('/admin/api/client-ledgers')}
      apiBase="/admin/api/client-ledgers"
      columns={[
        { key: 'id', label: '编号' }, { key: 'client_name', label: '客户名称' },
        { key: 'period', label: '账期' }, { key: 'balance', label: '余额' },
        { key: 'total_charged', label: '累计扣费' }, { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
