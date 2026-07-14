import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function WarehouseConsolePage() {
  return (
    <GenericListPage
      title="仓库控制台"
      queryKey={['admin-warehouse-console']}
      queryFn={() => client.get('/admin/api/warehouse-console')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'warehouse', label: '仓库' },
        { key: 'status', label: '状态' },
        { key: 'devices', label: '设备数' },
        { key: 'last_heartbeat', label: '最后心跳' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
