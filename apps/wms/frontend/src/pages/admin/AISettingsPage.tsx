import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function AISettingsPage() {
  return (
    <GenericListPage
      title="AI设置"
      queryKey={['admin-ai-settings']}
      queryFn={() => client.get('/admin/api/ai-settings')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'key', label: '配置项' },
        { key: 'value', label: '配置值' },
        { key: 'description', label: '描述' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
