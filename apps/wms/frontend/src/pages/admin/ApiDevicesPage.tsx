import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiDevicesPage() {
  return (
    <GenericListPage
      title="设备API"
      queryKey={['admin-api-devices']}
      queryFn={() => client.get('/admin/api/api-devices')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'type', label: '类型' },
        { key: 'serial_no', label: '序列号' },
        { key: 'status', label: '状态' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
