import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function LogisticsTrackingPage() {
  return (
    <GenericListPage title="物流追踪" queryKey={['admin-logistics-tracking']}
      queryFn={() => client.get('/admin/api/logistics-tracking')}
      apiBase="/admin/api/logistics-tracking"
      columns={[
        { key: 'id', label: '编号' }, { key: 'tracking_no', label: '快递单号' },
        { key: 'courier_name', label: '快递公司' }, { key: 'status', label: '物流状态' },
        { key: 'location', label: '当前位置' }, { key: 'updated_at', label: '更新时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
