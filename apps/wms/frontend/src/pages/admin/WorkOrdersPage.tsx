import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function WorkOrdersPage() {
  return (
    <GenericListPage
      title="工单管理"
      queryKey={['admin-work-orders']}
      queryFn={() => client.get('/admin/api/work-orders')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'wo_no', label: '工单号' },
        { key: 'type', label: '类型' },
        { key: 'assignee', label: '负责人' },
        { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
