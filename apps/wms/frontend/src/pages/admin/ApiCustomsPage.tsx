import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiCustomsPage() {
  return (
    <GenericListPage
      title="报关API"
      queryKey={['admin-api-customs']}
      queryFn={() => client.get('/admin/api/api-customs')}
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
