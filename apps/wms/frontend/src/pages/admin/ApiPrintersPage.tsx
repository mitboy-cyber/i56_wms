import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiPrintersPage() {
  return (
    <GenericListPage
      title="打印机API"
      queryKey={['admin-api-printers']}
      queryFn={() => client.get('/admin/api/system/api-printers')}
      apiBase="/admin/api/system/api-printers"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'name', label: '名称' },
        { key: 'provider', label: '提供商' },
        { key: 'endpoint', label: '接口地址' },
        { key: 'api_key', label: 'API密钥' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
