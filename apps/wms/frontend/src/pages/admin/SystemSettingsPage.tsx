import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function SystemSettingsPage() {
  return (
    <GenericListPage
      title="系统设置"
      queryKey={['admin-SystemSettingsPage']}
      queryFn={() => client.get('/admin/api/system/settings')}
      apiBase="/admin/api/system/settings"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'key', label: '键名' },
        { key: 'value', label: '值' },
        { key: 'group', label: '分组' },
        { key: 'label', label: '标签' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
