import axios from "axios"

const client = axios.create({
  withCredentials: true,
  headers: { "Content-Type": "application/x-www-form-urlencoded" },
})

client.interceptors.response.use(
  (r) => r,
  (err) => {
    if (err.response?.status === 401 || err.response?.status === 303) {
      window.location.href = "/admin/login"
    }
    return Promise.reject(err)
  }
)

export default client
