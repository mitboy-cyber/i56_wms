import client from '@/api/client';
import MinimalListPage from '@/components/MinimalListPage';

export default function RouteTemplatesPage() {
  return (
    <MinimalListPage title="路线模板" queryKey={['admin-route-templates']}
      queryFn={() => client.get('/admin/api/route-templates')}
      apiBase="/admin/api/route-templates"
      columns={[
        { key: 'id', label: '编号' }, { key: 'name', label: '名称' },
        { key: 'from', label: '起点' }, { key: 'to', label: '终点' },
        { key: 'carrier_id', label: '承运商' }, { key: 'est_days', label: '预计天数' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
