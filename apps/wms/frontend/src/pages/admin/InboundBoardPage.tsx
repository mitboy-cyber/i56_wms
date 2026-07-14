import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function InboundBoardPage() {
  return (
    <GenericListPage
      title="入库看板"
      queryKey={['admin-InboundBoardPage']}
      queryFn={() => client.get('/admin/api/inbound-board')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'parcel_no', label: 'Parcel No' },
        { key: 'warehouse', label: 'Warehouse' },
        { key: 'status', label: 'Status' },
        { key: 'expected_at', label: 'Expected At' },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
