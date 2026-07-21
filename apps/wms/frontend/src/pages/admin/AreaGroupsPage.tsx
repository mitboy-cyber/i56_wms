import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function AreaGroupsPage() {
  return (
    <MinimalListPage title="区域分组" queryKey={['admin-area-groups']}
      queryFn={() => client.get('/admin/api/area-groups')}
      apiBase="/admin/api/area-groups"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'code', label: '代码' }, { key: 'description', label: '描述' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
