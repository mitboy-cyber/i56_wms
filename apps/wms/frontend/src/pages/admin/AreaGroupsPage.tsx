import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function AreaGroupsPage() {
  return (
    <GenericListPage
      title="区域分组"
      queryKey={['admin-area-groups']}
      queryFn={() => client.get('/admin/api/area-groups')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'code', label: '编码' },
        { key: 'description', label: '描述' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
