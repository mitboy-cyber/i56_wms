import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useEffect, useState } from 'react';
import { useAuthStore } from '@/stores/auth';
import { useClientAuth } from '@/stores/clientAuth';
import { DashboardLayout } from '@/layouts/DashboardLayout';
import ClientLayout from '@/layouts/ClientLayout';

// Admin pages
import { LoginPage } from '@/pages/Login';
import { DashboardPage } from '@/pages/Dashboard';
import { EmployeesPage } from '@/pages/Employees';
import { WarehousesPage } from '@/pages/Warehouses';
import ClientsPage from '@/pages/admin/ClientsPage';
import RolesPage from '@/pages/admin/RolesPage';
import ParcelsPage from '@/pages/admin/ParcelsPage';
import CarriersPage from '@/pages/admin/CarriersPage';
import CouriersPage from '@/pages/admin/CouriersPage';
import DeclarantsPage from '@/pages/admin/DeclarantsPage';

// Client pages
import ClientLogin from '@/pages/client/ClientLogin';
import ClientDashboard from '@/pages/client/ClientDashboard';
import ClientParcels from '@/pages/client/ClientParcels';
import ClientOrders from '@/pages/client/ClientOrders';
import ClientLedger from '@/pages/client/ClientLedger';
import ClientDeclarants from '@/pages/client/ClientDeclarants';
import ClientMembers from '@/pages/client/ClientMembers';
import ClientAddresses from '@/pages/client/ClientAddresses';
import ClientWarehouses from '@/pages/client/ClientWarehouses';
import ClientCouriers from '@/pages/client/ClientCouriers';
import ClientServices from '@/pages/client/ClientServices';
import ClientCredentials from '@/pages/client/ClientCredentials';
import ClientRoutePrices from '@/pages/client/ClientRoutePrices';
import ClientDeliveryFees from '@/pages/client/ClientDeliveryFees';
import ClientSurcharges from '@/pages/client/ClientSurcharges';
import ClientWebhooks from '@/pages/client/ClientWebhooks';
import ClientOrderNew from '@/pages/client/ClientOrderNew';
import ClientOrderDetail from '@/pages/client/ClientOrderDetail';
import ClientWeightDashboard from '@/pages/client/ClientWeightDashboard';

// PDA pages
import PDALogin from '@/pages/pda/PDALogin';
import PDALayout from '@/pages/pda/PDALayout';
import PDADashboard from '@/pages/pda/PDADashboard';
import PDAReceive from '@/pages/pda/PDAReceive';
import PDAWeigh from '@/pages/pda/PDAWeigh';
import PDAPutaway from '@/pages/pda/PDAPutaway';
import PDAPick from '@/pages/pda/PDAPick';
import PDAPack from '@/pages/pda/PDAPack';
import PDALoad from '@/pages/pda/PDALoad';
import PDAException from '@/pages/pda/PDAException';
import PDAQuery from '@/pages/pda/PDAQuery';
import { usePDAAuth } from '@/stores/pdaAuth';

// Admin full pages (dynamically imported)
import { lazy, Suspense } from 'react';
const Lazy = (imp: () => Promise<any>) => {
  const Comp = lazy(imp);
  return <Suspense fallback={<div className="p-8 text-gray-400">加载中...</div>}><Comp /></Suspense>;
};

// Generate all admin module routes from page registry
const adminModuleRoutes = [
  { path: 'orders', page: () => import('@/pages/admin/OrdersPage') },
  { path: 'service-orders', page: () => import('@/pages/admin/ServiceOrdersPage') },
  { path: 'work-orders', page: () => import('@/pages/admin/WorkOrdersPage') },
  { path: 'route-templates', page: () => import('@/pages/admin/RouteTemplatesPage') },
  { path: 'area-groups', page: () => import('@/pages/admin/AreaGroupsPage') },
  { path: 'cargo-types', page: () => import('@/pages/admin/CargoTypesPage') },
  { path: 'transport-modes', page: () => import('@/pages/admin/TransportModesPage') },
  { path: 'container-loadings', page: () => import('@/pages/admin/ContainerLoadingsPage') },
  { path: 'shipping-providers', page: () => import('@/pages/admin/ShippingProvidersPage') },
  { path: 'customs-brokers', page: () => import('@/pages/admin/CustomsBrokersPage') },
  { path: 'customs-points', page: () => import('@/pages/admin/CustomsPointsPage') },
  { path: 'client-accounts', page: () => import('@/pages/admin/ClientAccountsPage') },
  { path: 'client-members', page: () => import('@/pages/admin/ClientMembersPage') },
  { path: 'client-recharges', page: () => import('@/pages/admin/ClientRechargesPage') },
  { path: 'balance-logs', page: () => import('@/pages/admin/BalanceLogsPage') },
  { path: 'client-pricing', page: () => import('@/pages/admin/ClientPricingPage') },
  { path: 'pricing/routes', page: () => import('@/pages/admin/PricingRoutesPage') },
  { path: 'pricing/delivery', page: () => import('@/pages/admin/PricingDeliveryPage') },
  { path: 'pricing/surcharges', page: () => import('@/pages/admin/PricingSurchargesPage') },
  { path: 'pricing/services', page: () => import('@/pages/admin/PricingServicesPage') },
  { path: 'monthly-statements', page: () => import('@/pages/admin/MonthlyStatementsPage') },
  { path: 'customer-addresses', page: () => import('@/pages/admin/CustomerAddressesPage') },
  { path: 'customer-declarants', page: () => import('@/pages/admin/CustomerDeclarantsPage') },
  { path: 'notifications', page: () => import('@/pages/admin/NotificationsPage') },
  { path: 'print-templates', page: () => import('@/pages/admin/PrintTemplatesPage') },
  { path: 'storage', page: () => import('@/pages/admin/StorageConfigPage') },
  { path: 'system/params', page: () => import('@/pages/admin/SystemParamsPage') },
  { path: 'system/api-couriers', page: () => import('@/pages/admin/ApiCouriersPage') },
  { path: 'system/api-customs', page: () => import('@/pages/admin/ApiCustomsPage') },
  { path: 'system/api-notifications', page: () => import('@/pages/admin/ApiNotificationsPage') },
  { path: 'system/api-printers', page: () => import('@/pages/admin/ApiPrintersPage') },
  { path: 'system/api-storage', page: () => import('@/pages/admin/ApiStoragePage') },
  { path: 'system/api-ezway', page: () => import('@/pages/admin/ApiEZWayPage') },
  { path: 'system/api-devices', page: () => import('@/pages/admin/ApiDevicesPage') },
  { path: 'system/brand', page: () => import('@/pages/admin/BrandSettingsPage') },
  { path: 'system/ai-settings', page: () => import('@/pages/admin/AISettingsPage') },
  { path: 'system/ai-chat', page: () => import('@/pages/admin/AIChatPage') },
  { path: 'system/scheduler', page: () => import('@/pages/admin/SchedulerPage') },
  { path: 'system/audit-logs', page: () => import('@/pages/admin/AuditLogsPage') },
  { path: 'system/reports', page: () => import('@/pages/admin/ReportsPage') },
  { path: 'report/order-profit', page: () => import('@/pages/admin/OrderProfitReport') },
  { path: 'report/route-profit', page: () => import('@/pages/admin/RouteProfitReport') },
  { path: 'report/client-profit', page: () => import('@/pages/admin/ClientProfitReport') },
  { path: 'report/service-profit', page: () => import('@/pages/admin/ServiceProfitReport') },
  { path: 'task-monitor', page: () => import('@/pages/admin/TaskMonitorPage') },
  { path: 'warehouse-board', page: () => import('@/pages/admin/WarehouseBoardPage') },
  { path: 'warehouse-console', page: () => import('@/pages/admin/WarehouseConsolePage') },
  { path: 'logistics-tracking', page: () => import('@/pages/admin/LogisticsTrackingPage') },
  // Cross-reference restored: old routes previously missing
  { path: 'ai-exceptions', page: () => import('@/pages/admin/AiExceptionsPage') },
  { path: 'client-permissions', page: () => import('@/pages/admin/ClientPermissionsPage') },
  { path: 'exception-reports', page: () => import('@/pages/admin/ExceptionReportsPage') },
  { path: 'exceptions', page: () => import('@/pages/admin/ExceptionsPage') },
  { path: 'inbound-board', page: () => import('@/pages/admin/InboundBoardPage') },
  { path: 'pda-sessions', page: () => import('@/pages/admin/PDASessionsPage') },
  { path: 'pda-workorder-templates', page: () => import('@/pages/admin/PDAWorkorderTemplatesPage') },
  { path: 'workflow-management', page: () => import('@/pages/admin/WorkflowManagementPage') },
  { path: 'client-members-list', page: () => import('@/pages/admin/ClientMembersPage') },
  { path: 'balance-logs', page: () => import('@/pages/admin/BalanceLogsPage') },
  { path: 'recharge-records', page: () => import('@/pages/admin/RechargeRecordsPage') },
  { path: 'containers', page: () => import('@/pages/admin/ContainersPage') },
  { path: 'client-panel-perms', page: () => import('@/pages/admin/ClientPanelPermsPage') },
  { path: 'printers', page: () => import('@/pages/admin/PrintersPage') },
  { path: 'service-templates', page: () => import('@/pages/admin/ServiceTemplatesPage') },
  { path: 'service-types', page: () => import('@/pages/admin/ServiceTypesPage') },
  { path: 'service-workorders', page: () => import('@/pages/admin/ServiceWorkordersPage') },
  { path: 'client-ledgers', page: () => import('@/pages/admin/ClientLedgersPage') },
  { path: 'system/customs-broker-api', page: () => import('@/pages/admin/CustomsBrokerAPIPage') },
  { path: 'system/logistics-api', page: () => import('@/pages/admin/LogisticsAPIPage') },
  { path: 'system/notification-channels', page: () => import('@/pages/admin/NotificationChannelsPage') },
  { path: 'system/printers', page: () => import('@/pages/admin/PrintersSystemPage') },
  { path: 'system/settings', page: () => import('@/pages/admin/SystemSettingsPage') },
];

const queryClient = new QueryClient();

function ProtectedAdmin() {
  const { user, checkSession } = useAuthStore();
  const [loading, setLoading] = useState(true);
  useEffect(() => { checkSession().finally(() => setLoading(false)); }, []);
  if (loading) return <div className="flex items-center justify-center h-screen"><p className="text-gray-400">加载中...</p></div>;
  if (!user) return <Navigate to="/admin/login" replace />;
  return <DashboardLayout />;
}

function ProtectedPDA() {
  const { operatorId, loading, checkSession } = usePDAAuth();
  useEffect(() => { checkSession(); }, []);
  if (loading) return <div className="flex items-center justify-center h-screen"><p className="text-gray-400">加载中...</p></div>;
  if (!operatorId) return <Navigate to="/pda/login" replace />;
  return <PDALayout />;
}

function ProtectedClient() {
  const { client, loading, checkSession } = useClientAuth();
  useEffect(() => { checkSession(); }, []);
  if (loading) return <div className="flex items-center justify-center h-screen"><p className="text-gray-400">加载中...</p></div>;
  if (!client) return <Navigate to="/client/login" replace />;
  return <ClientLayout />;
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          {/* Admin routes */}
          <Route path="/admin/login" element={<LoginPage />} />
          <Route path="/admin" element={<ProtectedAdmin />}>
            <Route index element={<Navigate to="/admin/dashboard" replace />} />
            <Route path="dashboard" element={<DashboardPage />} />
            <Route path="employees" element={<EmployeesPage />} />
            <Route path="warehouses" element={<WarehousesPage />} />
            <Route path="clients" element={<ClientsPage />} />
            <Route path="roles" element={<RolesPage />} />
            <Route path="parcels" element={<ParcelsPage />} />
            <Route path="carriers" element={<CarriersPage />} />
            <Route path="couriers" element={<CouriersPage />} />
            <Route path="declarants" element={<DeclarantsPage />} />
            {/* Auto-generated admin module routes */}
            {adminModuleRoutes.map((r) => (
              <Route key={r.path} path={r.path} element={Lazy(r.page)} />
            ))}
          </Route>

          {/* PDA routes */}
          <Route path="/pda/login" element={<PDALogin />} />
          <Route path="/pda" element={<ProtectedPDA />}>
            <Route index element={<Navigate to="/pda/dashboard" replace />} />
            <Route path="dashboard" element={<PDADashboard />} />
            <Route path="receive" element={<PDAReceive />} />
            <Route path="weigh" element={<PDAWeigh />} />
            <Route path="putaway" element={<PDAPutaway />} />
            <Route path="pick" element={<PDAPick />} />
            <Route path="pack" element={<PDAPack />} />
            <Route path="load" element={<PDALoad />} />
            <Route path="exception" element={<PDAException />} />
            <Route path="query" element={<PDAQuery />} />
          </Route>

          {/* Client routes */}
          <Route path="/client/login" element={<ClientLogin />} />
          <Route path="/client" element={<ProtectedClient />}>
            <Route index element={<Navigate to="/client/dashboard" replace />} />
            <Route path="dashboard" element={<ClientDashboard />} />
            <Route path="parcels" element={<ClientParcels />} />
            <Route path="orders" element={<ClientOrders />} />
            <Route path="orders/new" element={<ClientOrderNew />} />
            <Route path="orders/:id" element={<ClientOrderDetail />} />
            <Route path="ledger" element={<ClientLedger />} />
            <Route path="declarants" element={<ClientDeclarants />} />
            <Route path="members" element={<ClientMembers />} />
            <Route path="addresses" element={<ClientAddresses />} />
            <Route path="warehouses" element={<ClientWarehouses />} />
            <Route path="couriers" element={<ClientCouriers />} />
            <Route path="services" element={<ClientServices />} />
            <Route path="credentials" element={<ClientCredentials />} />
            <Route path="route-prices" element={<ClientRoutePrices />} />
            <Route path="delivery-fees" element={<ClientDeliveryFees />} />
            <Route path="surcharges" element={<ClientSurcharges />} />
            <Route path="webhooks" element={<ClientWebhooks />} />
            <Route path="weight-dashboard" element={<ClientWeightDashboard />} />
          </Route>

          <Route path="/" element={<Navigate to="/admin/dashboard" replace />} />
          <Route path="*" element={<div className="p-8 text-center text-gray-400">404</div>} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
