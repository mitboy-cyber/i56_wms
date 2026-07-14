import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ExceptionReportsPage() {
  return (
    <GenericListPage
      title="异常报告"
      queryKey={['admin-ExceptionReportsPage']}
      queryFn={() => client.get('/admin/api/exception-reports')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
