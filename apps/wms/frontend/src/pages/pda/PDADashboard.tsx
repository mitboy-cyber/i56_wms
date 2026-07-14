import { useQuery } from '@tanstack/react-query';
import pdaApi from '@/api/pdaApi';
export default function PDADashboard() {
  const { data } = useQuery({ queryKey: ['pda-dashboard'], queryFn: () => pdaApi.dashboard() } as any);
  const d = (data as any)?.data;
  return (
    <div>
      <h2 className="text-xl font-bold mb-4">PDA 作业台</h2>
      <div className="bg-white rounded-xl shadow-sm border p-4">
        <p className="text-gray-500">操作员 #{d?.op_id}</p>
        <pre className="text-xs mt-2 text-gray-600 overflow-auto max-h-60">{JSON.stringify(d?.stats, null, 2)}</pre>
      </div>
    </div>
  );
}
