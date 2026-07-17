import { create } from 'zustand';

interface ClientAuthState {
  client: string | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<boolean>;
  logout: () => Promise<void>;
  checkSession: () => Promise<void>;
}

export const useClientAuth = create<ClientAuthState>((set) => ({
  client: null,
  loading: true,
  login: async (email, password) => {
    const params = new URLSearchParams();
    params.append('email', email);
    params.append('password', password);
    const res = await fetch('/client/login', {
      method: 'POST',
      body: params,
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      credentials: 'include',
      redirect: 'manual',
    });
    if (res.status === 303 || res.ok) {
      await useClientAuth.getState().checkSession();
      return true;
    }
    return false;
  },
  logout: async () => {
    await fetch('/client/logout', { credentials: 'include' });
    set({ client: null });
  },
  checkSession: async () => {
    try {
      const res = await fetch('/client/api/me', { credentials: 'include' });
      if (!res.ok) throw new Error('no session');
      const data = await res.json();
      // Go backend returns array wrapped in {data: [...]} or plain object
      set({ client: data?.client || data?.data?.client || null, loading: false });
    } catch {
      set({ client: null, loading: false });
    }
  },
}));
