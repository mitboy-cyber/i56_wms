import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ClientLedgersPage() {
  return (
    <GenericListPage title="客户账本" queryKey={['admin-ClientLedgersPage']}
      queryFn={() => client.get('/admin/api/client-ledgers')}
      apiBase="/admin/api/client-ledgers"
      columns={[
        { key: 'id', label: '编号' }, { key: 'amount', label: '金额' },
        { key: 'balance_after', label: '余额' }, { key: 'type', label: '类型' },
        { key: 'description', label: '描述' }, { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
