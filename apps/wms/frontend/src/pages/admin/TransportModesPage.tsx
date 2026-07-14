import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function TransportModesPage() {
  return (
    <GenericListPage
      title="运输方式"
      queryKey={['admin-transport-modes']}
      queryFn={() => client.get('/admin/api/transport-modes')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'code', label: 'Code' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
