import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function ExceptionsPage() {
  return (
    <GenericListPage
      title="异常列表"
      queryKey={['admin-ExceptionsPage']}
      queryFn={() => client.get('/admin/api/exceptions')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'parcel_id', label: 'Parcel Id' },
        { key: 'type', label: 'Type' },
        { key: 'description', label: 'Description' },
        { key: 'status', label: 'Status' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
