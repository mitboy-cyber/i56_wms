import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function SystemParamsPage() {
  return (
    <GenericListPage
      title="系统参数"
      queryKey={['admin-system-params']}
      queryFn={() => client.get('/admin/api/system/params')}
      apiBase="/admin/api/system/params"
      columns={[
        { key: 'id', label: '编号' },
        { key: 'key', label: '键名' },
        { key: 'value', label: '值' },
        { key: 'group', label: '分组' },
        { key: 'label', label: '标签' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
