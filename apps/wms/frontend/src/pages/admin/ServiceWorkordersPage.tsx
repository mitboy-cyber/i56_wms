import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function ServiceWorkordersPage() {
  return (
    <MinimalListPage title="服务工单" queryKey={['admin-service-workorders']}
      queryFn={() => client.get('/admin/api/service-workorders')}
      apiBase="/admin/api/service-workorders"
      columns={[
        { key: 'id', label: '编号' }, { key: 'title', label: '工单标题' },
        { key: 'service_type', label: '服务类型' }, { key: 'assigned_to', label: '执行人' },
        { key: 'priority', label: '优先级' }, { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
