import axios from "axios"

const client = axios.create({
  withCredentials: true,
})

// Inject Bearer token from localStorage if available
client.interceptors.request.use((config) => {
  const token = localStorage.getItem("admin_token")
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Track redirect state
let redirecting = false

client.interceptors.response.use(
  (r) => r,
  async (err) => {
    const status = err.response?.status
    const url = err.config?.url || ""
    const isLoginPage = window.location.pathname === "/admin/login" || window.location.pathname === "/client/login"
    const isAuthEndpoint = url.includes("/api/me") || url.includes("/login")

    if (status === 401 && !isLoginPage && !isAuthEndpoint && !redirecting) {
      redirecting = true
      localStorage.removeItem("admin_token")
      window.location.replace("/admin/login")
    }
    return Promise.reject(err)
  }
)

export default client
