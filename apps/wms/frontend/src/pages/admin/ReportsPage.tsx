import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ReportsPage() {
  return (
    <GenericListPage
      title="报表管理"
      queryKey={['admin-reports']}
      queryFn={() => client.get('/admin/api/system/reports')}
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
