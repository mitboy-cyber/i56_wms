import { useQuery } from '@tanstack/react-query';
import client from '@/api/client';
import { Server, Database, HardDrive, Activity, Wifi, Clock } from 'lucide-react';

export default function TaskMonitorPage() {
  const { data: stats } = useQuery<any>({ queryKey: ['monitor-stats'], queryFn: () => client.get('/admin/api/dashboard/stats').then(r => r.data) });

  const serverTime = new Date().toLocaleString('zh-CN');

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-gray-800">系统监控</h1>
        <div className="flex items-center gap-2">
          <span className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
          <span className="text-sm text-green-700">系统正常</span>
        </div>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-6 gap-4">
        {[
          { icon: Server, label: '订单', value: stats?.total_orders || 0, color: 'bg-blue-50 border-blue-200' },
          { icon: HardDrive, label: '包裹', value: stats?.total_parcels || 0, color: 'bg-indigo-50 border-indigo-200' },
          { icon: Database, label: '客户', value: stats?.total_clients || 0, color: 'bg-green-50 border-green-200' },
          { icon: Wifi, label: '快递', value: stats?.total_couriers || 0, color: 'bg-teal-50 border-teal-200' },
          { icon: Clock, label: '服务模板', value: stats?.active_templates || 0, color: 'bg-purple-50 border-purple-200' },
          { icon: Activity, label: 'API状态', value: '正常', color: 'bg-emerald-50 border-emerald-200' },
        ].map(c => (
          <div key={c.label} className={`rounded-lg border p-3 ${c.color}`}>
            <div className="flex items-center gap-1.5 mb-1"><c.icon size={16} className="opacity-60" /><span className="text-xs font-medium opacity-70">{c.label}</span></div>
            <div className="text-xl font-bold">{c.value}</div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="bg-white rounded-lg shadow p-4">
          <h3 className="text-sm font-semibold text-gray-600 mb-3">服务状态</h3>
          <div className="space-y-3">
            {[
              { name: 'API Gateway', status: '正常', color: 'bg-green-500' },
              { name: '订单服务', status: '正常', color: 'bg-green-500' },
              { name: '仓储服务', status: '正常', color: 'bg-green-500' },
              { name: '物流服务', status: '正常', color: 'bg-green-500' },
              { name: '认证服务', status: '正常', color: 'bg-green-500' },
              { name: '通知服务', status: '就绪', color: 'bg-yellow-500' },
            ].map(s => (
              <div key={s.name} className="flex items-center justify-between">
                <span className="text-sm">{s.name}</span>
                <div className="flex items-center gap-2">
                  <span className={`h-2 w-2 rounded-full ${s.color}`} />
                  <span className="text-sm text-gray-500">{s.status}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-4">
          <h3 className="text-sm font-semibold text-gray-600 mb-3">系统信息</h3>
          <div className="space-y-3 text-sm">
            {[
              ['服务时间', serverTime],
              ['框架版本', 'I56 Framework v2.4.2'],
              ['存储模式', 'Memory (开发模式)'],
              ['租户模式', 'Single Tenant'],
              ['前端构建', 'v95 (74 chunks)'],
              ['运行状态', 'Production'],
            ].map(([k, v]) => (
              <div key={k} className="flex justify-between">
                <span className="text-gray-500">{k}</span>
                <span className="font-medium text-right">{v}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
