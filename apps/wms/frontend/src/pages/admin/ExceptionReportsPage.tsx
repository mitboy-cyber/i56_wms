import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function ExceptionReportsPage() {
  return (
    <MinimalListPage title="异常报告" queryKey={['admin-ExceptionReports']}
      queryFn={() => client.get('/admin/api/exception-reports')}
      apiBase="/admin/api/exception-reports"
      columns={[
        { key: 'id', label: '编号' }, { key: 'title', label: '异常标题' },
        { key: 'type', label: '类型' }, { key: 'count', label: '次数' },
        { key: 'last_occurred', label: '最后发生' }, { key: 'status', label: '处理状态' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
