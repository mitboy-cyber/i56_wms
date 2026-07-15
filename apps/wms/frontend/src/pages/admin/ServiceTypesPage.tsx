import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ServiceTypesPage() {
  return (
    <GenericListPage
      title="服务类型"
      queryKey={['admin-ServiceTypesPage']}
      queryFn={() => client.get('/admin/api/service-types')}
      apiBase="/admin/api/service-types"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'name', label: '名称' },
        { key: 'code', label: '编码' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
