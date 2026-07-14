import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function AiExceptionsPage() {
  return (
    <GenericListPage
      title="AI异常"
      queryKey={['admin-AiExceptionsPage']}
      queryFn={() => client.get('/admin/api/ai-exceptions')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
