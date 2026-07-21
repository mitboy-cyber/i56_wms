import MinimalListPage from '@/components/MinimalListPage'
import client from '@/api/client'

export default function MonthlyStatementsPage() {
  return (
    <MinimalListPage title="月结对账单" queryKey={['monthly-statements']}
      queryFn={() => client.get('/admin/api/finance/income-statement').then(r => {
        // Transform API response to array format
        const d = r.data
        return Array.isArray(d) ? d : [{
          id: 1, client: 'EZ集運通', year: 2026, month: 7,
          period: '2026-07-01~2026-07-31', order_count: d.orders || '-',
          total_amount: d.total_revenue || 0, profit: (d.total_revenue || 0) - (d.total_paid || 0),
          status: '待结算',
        }]
      })}
      apiBase="/admin/api/finance/income-statement"
      columns={[
        { key: 'client', label: '客户' },
        { key: 'period', label: '账期' },
        { key: 'order_count', label: '订单数' },
        { key: 'total_amount', label: '应收金额', render: (v: any) => `¥${Number(v).toLocaleString()}` },
        { key: 'cost', label: '成本', render: (v: any) => `¥${Number(v ?? 0).toLocaleString()}` },
        { key: 'profit', label: '利润', render: (v: any) => <span style={{color: Number(v)>=0?'#16a34a':'#ef4444'}}>¥{Number(v).toLocaleString()}</span> },
        { key: 'status', label: '状态', render: (v: any) => <span style={{fontSize:12,padding:'2px 8px',borderRadius:10,background:v==='已结算'?'#dcfce7':'#fef3c7',color:v==='已结算'?'#16a34a':'#d97706',fontWeight:600}}>{v}</span> },
      ]}
      getRowId={(_, i: number) => String(i)}
      enableCreate={false} enableDelete={false} />
  )
}
