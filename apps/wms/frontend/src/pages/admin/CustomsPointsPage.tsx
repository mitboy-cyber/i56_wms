import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CustomsPointsPage() {
  return (
    <GenericListPage title="海关口岸" queryKey={['admin-customs-points']}
      queryFn={() => client.get('/admin/api/customs-points')}
      apiBase="/admin/api/customs-points"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'code', label: '代码' }, { key: 'port', label: '港口' },
        { key: 'country', label: '国家' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
