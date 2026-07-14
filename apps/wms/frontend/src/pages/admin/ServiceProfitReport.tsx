import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ServiceProfitReport() {
  return (
    <GenericListPage
      title="服务利润报表"
      queryKey={['admin-service-profit']}
      queryFn={() => client.get('/admin/api/report/service-profit')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'type', label: 'Type' },
        { key: 'status', label: 'Status' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
