import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CustomsPointsPage() {
  return (
    <GenericListPage
      title="海关口岸"
      queryKey={['admin-customs-points']}
      queryFn={() => client.get('/admin/api/customs-points')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'code', label: 'Code' },
        { key: 'port', label: 'Port' },
        { key: 'country', label: 'Country' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
