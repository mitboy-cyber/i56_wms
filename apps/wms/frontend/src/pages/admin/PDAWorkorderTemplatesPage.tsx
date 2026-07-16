import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PDAWorkorderTemplatesPage() {
  return (
    <GenericListPage
      title="PDA 工单模板"
      queryKey={['admin-PDAWorkorderTemplates']}
      queryFn={() => client.get('/admin/api/pda-workorder-templates')}
      apiBase="/admin/api/pda-workorder-templates"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'warehouse', label: '仓库' },
        { key: 'template_id', label: '模板ID' },
        { key: 'name', label: '模板名称' },
        { key: 'work_type', label: '工种' },
        { key: 'workflow_id', label: '参与流程' },
        { key: 'default_priority', label: '默认优先级' },
        { key: 'is_enabled', label: '启用', render: (v: unknown) => <>{v ? '✅ 启用' : '❌ 禁用'}</> },
        { key: 'updated_at', label: '最近编辑' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
