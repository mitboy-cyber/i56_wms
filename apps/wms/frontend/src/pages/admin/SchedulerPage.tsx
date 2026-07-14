import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function SchedulerPage() {
  return (
    <GenericListPage
      title="定时任务"
      queryKey={['admin-scheduler']}
      queryFn={() => client.get('/admin/api/system/scheduler')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'cron', label: 'Cron' },
        { key: 'enabled', label: 'Enabled' },
        { key: 'last_run', label: 'Last Run' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
