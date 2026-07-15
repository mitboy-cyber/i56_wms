import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientRechargesPage() {
  return (
    <GenericListPage
      title="充值记录"
      queryKey={['admin-client-recharges']}
      queryFn={() => client.get('/admin/api/client-recharges')}
      apiBase="/admin/api/client-recharges"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'client_id', label: '客户编号' },
        { key: 'amount', label: '金额' },
        { key: 'method', label: '方式' },
        { key: 'remark', label: '备注' },
        { key: 'time', label: '时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
