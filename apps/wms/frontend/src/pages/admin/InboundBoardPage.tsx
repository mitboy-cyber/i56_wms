import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function InboundBoardPage() {
  return (
    <GenericListPage
      title="入库看板"
      queryKey={['admin-InboundBoardPage']}
      queryFn={() => client.get('/admin/api/inbound-board')}
      columns={[{ key: 'id', label: 'ID' }, { key: 'name', label: '名称' }]}
      getRowId={(_, i) => String(i)}
    />
  );
}
