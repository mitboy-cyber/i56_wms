import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function BrandSettingsPage() {
  return (
    <GenericListPage
      title="品牌设置"
      queryKey={['admin-brand-settings']}
      queryFn={() => client.get('/admin/api/system/brand')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'key', label: 'Key' },
        { key: 'value', label: 'Value' },
        { key: 'group', label: 'Group' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
