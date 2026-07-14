import client from '@/api/client';
import GenericListPage from '@/components/GenericListPage';

export default function SystemParamsPage() {
  return (
    <GenericListPage
      title="系统参数"
      queryKey={['admin-system-params']}
      queryFn={() => client.get('/admin/api/system/params')}
      columns={[
        { key: 'id', label: 'Id' },
        { key: 'key', label: 'Key' },
        { key: 'value', label: 'Value' },
        { key: 'group', label: 'Group' },
        { key: 'label', label: 'Label' },
      ]}
      getRowId={(r: any, i: number) => String(r.id || i)}
    />
  );
}
