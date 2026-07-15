import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function WorkOrdersPage() {
  return (
    <GenericListPage title="工单管理" queryKey={['admin-work-orders']}
      queryFn={() => client.get('/admin/api/work-orders')}
      apiBase="/admin/api/work-orders"
      columns={[
        { key: 'id', label: '编号' }, { key: 'title', label: '标题' },
        { key: 'workflow_process', label: '工作流' }, { key: 'priority', label: '优先级' },
        { key: 'status', label: '状态' }, { key: 'assigned_worker', label: '执行人' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
