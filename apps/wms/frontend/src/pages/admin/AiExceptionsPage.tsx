import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function AiExceptionsPage() {
  return (
    <GenericListPage
      title="AI异常"
      queryKey={['admin-AiExceptionsPage']}
      queryFn={() => client.get('/admin/api/ai-exceptions')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'parcel_id', label: 'Parcel Id' },
        { key: 'reason', label: 'Reason' },
        { key: 'confidence', label: 'Confidence' },
        { key: 'reviewed', label: 'Reviewed' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
