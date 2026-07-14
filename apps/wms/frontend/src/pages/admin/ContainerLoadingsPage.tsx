import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ContainerLoadingsPage() {
  return (
    <GenericListPage
      title="集装箱装货"
      queryKey={['admin-container-loadings']}
      queryFn={() => client.get('/admin/api/container-loadings')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'container_no', label: 'Container No' },
        { key: 'vessel', label: 'Vessel' },
        { key: 'port_from', label: 'Port From' },
        { key: 'port_to', label: 'Port To' },
        { key: 'parcel_count', label: 'Parcel Count' },
        { key: 'loaded_at', label: 'Loaded At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
