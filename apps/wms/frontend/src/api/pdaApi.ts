import client from './client';

interface PendingItems {
  receive: unknown[];
  putaway: unknown[];
  weigh: unknown[];
  pick: unknown[];
  pack: unknown[];
}

const pdaApi = {
  login: (code: string, pin: string) =>
    client.post('/pda/api/login', { code, pin }),
  logout: () => client.post('/pda/api/logout'),
  me: () => client.get<{ operator_id: number }>('/pda/api/me'),
  dashboard: () => client.get('/pda/api/dashboard'),

  receive: (data: { scan: string; weight: number; length: number; width: number; height: number }) =>
    client.post('/pda/api/receive', data),
  weigh: (data: { scan: string; weight: number }) =>
    client.post('/pda/api/weigh', data),
  putaway: (data: { scan: string; location_barcode: string }) =>
    client.post('/pda/api/putaway', data),
  pick: (data: { order_no: string }) =>
    client.post('/pda/api/pick', data),
  pack: (data: { order_no: string }) =>
    client.post('/pda/api/pack', data),
  load: (data: { container_no: string; order_no: string }) =>
    client.post('/pda/api/load', data),
  exception: (data: { scan: string; reason: string }) =>
    client.post('/pda/api/exception', data),
  query: (data: { scan: string }) =>
    client.post('/pda/api/query', data),
  pending: () => client.get<PendingItems>('/pda/api/pending'),
};

export default pdaApi;
