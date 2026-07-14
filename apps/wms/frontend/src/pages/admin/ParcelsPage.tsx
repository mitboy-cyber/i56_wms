import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ParcelsPage() {
  return (
    <GenericListPage
      title="包裹管理"
      queryKey={['admin-parcels']}
      queryFn={() => client.get('/admin/api/parcels')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'tracking_number', label: 'Tracking Number' },
        { key: 'product_name', label: 'Product Name' },
        { key: 'status', label: 'Status' },
        { key: 'warehouse_id', label: 'Warehouse Id' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
