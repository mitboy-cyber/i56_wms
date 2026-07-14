import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function RouteTemplatesPage() {
  return (
    <GenericListPage
      title="路线模板"
      queryKey={['admin-route-templates']}
      queryFn={() => client.get('/admin/api/route-templates')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'from', label: 'From' },
        { key: 'to', label: 'To' },
        { key: 'carrier_id', label: 'Carrier Id' },
        { key: 'est_days', label: 'Est Days' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
