import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function RouteProfitReport() {
  return (
    <GenericListPage
      title="路线利润报表"
      queryKey={['admin-route-profit']}
      queryFn={() => client.get('/admin/api/report/route-profit')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'route', label: '路线' },
        { key: 'revenue', label: '收入' },
        { key: 'cost', label: '成本' },
        { key: 'profit', label: '利润' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
