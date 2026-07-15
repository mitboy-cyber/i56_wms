import axios from "axios"

const client = axios.create({
  withCredentials: true,
  headers: { "Content-Type": "application/x-www-form-urlencoded" },
})

// Track whether a redirect is already in progress to prevent loops
let redirecting = false

client.interceptors.response.use(
  (r) => r,
  (err) => {
    const url = err.config?.url || ""
    const isLoginPage = window.location.pathname === "/admin/login"
    // Never redirect on auth-related endpoints (they're expected to fail sometimes)
    const isAuthEndpoint = url.includes("/admin/api/me") || url.includes("/admin/login")

    if (
      (err.response?.status === 401 || err.response?.status === 303) &&
      !isLoginPage && !isAuthEndpoint && !redirecting
    ) {
      redirecting = true
      window.location.replace("/admin/login")
    }
    return Promise.reject(err)
  }
)

export default client
