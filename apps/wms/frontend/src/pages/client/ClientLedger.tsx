import { useQuery } from '@tanstack/react-query';
import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';

export default function ClientLedger() {
  return (
    <GenericListPage title="余额明细" queryKey={['client-ledger']}
      queryFn={() => clientApi.ledger()}
      apiBase="/client/api/ledger"
      columns={[
        { key: 'type', label: '类型' },
        { key: 'amount', label: '金额', render: (v:unknown)=>`¥${Number(v).toFixed(2)}` },
        { key: 'balance_after', label: '余额', render: (v:unknown)=>`¥${Number(v).toFixed(2)}` },
        { key: 'description', label: '描述' },
        { key: 'created_at', label: '时间' },
      ]}
      getRowId={(r:any)=>String(r.id)} />
  );
}
