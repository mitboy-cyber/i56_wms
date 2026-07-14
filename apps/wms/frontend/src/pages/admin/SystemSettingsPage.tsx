import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function SystemSettingsPage() {
  return (
    <GenericListPage
      title="系统设置"
      queryKey={['admin-SystemSettingsPage']}
      queryFn={() => client.get('/admin/api/system/settings')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'key', label: 'Key' },
        { key: 'value', label: 'Value' },
        { key: 'group', label: 'Group' },
        { key: 'label', label: 'Label' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
