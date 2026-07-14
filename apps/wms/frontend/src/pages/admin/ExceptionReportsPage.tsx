import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ExceptionReportsPage() {
  return (
    <GenericListPage
      title="异常报告"
      queryKey={['admin-ExceptionReportsPage']}
      queryFn={() => client.get('/admin/api/exception-reports')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'type', label: 'Type' },
        { key: 'count', label: 'Count' },
        { key: 'period', label: 'Period' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
