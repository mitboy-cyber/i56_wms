import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientWebhooks() {
  return <GenericListPage title="Webhook" queryKey={['client-webhooks']} queryFn={() => clientApi.webhooks()}
    columns={[{key:'id',label:'#'}]}
    getRowId={(_:unknown,i:number)=>String(i)} />;
}
