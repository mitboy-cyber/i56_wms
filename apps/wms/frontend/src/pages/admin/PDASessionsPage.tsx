import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PDASessionsPage() {
  return (
    <GenericListPage title="PDA在线会话" queryKey={['admin-pda-sessions']}
      queryFn={() => client.get('/admin/api/pda-sessions')}
      apiBase="/admin/api/pda-sessions"
      columns={[
        { key: 'id', label: '编号' }, { key: 'device_id', label: '设备编号' },
        { key: 'worker_name', label: '操作员' }, { key: 'warehouse', label: '仓库' },
        { key: 'task_count', label: '任务数' }, { key: 'status', label: '在线状态' },
        { key: 'last_active', label: '最后活跃' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
