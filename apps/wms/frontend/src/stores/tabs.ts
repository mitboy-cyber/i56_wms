import { create } from 'zustand';
import type { TabItem } from '@/types';

interface TabState {
  tabs: TabItem[];
  activeTabId: string | null;
  /** Open a tab (or focus existing) */
  openTab: (tab: Omit<TabItem, 'openedAt'>) => void;
  /** Close a tab by id */
  closeTab: (id: string) => void;
  /** Set the active tab */
  setActiveTab: (id: string) => void;
  /** Close all tabs */
  closeAll: () => void;
  /** Close tabs except the given id */
  closeOthers: (id: string) => void;
}

export const useTabStore = create<TabState>((set, get) => ({
  tabs: [],
  activeTabId: null,

  openTab: (tab) => {
    const existing = get().tabs.find((t) => t.id === tab.id);
    if (existing) {
      set({ activeTabId: tab.id });
      return;
    }
    const newTab: TabItem = { ...tab, openedAt: Date.now() };
    set((s) => ({
      tabs: [...s.tabs, newTab],
      activeTabId: tab.id,
    }));
  },

  closeTab: (id) => {
    set((s) => {
      const next = s.tabs.filter((t) => t.id !== id);
      let nextActive = s.activeTabId;
      if (s.activeTabId === id) {
        // Activate the last tab
        nextActive = next.length > 0 ? next[next.length - 1].id : null;
      }
      return { tabs: next, activeTabId: nextActive };
    });
  },

  setActiveTab: (id) => set({ activeTabId: id }),

  closeAll: () => set({ tabs: [], activeTabId: null }),

  closeOthers: (id) => {
    set((s) => ({
      tabs: s.tabs.filter((t) => t.id === id),
      activeTabId: id,
    }));
  },
}));
