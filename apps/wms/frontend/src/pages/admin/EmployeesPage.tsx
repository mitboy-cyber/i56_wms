import MinimalListPage from '@/components/MinimalListPage';
import client from '@/api/client';

export default function EmployeesPage() {
  return (
    <MinimalListPage title="员工管理" queryKey={['admin-employees']}
      queryFn={() => client.get('/admin/api/employees')}
      apiBase="/admin/api/employees"
      columns={[
        { key: 'id', label: '编号' }, { key: 'warehouse', label: '仓库' },
        { key: 'name', label: '姓名' }, { key: 'username', label: '账号' },
        { key: 'phone', label: '电话' }, { key: 'role', label: '角色' },
        { key: 'created_at', label: '创建时间' },
      ]} getRowId={(r: any, i: number) => String(r.id || i)} />
  );
}
