import client from './client';

export interface AdminListItem {
  id: number;
  [key: string]: any;
}

const adminApi = {
  // Core CRUD
  warehouses: {
    list: () => client.get('/admin/api/warehouses'),
    create: (data: any) => client.post('/admin/api/warehouses', data),
    update: (id: number, data: any) => client.put(`/admin/api/warehouses/${id}`, data),
    delete: (id: number) => client.delete(`/admin/api/warehouses/${id}`),
  },
  parcels: {
    list: () => client.get('/admin/api/parcels'),
    create: (data: any) => client.post('/admin/api/parcels', data),
  },
  orders: {
    list: () => client.get('/admin/api/orders'),
    create: (data: any) => client.post('/admin/api/orders', data),
    detail: (id: string) => client.get(`/admin/api/orders/${id}`),
  },
  clients: {
    list: () => client.get('/admin/api/clients'),
    create: (data: any) => client.post('/admin/api/clients', data),
  },
  employees: {
    list: () => client.get('/admin/api/employees'),
    create: (data: any) => client.post('/admin/api/employees', data),
  },
  roles: {
    list: () => client.get('/admin/api/roles'),
  },
  carriers: {
    list: () => client.get('/admin/api/carriers'),
  },
  declarants: {
    list: () => client.get('/admin/api/declarants'),
  },
  members: {
    list: () => client.get('/admin/api/members'),
  },
  addresses: {
    list: () => client.get('/admin/api/addresses'),
  },
  ledger: {
    list: () => client.get('/admin/api/ledger'),
  },
  serviceOrders: {
    list: () => client.get('/admin/api/service-orders'),
  },
  workOrders: {
    list: () => client.get('/admin/api/work-orders'),
  },
  credentials: {
    list: () => client.get('/admin/api/credentials'),
  },
  pricing: {
    routes: () => client.get('/admin/api/pricing/routes'),
    delivery: () => client.get('/admin/api/pricing/delivery'),
    surcharges: () => client.get('/admin/api/pricing/surcharges'),
  },
};

export default adminApi;
