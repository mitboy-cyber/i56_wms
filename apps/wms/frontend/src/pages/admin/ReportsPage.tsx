import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ReportsPage() {
  return (
    <GenericListPage
      title="报表管理"
      queryKey={['admin-reports']}
      queryFn={() => client.get('/admin/api/reports')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'name', label: '名称' },
        { key: 'type', label: '类型' },
        { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
