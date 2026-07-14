import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
export default function ClientOrderDetail() {
  const { id } = useParams();
  const { data } = useQuery({ queryKey: ['client-order', id], queryFn: () => clientApi.orderDetail(id!), enabled: !!id } as any);
  const o: any = (data as any)?.data;
  if (!o) return <div className="text-center py-8 text-gray-400">加载中...</div>;
  return (
    <div>
      <h2 className="text-xl font-bold text-gray-800 mb-4">订单详情</h2>
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 space-y-2">
        {Object.entries(o).map(([k,v]) => (
          <div key={k} className="flex gap-4"><span className="text-gray-500 w-24 text-sm">{k}:</span><span className="text-gray-800">{String(v)}</span></div>
        ))}
      </div>
    </div>
  );
}
