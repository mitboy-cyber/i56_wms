import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientProfitReport() {
  return (
    <GenericListPage
      title="客户利润报表"
      queryKey={['admin-client-profit']}
      queryFn={() => client.get('/admin/api/report/client-profit')}
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
