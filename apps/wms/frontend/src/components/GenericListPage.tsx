import { useQuery } from '@tanstack/react-query';

interface ColumnDef {
  key: string;
  label: string;
  render?: (value: any, row: any) => React.ReactNode;
}

interface GenericListPageProps {
  title: string;
  queryKey: string[];
  queryFn: () => Promise<any>;
  columns: ColumnDef[];
  getRowId: (row: any, i: number) => string;
}

export default function GenericListPage({
  title, queryKey, queryFn, columns, getRowId,
}: GenericListPageProps) {
  const result = useQuery({ queryKey, queryFn } as any);
  const rows: any[] = (result as any).data?.data ?? [];
  const isLoading = (result as any).isLoading;

  return (
    <div>
      <h2 className="text-xl font-bold text-gray-800 mb-4">{title}</h2>
      {isLoading ? (
        <div className="text-center py-8 text-gray-400">加载中...</div>
      ) : rows.length === 0 ? (
        <div className="text-center py-8 text-gray-400">暂无数据</div>
      ) : (
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-gray-50 border-b border-gray-200">
                {columns.map((c) => (
                  <th key={c.key} className="text-left px-4 py-3 font-medium text-gray-600">{c.label}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((row: any, i: number) => (
                <tr key={getRowId(row, i)} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                  {columns.map((c) => (
                    <td key={c.key} className="px-4 py-2.5 text-gray-700">
                      {c.render ? c.render(row[c.key], row) : String(row[c.key] ?? '')}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
