import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function AreaGroupsPage() {
  return (
    <GenericListPage
      title="区域分组"
      queryKey={['admin-area-groups']}
      queryFn={() => client.get('/admin/api/area-groups')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'code', label: 'Code' },
        { key: 'description', label: 'Description' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
