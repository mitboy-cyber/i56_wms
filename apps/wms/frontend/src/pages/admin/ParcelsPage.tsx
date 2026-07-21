import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function ParcelsPage() {
  return (
    <MinimalListPage
      title="包裹列表"
      queryKey={['admin-parcels']}
      queryFn={() => client.get('/admin/api/parcels')}
      apiBase="/admin/api/parcels"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'tracking_number', label: '快递单号' },
        { key: 'product_name', label: '品名' },
        { key: 'parcel_name', label: '包裹名' },
        { key: 'status', label: '状态' },
        { key: 'cargo_type', label: '货物类型' },
        { key: 'actual_weight', label: '实重(kg)' },
        { key: 'courier_code', label: '快递公司' },
        { key: 'created_at', label: '入库时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
