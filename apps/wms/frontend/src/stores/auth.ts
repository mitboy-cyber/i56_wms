import { create } from "zustand"

interface User {
  id: number
  username: string
  real_name: string
  role_id: number
  role_name: string
}

interface AuthState {
  user: User | null
  loading: boolean
  login: (username: string, password: string) => Promise<boolean>
  logout: () => Promise<void>
  checkSession: () => Promise<void>
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  loading: true,

  login: async (username, password) => {
    const params = new URLSearchParams()
    params.append("username", username)
    params.append("password", password)
    const res = await fetch("/admin/login", {
      method: "POST",
      body: params,
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      credentials: "include",
      redirect: "manual",
    })
    if (res.ok) {
      // Small delay to ensure cookie is set before checkSession
      await new Promise(r => setTimeout(r, 50))
      await useAuthStore.getState().checkSession()
      return true
    }
    return false
  },

  logout: async () => {
    await fetch("/admin/logout", { credentials: "include" })
    set({ user: null })
  },

  checkSession: async () => {
    try {
      const res = await fetch("/admin/api/me", { credentials: "include" })
      if (!res.ok) throw new Error("no session")
      const data = await res.json()
      set({ user: data, loading: false })
    } catch {
      set({ user: null, loading: false })
    }
  },
}))
