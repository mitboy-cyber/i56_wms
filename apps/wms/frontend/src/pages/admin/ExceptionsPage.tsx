import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ExceptionsPage() {
  return (
    <GenericListPage
      title="异常列表"
      queryKey={['admin-ExceptionsPage']}
      queryFn={() => client.get('/admin/api/exceptions')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
