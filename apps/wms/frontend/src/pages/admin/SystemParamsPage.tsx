import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function SystemParamsPage() {
  return (
    <GenericListPage
      title="系统参数"
      queryKey={['admin-system-params']}
      queryFn={() => client.get('/admin/api/system-params')}
      columns={[
        { key: 'id', label: 'ID' },
        { key: 'key', label: '参数键' },
        { key: 'value', label: '参数值' },
        { key: 'description', label: '描述' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
