import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CustomsPointsPage() {
  return (
    <GenericListPage
      title="海关口岸"
      queryKey={['admin-customs-points']}
      queryFn={() => client.get('/admin/api/customs-points')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'code', label: '编码' },
        { key: 'country', label: '国家' },
        { key: 'port_type', label: '口岸类型' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
