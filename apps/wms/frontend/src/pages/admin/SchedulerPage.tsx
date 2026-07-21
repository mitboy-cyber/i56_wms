import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function SchedulerPage() {
  return (
    <MinimalListPage title="定时任务" queryKey={['admin-scheduler']}
      queryFn={() => client.get('/admin/api/system/scheduler')}
      apiBase="/admin/api/system/scheduler"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '任务名称' },
        { key: 'cron', label: 'Cron表达式' }, { key: 'enabled', label: '启用' },
        { key: 'last_run', label: '上次执行' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
