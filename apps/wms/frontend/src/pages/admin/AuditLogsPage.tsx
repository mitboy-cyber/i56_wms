import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function AuditLogsPage() {
  return (
    <GenericListPage
      title="审计日志"
      queryKey={['admin-audit-logs']}
      queryFn={() => client.get('/admin/api/system/audit-logs')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'user_id', label: 'User Id' },
        { key: 'action', label: 'Action' },
        { key: 'resource', label: 'Resource' },
        { key: 'detail', label: 'Detail' },
        { key: 'created_at', label: 'Created At' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
