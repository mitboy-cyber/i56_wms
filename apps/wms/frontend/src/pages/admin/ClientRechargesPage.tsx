import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientRechargesPage() {
  return (
    <GenericListPage
      title="充值记录"
      queryKey={['admin-client-recharges']}
      queryFn={() => client.get('/admin/api/client-recharges')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'client_name', label: '客户' },
        { key: 'amount', label: '金额' },
        { key: 'method', label: '充值方式' },
        { key: 'status', label: '状态' },
        { key: 'created_at', label: '充值时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
