import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function BrandSettingsPage() {
  return (
    <GenericListPage
      title="品牌设置"
      queryKey={['admin-brand-settings']}
      queryFn={() => client.get('/admin/api/brand-settings')}
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
