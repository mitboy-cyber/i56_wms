import client from './client';

export interface Parcel {
  id: number;
  tracking_number: string;
  product_name: string;
  status: string;
  actual_weight: number;
  warehouse_id: number;
  client_id: number;
}

export interface Order {
  id: number;
  order_no: string;
  recipient_name: string;
  parcel_count: number;
  route_id: number;
  total_actual_weight: number;
  total_price: number;
  status: string;
  created_at: string;
}

export interface Warehouse {
  id: number;
  name: string;
  code: string;
  address: string;
  contact: string;
  phone: string;
}

export interface Declarant {
  id: number;
  name: string;
  id_number: string;
  type: string;
  phone: string;
  client_id: number;
}

export interface Member {
  id: number;
  name: string;
  client_id: number;
}

export interface Address {
  id: number;
  recipient_name: string;
  phone: string;
  city: string;
  district: string;
  address: string;
}

export interface RoutePrice {
  route_name: string;
  transport_type: string;
  base_weight_price: number;
  base_volume_price: number;
}

const clientApi = {
  // Dashboard
  dashboard: () => client.get('/client/api/dashboard'),
  me: () => client.get('/client/api/me'),

  // Parcels
  parcels: () => client.get<Parcel[]>('/client/api/parcels'),
  predeclare: (data: { tracking_number: string; product_name: string; warehouse_id?: number; courier_code?: string }) =>
    client.post<Parcel>('/client/api/parcels/predeclare', data),

  // Orders
  orders: () => client.get<Order[]>('/client/api/orders'),
  orderDetail: (id: string) => client.get<Order>(`/client/api/orders/${id}`),

  // Ledger
  ledger: () => client.get('/client/api/ledger'),

  // Declarants
  declarants: () => client.get<Declarant[]>('/client/api/declarants'),

  // Members
  members: () => client.get<Member[]>('/client/api/members'),

  // Addresses
  addresses: () => client.get<Address[]>('/client/api/addresses'),

  // Warehouses
  warehouses: () => client.get<Warehouse[]>('/client/api/warehouses'),

  // Route prices
  routePrices: () => client.get<RoutePrice[]>('/client/api/route-prices'),

  // Service orders
  serviceOrders: () => client.get('/client/api/service-orders'),

  // Webhooks
  webhooks: () => client.get('/client/api/webhooks'),

  // Credentials
  credentials: () => client.get('/client/api/credentials'),

  // Delivery fees
  deliveryFees: () => client.get('/client/api/delivery-fees'),

  // Surcharges
  surcharges: () => client.get('/client/api/surcharges'),

  // Couriers
  couriers: () => client.get('/client/api/couriers'),
};

export default clientApi;
