import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

const STATUS_MAP: Record<string, string> = {
  pending: '待处理', in_progress: '进行中', completed: '已完成', cancelled: '已取消'
};
const PRIORITY_MAP: Record<number, string> = { 1: '普通', 2: '中', 3: '高', 4: '紧急', 5: '特急' };
const WH_MAP: Record<number, string> = { 1: '厦门仓', 2: '深圳仓', 3: '上海仓' };

export default function WorkOrdersPage() {
  return (
    <GenericListPage title="工单管理" queryKey={['admin-work-orders']}
      queryFn={() => client.get('/admin/api/work-orders')}
      apiBase="/admin/api/work-orders"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'title', label: '标题' },
        { key: 'warehouse_id', label: '仓库', render: (v: unknown) => <span>{WH_MAP[Number(v)] || String(v)}</span> },
        { key: 'status', label: '状态', render: (v: unknown) => <span className={`text-sm font-medium ${String(v)==='in_progress'?'text-blue-600':String(v)==='completed'?'text-green-600':'text-amber-600'}`}>{STATUS_MAP[String(v)] || String(v)}</span> },
        { key: 'assigned_to', label: '当前工人', render: (v: unknown) => <span>{v ? String(v) : '—'}</span> },
        { key: 'priority', label: '优先级', render: (v: unknown) => <span>{PRIORITY_MAP[Number(v)] || '普通'}</span> },
        { key: 'description', label: '业务单号/描述', render: (v: unknown) => <span className="text-xs text-gray-500 max-w-xs truncate block">{String(v)}</span> },
        { key: 'created_at', label: '派发时间', render: (v: unknown) => <span className="text-xs">{v ? new Date(String(v)).toLocaleString('zh-CN') : '—'}</span> },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
