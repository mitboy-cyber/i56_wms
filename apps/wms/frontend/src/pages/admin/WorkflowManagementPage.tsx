import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function WorkflowManagementPage() {
  return (
    <GenericListPage
      title="工单流程管理"
      queryKey={['admin-WorkflowManagement']}
      queryFn={() => client.get('/admin/api/workflow-management').then(r => r.data)}
      columns={[
        { key: 'id', label: '编号' },
        { key: 'warehouse', label: '仓库' },
        { key: 'process_id', label: '流程ID' },
        { key: 'name', label: '流程名称' },
        { key: 'steps', label: '流程工单' },
        { key: 'trigger_event', label: '触发' },
        { key: 'is_enabled', label: '启用', render: (v: unknown) => <>{v ? '✅ 启用' : '❌ 禁用'}</> },
        { key: 'updated_at', label: '最近编辑' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
