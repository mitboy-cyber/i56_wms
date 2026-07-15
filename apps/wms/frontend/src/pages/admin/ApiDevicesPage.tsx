import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ApiDevicesPage() {
  return (
    <GenericListPage title="设备管理" queryKey={['admin-api-devices']}
      queryFn={() => client.get('/admin/api/devices')}
      apiBase="/admin/api/devices"
      columns={[
        { key: 'id', label: '编号' }, { key: 'device_name', label: '设备名称' },
        { key: 'device_type', label: '设备类型' }, { key: 'device_code', label: '设备编号' },
        { key: 'ip_address', label: 'IP地址' }, { key: 'status', label: '状态' },
        { key: 'warehouse_id', label: '所属仓库' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
