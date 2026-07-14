import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function CustomsBrokersPage() {
  return (
    <GenericListPage
      title="报关行管理"
      queryKey={['admin-customs-brokers']}
      queryFn={() => client.get('/admin/api/customs-brokers')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'name', label: 'Name' },
        { key: 'license', label: 'License' },
        { key: 'contact', label: 'Contact' },
        { key: 'phone', label: 'Phone' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
