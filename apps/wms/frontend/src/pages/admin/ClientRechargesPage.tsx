import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientRechargesPage() {
  return (
    <GenericListPage
      title="充值记录"
      queryKey={['admin-client-recharges']}
      queryFn={() => client.get('/admin/api/client-recharges')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'client_id', label: 'Client Id' },
        { key: 'amount', label: 'Amount' },
        { key: 'method', label: 'Method' },
        { key: 'remark', label: 'Remark' },
        { key: 'time', label: 'Time' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
