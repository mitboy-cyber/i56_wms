import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiStoragePage() {
  return (
    <GenericListPage
      title="仓储API"
      queryKey={['admin-api-storage']}
      queryFn={() => client.get('/admin/api/api-storage')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'provider', label: '服务商' },
        { key: 'endpoint', label: '接口地址' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
