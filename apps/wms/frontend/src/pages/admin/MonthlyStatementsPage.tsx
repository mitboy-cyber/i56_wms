import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function MonthlyStatementsPage() {
  return (
    <MinimalListPage title="月结对账单" queryKey={['admin-monthly-statements']}
      queryFn={() => client.get('/admin/api/monthly-statements')}
      apiBase="/admin/api/monthly-statements"
      columns={[
        { key: 'id', label: '编号' }, { key: 'client_id', label: '客户编号' },
        { key: 'period', label: '期间', render: (v: unknown) => <span className="font-medium">{String(v)}</span> },
        { key: 'total', label: '账单金额', render: (v: unknown) => <span className="font-mono text-gray-800">¥{Number(v).toFixed(2)}</span> },
        { key: 'paid_amount', label: '已付金额', render: (v: unknown) => <span className="font-mono text-green-600">¥{Number(v).toFixed(2)}</span> },
        { key: 'status', label: '状态', render: (v: unknown) => <span className={String(v) === '已结算' ? 'text-green-600 font-medium' : 'text-amber-600 font-medium'}>{String(v)}</span> },
        { key: 'created_at', label: '生成时间', render: (v: unknown) => <span className="text-xs text-gray-500">{v ? new Date(String(v)).toLocaleString('zh-CN') : '—'}</span> },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
