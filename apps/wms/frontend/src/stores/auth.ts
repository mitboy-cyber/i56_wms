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
  token: string | null
  login: (username: string, password: string) => Promise<boolean>
  logout: () => Promise<void>
  checkSession: () => Promise<void>
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  loading: true,
  token: localStorage.getItem("admin_token"),

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

    const data = await res.json()
    if (data.token) {
      localStorage.setItem("admin_token", data.token)
      set({ token: data.token })
      // Fetch user info with Bearer token
      try {
        const meRes = await fetch("/admin/api/me", {
          headers: { Authorization: `Bearer ${data.token}` },
          credentials: "include",
        })
        if (meRes.ok) {
          const user: User = await meRes.json()
          set({ user, loading: false })
          return true
        }
      } catch { /* fallback */ }
    }

    // Fallback: try cookie-based session
    try {
      const meRes = await fetch("/admin/api/me", { credentials: "include" })
      if (meRes.ok) {
        const user: User = await meRes.json()
        set({ user, loading: false })
        return true
      }
    } catch { /* fallback failed */ }

    return false
  },

  logout: async () => {
    localStorage.removeItem("admin_token")
    await fetch("/admin/logout", { credentials: "include" })
    set({ user: null, token: null })
  },

  checkSession: async () => {
    const token = get().token || localStorage.getItem("admin_token")
    try {
      const headers: Record<string, string> = {}
      if (token) headers["Authorization"] = `Bearer ${token}`
      const res = await fetch("/admin/api/me", { headers, credentials: "include" })
      if (!res.ok) throw new Error("no session")
      const data: User = await res.json()
      set({ user: data, token: token || undefined, loading: false })
    } catch {
      set({ user: null, loading: false })
    }
  },
}))
