import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { useEffect } from "react"
import { useAuthStore } from "@/stores/auth"
import { DashboardLayout } from "@/layouts/DashboardLayout"
import { LoginPage } from "@/pages/Login"
import { DashboardPage } from "@/pages/Dashboard"
import { EmployeesPage } from "@/pages/Employees"
import { WarehousesPage } from "@/pages/Warehouses"
import { ClientsPage } from "@/pages/Clients"
import { RolesPage } from "@/pages/Roles"
import { ParcelsPage } from "@/pages/Parcels"
import { CarriersPage, CouriersPage, DeclarantsPage } from "@/pages/Transport"

const queryClient = new QueryClient()

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user, loading, checkSession } = useAuthStore()

  useEffect(() => { checkSession() }, [])

  if (loading) return <div className="flex items-center justify-center h-screen">加载中...</div>
  if (!user) return <Navigate to="/admin/login" replace />

  return <>{children}</>
}

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/admin/login" element={<LoginPage />} />
          <Route
            path="/admin"
            element={
              <ProtectedRoute>
                <DashboardLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Navigate to="dashboard" replace />} />
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
          <Route path="*" element={<Navigate to="/admin" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
