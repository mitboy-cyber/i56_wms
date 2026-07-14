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
import { ClientsPage } from '@/pages/Clients';
import { RolesPage } from '@/pages/Roles';
import { ParcelsPage } from '@/pages/Parcels';
import { CarriersPage, CouriersPage, DeclarantsPage } from '@/pages/Transport';

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

const queryClient = new QueryClient();

function ProtectedAdmin() {
  const { user, checkSession } = useAuthStore();
  const [loading, setLoading] = useState(true);
  useEffect(() => { checkSession().finally(() => setLoading(false)); }, []);
  if (loading) return <div className="flex items-center justify-center h-screen"><p className="text-gray-400">加载中...</p></div>;
  if (!user) return <Navigate to="/admin/login" replace />;
  return <DashboardLayout />;
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
          </Route>

          {/* Client routes */}
          <Route path="/client/login" element={<ClientLogin />} />
          <Route path="/client" element={<ProtectedClient />}>
            <Route index element={<Navigate to="/client/dashboard" replace />} />
            <Route path="dashboard" element={<ClientDashboard />} />
            <Route path="parcels" element={<ClientParcels />} />
            <Route path="orders" element={<ClientOrders />} />
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
          </Route>

          <Route path="/" element={<Navigate to="/admin/dashboard" replace />} />
          <Route path="*" element={<div className="p-8 text-center text-gray-400">404</div>} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
