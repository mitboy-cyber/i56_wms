import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function StorageConfigPage() {
  return (
    <GenericListPage title="存储配置" queryKey={['admin-storage']}
      queryFn={() => client.get('/admin/api/storage')}
      apiBase="/admin/api/storage"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'provider', label: '提供商' }, { key: 'bucket', label: '存储桶' },
        { key: 'region', label: '区域' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
