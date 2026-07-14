import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ParcelsPage() {
  return (
    <GenericListPage
      title="包裹管理"
      queryKey={['admin-parcels']}
      queryFn={() => client.get('/admin/api/parcels')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'tracking_no', label: '运单号' },
        { key: 'order_no', label: '订单号' },
        { key: 'client_name', label: '客户' },
        { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
