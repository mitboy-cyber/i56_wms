import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function ReportsPage() {
  return (
    <GenericListPage title="报表管理" queryKey={['admin-reports']}
      queryFn={() => client.get('/admin/api/system/reports')}
      apiBase="/admin/api/system/reports"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '报表名称' },
        { key: 'type', label: '类型' }, { key: 'status', label: '状态' },
        { key: 'created_at', label: '生成时间' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
