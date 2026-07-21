import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function AuditLogsPage() {
  return (
    <MinimalListPage title="审计日志" queryKey={['admin-audit-logs']}
      queryFn={() => client.get('/admin/api/system/audit-logs')}
      apiBase="/admin/api/system/audit-logs"
      columns={[
        { key: 'id', label: '编号' }, { key: 'user_id', label: '用户' },
        { key: 'action', label: '操作' }, { key: 'resource', label: '资源' },
        { key: 'detail', label: '详情' }, { key: 'created_at', label: '时间' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
