import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ServiceProfitReport() {
  return (
    <GenericListPage
      title="服务利润报表"
      queryKey={['admin-service-profit']}
      queryFn={() => client.get('/admin/api/report/service-profit')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'service_type', label: '服务类型' },
        { key: 'revenue', label: '收入' },
        { key: 'cost', label: '成本' },
        { key: 'profit', label: '利润' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
