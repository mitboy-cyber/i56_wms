import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

const PAGE_CN: Record<string, string> = {
  receive: '收货', weighing: '称重', shelve: '上架',
  pick: '拣货', pack: '打包', load: '装车', dispatch: '出库'
};

export default function PDASessionsPage() {
  return (
    <GenericListPage title="PDA 在线会话" queryKey={['admin-pda-sessions']}
      queryFn={() => client.get('/admin/api/pda-sessions')}
      apiBase="/admin/api/pda-sessions"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'warehouse', label: '仓库' },
        { key: 'worker_name', label: '操作员' },
        { key: 'device', label: '设备编号' },
        { key: 'login_at', label: '登录时间', render: (v: unknown) => <span className="text-xs">{v ? new Date(String(v)).toLocaleString('zh-CN') : '—'}</span> },
        { key: 'last_heartbeat', label: '最近心跳', render: (v: unknown) => <span className="text-xs">{v ? new Date(String(v)).toLocaleString('zh-CN') : '—'}</span> },
        { key: 'online_duration', label: '在线时长' },
        { key: 'current_page', label: '当前页面', render: (v: unknown) => <span className="text-blue-600 font-medium">{PAGE_CN[String(v)] || String(v)}</span> },
        { key: 'current_area', label: '当前区域' },
        { key: 'current_location', label: '当前货位' },
        { key: 'is_online', label: '在线', render: (v: unknown) => <span className={v ? 'text-green-600' : 'text-red-500'}>{v ? '🟢 在线' : '⚫ 离线'}</span> },
        { key: 'logout_at', label: '登出时间', render: (v: unknown) => <span className="text-xs text-gray-400">{v ? new Date(String(v)).toLocaleString('zh-CN') : '—'}</span> },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
