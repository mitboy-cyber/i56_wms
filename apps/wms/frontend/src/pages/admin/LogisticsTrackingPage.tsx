import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function LogisticsTrackingPage() {
  return (
    <GenericListPage
      title="物流追踪"
      queryKey={['admin-logistics-tracking']}
      queryFn={() => client.get('/admin/api/logistics-tracking')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'tracking_no', label: 'Tracking No' },
        { key: 'location', label: 'Location' },
        { key: 'status', label: 'Status' },
        { key: 'updated_at', label: 'Updated At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
