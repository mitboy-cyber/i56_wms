import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function AuditLogsPage() {
  return (
    <GenericListPage
      title="审计日志"
      queryKey={['admin-audit-logs']}
      queryFn={() => client.get('/admin/api/system/audit-logs')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'user', label: '操作人' },
        { key: 'action', label: '操作' },
        { key: 'resource', label: '资源' },
        { key: 'ip', label: 'IP' },
        { key: 'created_at', label: '时间' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
