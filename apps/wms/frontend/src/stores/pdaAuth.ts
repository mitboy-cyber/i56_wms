import { create } from 'zustand';
import pdaApi from '@/api/pdaApi';

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
      await pdaApi.login(code, pin);
      await usePDAAuth.getState().checkSession();
      return true;
    } catch {
      return false;
    }
  },
  logout: async () => {
    await pdaApi.logout();
    set({ operatorId: null });
  },
  checkSession: async () => {
    try {
      const res = await pdaApi.me();
      set({ operatorId: res.data.operator_id, loading: false });
    } catch {
      set({ operatorId: null, loading: false });
    }
  },
}));
