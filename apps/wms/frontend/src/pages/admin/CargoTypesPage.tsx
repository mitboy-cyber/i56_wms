import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function CargoTypesPage() {
  return (
    <MinimalListPage title="货物类型" queryKey={['admin-cargo-types']}
      queryFn={() => client.get('/admin/api/cargo-types')}
      apiBase="/admin/api/cargo-types"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'code', label: '编码' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
