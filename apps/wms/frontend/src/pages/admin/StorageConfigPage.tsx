import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function StorageConfigPage() {
  return (
    <GenericListPage
      title="仓储配置"
      queryKey={['admin-storage-config']}
      queryFn={() => client.get('/admin/api/storage')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'warehouse_name', label: '仓库' },
        { key: 'config_key', label: '配置项' },
        { key: 'config_value', label: '配置值' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
