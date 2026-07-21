import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function OrdersPage() {
  return (
    <MinimalListPage
      title="集运订单"
      queryKey={['admin-orders']}
      queryFn={() => client.get('/admin/api/orders')}
      apiBase="/admin/api/orders"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'order_no', label: '订单号' },
        { key: 'recipient_name', label: '收件人' },
        { key: 'parcel_count', label: '包裹数' },
        { key: 'total_price', label: '金额' },
        { key: 'total_actual_weight', label: '实重(kg)' },
        { key: 'total_chargeable_weight', label: '计费量(kg)' },
        { key: 'status', label: '状态' },
        { key: 'remark', label: '备注' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
