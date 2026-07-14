import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CargoTypesPage() {
  return (
    <GenericListPage
      title="货物类型"
      queryKey={['admin-cargo-types']}
      queryFn={() => client.get('/admin/api/cargo-types')}
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
