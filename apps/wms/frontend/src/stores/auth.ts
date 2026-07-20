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

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  loading: true,

  login: async (username: string, password: string): Promise<boolean> => {
    const params = new URLSearchParams()
    params.append("username", username)
    params.append("password", password)
    const res = await fetch("/admin/login", {
      method: "POST",
      body: params,
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      credentials: "include",
    })
    if (!res.ok) return false

    // Retry checkSession up to 3 times with increasing delays
    for (let i = 0; i < 3; i++) {
      await new Promise<void>(r => setTimeout(r, 100 * (i + 1)))
      try {
        const meRes = await fetch("/admin/api/me", { credentials: "include" })
        if (meRes.ok) {
          const data: User = await meRes.json()
          set({ user: data, loading: false })
          return true
        }
      } catch { /* retry */ }
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
      const data: User = await res.json()
      set({ user: data, loading: false })
    } catch {
      set({ user: null, loading: false })
    }
  },
}))
