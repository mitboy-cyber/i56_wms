import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function RouteTemplatesPage() {
  return (
    <GenericListPage
      title="路线模板"
      queryKey={['admin-route-templates']}
      queryFn={() => client.get('/admin/api/route-templates')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'origin', label: '起运地' },
        { key: 'destination', label: '目的地' },
        { key: 'transport_mode', label: '运输方式' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
