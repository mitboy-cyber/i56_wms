import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function LogisticsTrackingPage() {
  return (
    <GenericListPage
      title="物流追踪"
      queryKey={['admin-logistics-tracking']}
      queryFn={() => client.get('/admin/api/logistics-tracking')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'tracking_no', label: '运单号' },
        { key: 'order_no', label: '订单号' },
        { key: 'status', label: '状态' },
        { key: 'location', label: '当前位置' },
        { key: 'updated_at', label: '更新时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
