import { create } from "zustand"
import axios from "@/api/client"

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
  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
  checkSession: () => Promise<void>
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  loading: true,

  login: async (username, password) => {
    const form = new URLSearchParams()
    form.append("username", username)
    form.append("password", password)
    await axios.post("/admin/login", form)
    await useAuthStore.getState().checkSession()
  },

  logout: async () => {
    await axios.get("/admin/logout")
    set({ user: null })
  },

  checkSession: async () => {
    try {
      const res = await axios.get("/admin/api/me")
      set({ user: res.data, loading: false })
    } catch {
      set({ user: null, loading: false })
    }
  },
}))
