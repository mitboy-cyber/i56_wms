import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function WarehouseConsolePage() {
  return (
    <GenericListPage
      title="仓库控制台"
      queryKey={['admin-warehouse-console']}
      queryFn={() => client.get('/admin/api/warehouse-console')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'warehouse_id', label: 'Warehouse Id' },
        { key: 'name', label: 'Name' },
        { key: 'machine', label: 'Machine' },
        { key: 'status', label: 'Status' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
