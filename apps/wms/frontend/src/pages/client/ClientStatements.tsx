import GenericListPage from '@/components/GenericListPage';
import { clientApi } from '@/api/client';

export default function ClientStatements() {
  return (
    <GenericListPage title="月结对账单" queryKey={['client-statements']}
      queryFn={() => clientApi.get('/client/api/statements')}
      columns={[
        { key: 'id', label: '编号' },
        { key: 'period', label: '期间' },
        { key: 'open_balance', label: '期初余额', render: (v: unknown) => <span>¥{Number(v).toLocaleString()}</span> },
        { key: 'credit', label: '入账(充值)', render: (v: unknown) => <span className="text-green-600">+¥{Number(v).toLocaleString()}</span> },
        { key: 'debit', label: '扣账(订单)', render: (v: unknown) => <span className="text-red-500">-¥{Number(v).toLocaleString()}</span> },
        { key: 'close_balance', label: '期末余额', render: (v: unknown) => <span className="font-medium">¥{Number(v).toLocaleString()}</span> },
        { key: 'txn_count', label: '流水数' },
        { key: 'status', label: '状态' },
        { key: 'generated_at', label: '生成时间' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
