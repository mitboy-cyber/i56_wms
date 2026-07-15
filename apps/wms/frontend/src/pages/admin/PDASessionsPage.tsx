import GenericListPage from '@/components/GenericListPage';
import client from '@/api/client';

export default function PDASessionsPage() {
  return (
    <GenericListPage
      title="PDA 在线会话"
      queryKey={['admin-PDASessions']}
      queryFn={() => client.get('/admin/api/pda-sessions').then(r => r.data)}
      columns={[
        { key: 'id', label: '编号' },
        { key: 'warehouse', label: '仓库' },
        { key: 'worker_name', label: '工人' },
        { key: 'device', label: '设备' },
        { key: 'login_at', label: '登录时间' },
        { key: 'last_heartbeat', label: '最近心跳' },
        { key: 'online_duration', label: '在线时长' },
        { key: 'current_page', label: '当前页面' },
        { key: 'current_area', label: '当前区域' },
        { key: 'current_location', label: '当前货位' },
        { key: 'is_online', label: '在线', render: (v: unknown) => <>{v ? '✅ 在线' : '❌ 离线'}</> },
        { key: 'logout_at', label: '登出时间', render: (v: unknown) => <>{v || '—'}</> },
      ]}
      getRowId={(_, i) => String(i)}
    />
  );
}
