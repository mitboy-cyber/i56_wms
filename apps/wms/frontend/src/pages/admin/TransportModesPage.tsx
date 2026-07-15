import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function TransportModesPage() {
  return (
    <GenericListPage
      title="运输方式"
      queryKey={['admin-transport-modes']}
      queryFn={() => client.get('/admin/api/transport-modes')}
      apiBase="/admin/api/transport-modes"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'name', label: '名称' },
        { key: 'code', label: '编码' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
