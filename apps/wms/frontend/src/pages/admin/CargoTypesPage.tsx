import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CargoTypesPage() {
  return (
    <GenericListPage
      title="货物类型"
      queryKey={['admin-cargo-types']}
      queryFn={() => client.get('/admin/api/cargo-types')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'code', label: 'Code' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
