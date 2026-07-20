import axios from "axios"

const client = axios.create({
  withCredentials: true,
})

// Track whether a redirect is already in progress
let redirecting = false

client.interceptors.response.use(
  (r) => r,
  async (err) => {
    const status = err.response?.status
    const url = err.config?.url || ""
    const isLoginPage = window.location.pathname === "/admin/login" || window.location.pathname === "/client/login"
    const isAuthEndpoint = url.includes("/api/me") || url.includes("/login")
    const isClientPage = window.location.pathname.startsWith("/client")

    // Only redirect on 401 when NOT on login page, NOT auth endpoint, and NOT already redirecting
    if (status === 401 && !isLoginPage && !isAuthEndpoint && !redirecting) {
      // Verify session is truly invalid before redirecting
      try {
        const checkUrl = isClientPage ? "/client/api/me" : "/admin/api/me"
        const checkRes = await axios.get(checkUrl, { withCredentials: true })
        if (checkRes.status !== 200) throw new Error("invalid")
        // Session still valid — this 401 was just for this specific resource
        return Promise.reject(err)
      } catch {
        // Session truly invalid — redirect
        redirecting = true
        window.location.replace(isClientPage ? "/client/login" : "/admin/login")
      }
    }
    return Promise.reject(err)
  }
)

export default client
