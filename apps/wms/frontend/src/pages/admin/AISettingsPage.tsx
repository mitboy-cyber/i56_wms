import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function AISettingsPage() {
  return (
    <GenericListPage
      title="AI设置"
      queryKey={['admin-ai-settings']}
      queryFn={() => client.get('/admin/api/system/ai-settings')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'key', label: 'Key' },
        { key: 'value', label: 'Value' },
        { key: 'group', label: 'Group' },
        { key: 'label', label: 'Label' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
