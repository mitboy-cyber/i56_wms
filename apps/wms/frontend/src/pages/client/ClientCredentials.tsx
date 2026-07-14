import clientApi from '@/api/clientApi';
import GenericListPage from '@/components/GenericListPage';
export default function ClientCredentials() {
  return <GenericListPage title="API凭证" queryKey={['client-credentials']} queryFn={() => clientApi.credentials()}
    columns={[{key:'app_key',label:'App Key'},{key:'masked_secret',label:'Secret'},{key:'active',label:'状态',render:(v:unknown)=>v?'✅ 启用':'❌ 停用'}]}
    getRowId={(r:Record<string,unknown>)=>String(r.app_key)} />;
}
