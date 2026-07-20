import { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom"
import { useAuthStore } from "@/stores/auth"
import { Package, Eye, EyeOff, LogIn, AlertCircle } from "lucide-react"

export function LoginPage() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [showPwd, setShowPwd] = useState(false)
  const [remember, setRemember] = useState(false)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)
  const [mounted, setMounted] = useState(false)
  const navigate = useNavigate()
  const { login } = useAuthStore()

  useEffect(() => { setMounted(true) }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    if (!username.trim() || !password.trim()) {
      setError("请输入账号和密码")
      return
    }
    setLoading(true)
    try {
      const ok = await login(username, password)
      if (ok) navigate("/admin/dashboard")
      else setError("账号或密码错误")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex bg-gradient-to-br from-slate-900 via-slate-800 to-emerald-900">
      {/* Left: Branding */}
      <div className="hidden lg:flex lg:w-5/12 flex-col justify-center px-16 text-white relative overflow-hidden">
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_rgba(16,185,129,0.15),transparent_60%)]" />
        <div className={`relative transition-all duration-1000 ${mounted ? "translate-x-0 opacity-100" : "-translate-x-8 opacity-0"}`}>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-12 h-12 bg-emerald-500 rounded-xl flex items-center justify-center shadow-lg shadow-emerald-500/25">
              <Package className="w-7 h-7 text-white" />
            </div>
            <div>
              <h2 className="text-2xl font-semibold tracking-tight">I56 WMS</h2>
              <p className="text-emerald-400 text-sm">Enterprise Warehouse Management</p>
            </div>
          </div>
          <h1 className="text-4xl font-bold leading-tight mb-4">
            智能仓储
            <br />
            <span className="text-emerald-400">管理平台</span>
          </h1>
          <p className="text-slate-400 text-lg leading-relaxed max-w-sm">
            实时追踪、智能分配、精准履约。
            <br />
            让仓库运营更高效。
          </p>
          <div className="mt-12 flex gap-6 text-sm text-slate-500">
            <div>
              <div className="text-2xl font-bold text-white">99.9%</div>
              <div>系统可用性</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-white">毫秒级</div>
              <div>响应延迟</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-white">7×24</div>
              <div>运行保障</div>
            </div>
          </div>
        </div>
      </div>

      {/* Right: Login Form */}
      <div className="flex-1 flex items-center justify-center px-6 bg-white lg:rounded-l-[3rem] shadow-2xl">
        <div className={`w-full max-w-md transition-all duration-700 delay-200 ${mounted ? "translate-y-0 opacity-100" : "translate-y-4 opacity-0"}`}>
          {/* Mobile branding */}
          <div className="lg:hidden flex items-center gap-3 mb-10 justify-center">
            <div className="w-10 h-10 bg-emerald-500 rounded-xl flex items-center justify-center">
              <Package className="w-6 h-6 text-white" />
            </div>
            <h2 className="text-xl font-bold text-slate-800">I56 WMS</h2>
          </div>

          <div className="mb-8">
            <h2 className="text-2xl font-bold text-slate-900">欢迎回来</h2>
            <p className="text-slate-500 mt-1">登录管理后台</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            {/* Error */}
            {error && (
              <div className="flex items-center gap-2 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm animate-shake">
                <AlertCircle className="w-4 h-4 flex-shrink-0" />
                <span>{error}</span>
              </div>
            )}

            {/* Username */}
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1.5">
                员工编号
              </label>
              <input
                type="text"
                value={username}
                onChange={(e) => { setUsername(e.target.value); setError("") }}
                className="w-full px-4 py-3 border border-slate-300 rounded-xl text-sm
                  focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500
                  placeholder:text-slate-400 transition-colors"
                placeholder="请输入员工编号"
                autoComplete="username"
                autoFocus
              />
            </div>

            {/* Password */}
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1.5">
                密码
              </label>
              <div className="relative">
                <input
                  type={showPwd ? "text" : "password"}
                  value={password}
                  onChange={(e) => { setPassword(e.target.value); setError("") }}
                  className="w-full px-4 py-3 pr-12 border border-slate-300 rounded-xl text-sm
                    focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500
                    placeholder:text-slate-400 transition-colors"
                  placeholder="请输入密码"
                  autoComplete="current-password"
                />
                <button
                  type="button"
                  onClick={() => setShowPwd(!showPwd)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors"
                  tabIndex={-1}
                >
                  {showPwd ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
            </div>

            {/* Remember + Forgot */}
            <div className="flex items-center justify-between">
              <label className="flex items-center gap-2 text-sm text-slate-600 cursor-pointer">
                <input
                  type="checkbox"
                  checked={remember}
                  onChange={(e) => setRemember(e.target.checked)}
                  className="rounded border-slate-300 text-emerald-600 focus:ring-emerald-500"
                />
                保持登录状态
              </label>
              <a href="#" className="text-sm text-emerald-600 hover:text-emerald-700 transition-colors">
                忘记密码?
              </a>
            </div>

            {/* Submit */}
            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 bg-emerald-600 hover:bg-emerald-700 disabled:bg-emerald-400
                text-white font-medium rounded-xl transition-all duration-200
                focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2
                shadow-lg shadow-emerald-600/25 hover:shadow-emerald-600/40
                active:scale-[0.98]"
            >
              {loading ? (
                <span className="flex items-center justify-center gap-2">
                  <svg className="animate-spin w-4 h-4" viewBox="0 0 24 24" fill="none">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                  </svg>
                  登录中...
                </span>
              ) : (
                <span className="flex items-center justify-center gap-2">
                  <LogIn className="w-4 h-4" />
                  登录
                </span>
              )}
            </button>
          </form>

          {/* Footer */}
          <p className="mt-8 text-center text-xs text-slate-400">
            I56 Framework 1.0 LTS · Enterprise Application Platform
          </p>
        </div>
      </div>
    </div>
  )
}
