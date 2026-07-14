import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function TaskMonitorPage() {
  return (
    <GenericListPage
      title="任务监控"
      queryKey={['admin-task-monitor']}
      queryFn={() => client.get('/admin/api/task-monitor')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'task_name', label: '任务名称' },
        { key: 'type', label: '类型' },
        { key: 'status', label: '状态' },
        { key: 'progress', label: '进度' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
