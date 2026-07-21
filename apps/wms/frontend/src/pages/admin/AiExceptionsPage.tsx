import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function AiExceptionsPage() {
  return (
    <MinimalListPage title="AI 异常检测" queryKey={['admin-AiExceptionsPage']}
      queryFn={() => client.get('/admin/api/ai-exceptions')}
      apiBase="/admin/api/ai-exceptions"
      columns={[
        { key: 'id', label: '编号' }, { key: 'parcel_id', label: '包裹编号' },
        { key: 'reason', label: '异常原因' }, { key: 'confidence', label: '置信度' },
        { key: 'reviewed', label: '已审核' }, { key: 'created_at', label: '时间' },
      ]} getRowId={(_, i) => String(i)} />
  );
}
