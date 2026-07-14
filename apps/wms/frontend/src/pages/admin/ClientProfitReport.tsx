import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ClientProfitReport() {
  return (
    <GenericListPage
      title="客户利润报表"
      queryKey={['admin-client-profit']}
      queryFn={() => client.get('/admin/api/report/client-profit')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'client_name', label: '客户' },
        { key: 'revenue', label: '收入' },
        { key: 'cost', label: '成本' },
        { key: 'profit', label: '利润' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
