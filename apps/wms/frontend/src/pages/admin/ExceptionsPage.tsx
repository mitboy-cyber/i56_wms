import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function ExceptionsPage() {
  return (
    <MinimalListPage title="异常记录" queryKey={['admin-ExceptionsPage']}
      queryFn={() => client.get('/admin/api/exceptions')}
      apiBase="/admin/api/exceptions"
      columns={[
        { key: 'id', label: '编号' }, { key: 'parcel_tracking', label: '关联包裹' },
        { key: 'type', label: '异常类型' }, { key: 'description', label: '描述' },
        { key: 'severity', label: '严重程度' }, { key: 'status', label: '状态' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
