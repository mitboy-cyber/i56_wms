import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function WorkflowManagementPage() {
  return (
    <MinimalListPage title="工作流管理" queryKey={['admin-workflow']}
      queryFn={() => client.get('/admin/api/workflow')}
      apiBase="/admin/api/workflow"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '流程名称' },
        { key: 'steps', label: '步骤数' }, { key: 'module', label: '所属模块' },
        { key: 'status', label: '状态' }, { key: 'updated_at', label: '更新时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
