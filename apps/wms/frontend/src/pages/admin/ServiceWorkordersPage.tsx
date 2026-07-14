import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ServiceWorkordersPage() {
  return (
    <GenericListPage
      title="服务工单"
      queryKey={['admin-ServiceWorkordersPage']}
      queryFn={() => client.get('/admin/api/service-workorders')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'order_no', label: 'Order No' },
        { key: 'type', label: 'Type' },
        { key: 'status', label: 'Status' },
        { key: 'assigned_to', label: 'Assigned To' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
