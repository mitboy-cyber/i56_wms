import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function WarehouseBoardPage() {
  return (
    <GenericListPage
      title="仓库看板"
      queryKey={['admin-warehouse-board']}
      queryFn={() => client.get('/admin/api/warehouse-board')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'warehouse', label: '仓库' },
        { key: 'inbound', label: '入库' },
        { key: 'outbound', label: '出库' },
        { key: 'stock', label: '库存' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
