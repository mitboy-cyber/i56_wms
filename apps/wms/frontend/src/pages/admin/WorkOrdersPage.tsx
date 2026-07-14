import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function WorkOrdersPage() {
  return (
    <GenericListPage
      title="工单管理"
      queryKey={['admin-work-orders']}
      queryFn={() => client.get('/admin/api/work-orders')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'order_no', label: 'Order No' },
        { key: 'type', label: 'Type' },
        { key: 'status', label: 'Status' },
        { key: 'assigned_to', label: 'Assigned To' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
