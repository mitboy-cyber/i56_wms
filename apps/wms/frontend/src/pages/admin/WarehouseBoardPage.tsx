import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function WarehouseBoardPage() {
  return (
    <GenericListPage
      title="仓库看板"
      queryKey={['admin-warehouse-board']}
      queryFn={() => client.get('/admin/api/warehouse-board')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'pending_receive', label: 'Pending Receive' },
        { key: 'in_stock', label: 'In Stock' },
        { key: 'picking', label: 'Picking' },
        { key: 'outbound', label: 'Outbound' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
