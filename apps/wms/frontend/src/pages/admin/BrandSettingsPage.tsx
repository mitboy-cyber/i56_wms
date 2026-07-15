import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function BrandSettingsPage() {
  return (
    <GenericListPage
      title="品牌设置"
      queryKey={['admin-brand-settings']}
      queryFn={() => client.get('/admin/api/system/brand')}
      apiBase="/admin/api/system/brand"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'key', label: '键名' },
        { key: 'value', label: '值' },
        { key: 'group', label: '分组' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
