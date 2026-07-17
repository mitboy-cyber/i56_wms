import { create } from 'zustand';

interface PDAAuthState {
  operatorId: number | null;
  loading: boolean;
  login: (code: string, pin: string) => Promise<boolean>;
  logout: () => Promise<void>;
  checkSession: () => Promise<void>;
}

export const usePDAAuth = create<PDAAuthState>((set) => ({
  operatorId: null,
  loading: true,
  login: async (code, pin) => {
    try {
      const res = await fetch('/pda/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code, pin }),
        credentials: 'include',
      });
      if (!res.ok) return false;
      await usePDAAuth.getState().checkSession();
      return true;
    } catch {
      return false;
    }
  },
  logout: async () => {
    await fetch('/pda/api/logout', { method: 'POST', credentials: 'include' });
    set({ operatorId: null });
  },
  checkSession: async () => {
    try {
      const res = await fetch('/pda/api/me', { credentials: 'include' });
      if (!res.ok) throw new Error('no session');
      const data = await res.json();
      // Go returns {"code":200,"data":{"operator_id":1}}
      set({ operatorId: data?.data?.operator_id || data?.operator_id || null, loading: false });
    } catch {
      set({ operatorId: null, loading: false });
    }
  },
}));
