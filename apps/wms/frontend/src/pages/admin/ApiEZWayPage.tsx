import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiEZWayPage() {
  return (
    <GenericListPage
      title="EZWay API"
      queryKey={['admin-api-ezway']}
      queryFn={() => client.get('/admin/api/api-ezway')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'endpoint', label: '接口地址' },
        { key: 'api_key', label: 'API Key' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
