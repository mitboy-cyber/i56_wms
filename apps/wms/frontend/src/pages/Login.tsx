import { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom"
import { useAuthStore } from "@/stores/auth"

export function LoginPage() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [showPwd, setShowPwd] = useState(false)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)
  const [mounted, setMounted] = useState(false)
  const navigate = useNavigate()
  const { login } = useAuthStore()

  useEffect(() => { setMounted(true) }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!username || !password) {
      setError("请输入员工编号和密码")
      return
    }
    setLoading(true)
    setError("")
    try {
      const ok = await login(username, password)
      if (ok) navigate("/admin")
      else setError("账号或密码错误")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ minHeight: "100vh", display: "flex", background: "linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #064e3b 100%)" }}>
      {/* Brand panel */}
      <div style={{ display: "none", flexDirection: "column", justifyContent: "center", padding: "0 64px", color: "white", width: "42%", position: "relative", overflow: "hidden" }}
        className="lg:flex">
        <div style={{ position: "absolute", inset: 0, background: "radial-gradient(ellipse at top right, rgba(16,185,129,0.15), transparent 60%)" }} />
        <div style={{ position: "relative", transition: "all 1s", transform: mounted ? "translateX(0)" : "translateX(-32px)", opacity: mounted ? 1 : 0 }}>
          <div style={{ display: "flex", alignItems: "center", gap: 12, marginBottom: 32 }}>
            <div style={{ width: 48, height: 48, background: "#10b981", borderRadius: 12, display: "flex", alignItems: "center", justifyContent: "center", boxShadow: "0 4px 16px rgba(16,185,129,0.3)" }}>
              <span style={{ fontSize: 24, fontWeight: "bold", color: "white" }}>I</span>
            </div>
            <div>
              <h2 style={{ fontSize: 24, fontWeight: 600, margin: 0 }}>I56 WMS</h2>
              <p style={{ color: "#34d399", fontSize: 13, margin: 0 }}>Enterprise Warehouse Management</p>
            </div>
          </div>
          <h1 style={{ fontSize: 40, fontWeight: "bold", lineHeight: 1.2, marginBottom: 16 }}>
            智能仓储<br /><span style={{ color: "#34d399" }}>管理平台</span>
          </h1>
          <p style={{ color: "#94a3b8", fontSize: 17, maxWidth: 320, lineHeight: 1.6 }}>
            实时追踪、智能分配、精准履约。<br />让仓库运营更高效。
          </p>
          <div style={{ display: "flex", gap: 32, marginTop: 48, fontSize: 13, color: "#64748b" }}>
            <div><div style={{ fontSize: 24, fontWeight: "bold", color: "white" }}>99.9%</div><div>系统可用性</div></div>
            <div><div style={{ fontSize: 24, fontWeight: "bold", color: "white" }}>毫秒级</div><div>响应延迟</div></div>
            <div><div style={{ fontSize: 24, fontWeight: "bold", color: "white" }}>7×24</div><div>运行保障</div></div>
          </div>
        </div>
      </div>

      {/* Login form */}
      <div style={{ flex: 1, display: "flex", alignItems: "center", justifyContent: "center", padding: 24, background: "white", borderTopLeftRadius: "3rem", borderBottomLeftRadius: "3rem", boxShadow: "-8px 0 32px rgba(0,0,0,0.1)" }}>
        <div style={{ width: "100%", maxWidth: 400, transition: "all .7s .2s", transform: mounted ? "translateY(0)" : "translateY(16px)", opacity: mounted ? 1 : 0 }}>
          <h2 style={{ fontSize: 28, fontWeight: "bold", color: "#1f2937", marginBottom: 4 }}>欢迎回来</h2>
          <p style={{ color: "#6b7280", marginBottom: 32 }}>登录管理后台</p>
          {error && <div style={{ background: "#fef2f2", color: "#dc2626", padding: "10px 16px", borderRadius: 8, marginBottom: 16, fontSize: 14 }}>{error}</div>}
          <form onSubmit={handleSubmit}>
            <div style={{ marginBottom: 16 }}>
              <label style={{ display: "block", fontSize: 13, fontWeight: 500, color: "#374151", marginBottom: 6 }}>员工编号</label>
              <input type="text" value={username} onChange={e => setUsername(e.target.value)}
                placeholder="请输入员工编号" autoFocus
                style={{ width: "100%", padding: "10px 14px", border: "1px solid #d1d5db", borderRadius: 8, fontSize: 15, outline: "none", boxSizing: "border-box" }} />
            </div>
            <div style={{ marginBottom: 16 }}>
              <label style={{ display: "block", fontSize: 13, fontWeight: 500, color: "#374151", marginBottom: 6 }}>密码</label>
              <div style={{ position: "relative" }}>
                <input type={showPwd ? "text" : "password"} value={password} onChange={e => setPassword(e.target.value)}
                  placeholder="请输入密码"
                  style={{ width: "100%", padding: "10px 44px 10px 14px", border: "1px solid #d1d5db", borderRadius: 8, fontSize: 15, outline: "none", boxSizing: "border-box" }} />
                <button type="button" onClick={() => setShowPwd(!showPwd)}
                  style={{ position: "absolute", right: 10, top: "50%", transform: "translateY(-50%)", background: "none", border: "none", cursor: "pointer", color: "#9ca3af", fontSize: 18, padding: 4 }}>
                  {showPwd ? "🙈" : "👁"}
                </button>
              </div>
            </div>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 20 }}>
              <label style={{ display: "flex", alignItems: "center", gap: 6, fontSize: 13, color: "#6b7280", cursor: "pointer" }}>
                <input type="checkbox" style={{ accentColor: "#10b981" }} /> 保持登录状态
              </label>
              <a href="#" style={{ fontSize: 13, color: "#10b981", textDecoration: "none" }}>忘记密码?</a>
            </div>
            <button type="submit" disabled={loading}
              style={{ width: "100%", padding: 12, background: loading ? "#6ee7b7" : "#10b981", color: "white", border: "none", borderRadius: 8, fontSize: 16, fontWeight: 600, cursor: loading ? "not-allowed" : "pointer" }}>
              {loading ? "登录中..." : "登录"}
            </button>
          </form>
          <p style={{ textAlign: "center", color: "#9ca3af", fontSize: 12, marginTop: 24 }}>I56 Framework 1.0 LTS · Enterprise Application Platform</p>
        </div>
      </div>
    </div>
  )
}
