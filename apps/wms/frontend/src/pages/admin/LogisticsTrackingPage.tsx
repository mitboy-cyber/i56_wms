import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

const STATUS_MAP: Record<string, string> = {
  '运输中': '📦 运输中', '已签收': '✅ 已签收', '派送中': '🚚 派送中',
  '已出库': '📤 已出库', '在途': '✈️ 在途',
};

export default function LogisticsTrackingPage() {
  return (
    <GenericListPage title="物流追踪" queryKey={['admin-logistics-tracking']}
      queryFn={() => client.get('/admin/api/logistics-tracking')}
      apiBase="/admin/api/logistics-tracking"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'tracking_no', label: '快递单号', render: (v: unknown) => <span className="font-mono text-blue-600">{String(v)}</span> },
        { key: 'route', label: '路线' },
        { key: 'region', label: '区域' },
        { key: 'company_name', label: '物流公司' },
        { key: 'status', label: '状态', render: (v: unknown) => <span>{STATUS_MAP[String(v)] || String(v)}</span> },
        { key: 'detail', label: '最新轨迹', render: (v: unknown) => <span className="text-xs text-gray-600 max-w-xs truncate block">{String(v)}</span> },
        { key: 'order_no', label: '关联订单', render: (v: unknown) => <span className="text-xs">{v && String(v) !== '' ? String(v) : '—'}</span> },
        { key: 'error_count', label: '失败次数', render: (v: unknown) => Number(v) > 0 ? <span className="text-red-500 font-medium">{String(v)}</span> : <span className="text-gray-400">0</span> },
        { key: 'updated_at', label: '更新时间', render: (v: unknown) => <span className="text-xs">{v ? new Date(String(v)).toLocaleString('zh-CN') : '—'}</span> },
      ]} getRowId={(_, i) => String(i)} />
  );
}
