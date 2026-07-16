import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function RechargeRecordsPage() {
  return (
    <GenericListPage title="充值记录" queryKey={['admin-RechargeRecs']}
      queryFn={() => client.get('/admin/api/recharge-records')}
      apiBase="/admin/api/recharge-records"
      columns={[
        { key: 'id', label: '编号' }, { key: 'client_id', label: '客户' },
        { key: 'amount', label: '充值金额' }, { key: 'method', label: '支付方式' },
        { key: 'remark', label: '备注' }, { key: 'time', label: '时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
