import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function SchedulerPage() {
  return (
    <GenericListPage
      title="定时任务"
      queryKey={['admin-scheduler']}
      queryFn={() => client.get('/admin/api/system/scheduler')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'job_name', label: '任务名称' },
        { key: 'cron_expr', label: 'Cron表达式' },
        { key: 'status', label: '状态' },
        { key: 'last_run', label: '上次执行' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
