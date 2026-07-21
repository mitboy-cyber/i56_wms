import { useState } from "react"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import client from "@/api/client"

export default function ClientLedgersPage() {
  const qc = useQueryClient()
  const [showRecharge, setShowRecharge] = useState(false)
  const [amount, setAmount] = useState("")
  const [method, setMethod] = useState("银行转账")

  const { data: entries = [] } = useQuery<any[]>({
    queryKey: ["client-ledgers"],
    queryFn: () => client.get("/admin/api/client-ledgers?client_id=1").then(r => r.data),
    retry: false,
  })

  const currentBalance = entries.length > 0 ? entries[0].balance_after : 0

  const rechargeMut = useMutation({
    mutationFn: (body: { client_id: number; amount: number; method: string }) =>
      client.post("/admin/api/ledger-recharge", body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["client-ledgers"] })
      setShowRecharge(false); setAmount("")
    },
  })

  return (
    <div>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
        <h1 style={{ fontSize: 20, fontWeight: "bold" }}>余额日志</h1>
        <button onClick={() => setShowRecharge(true)} style={{ background: "#10b981", color: "white", border: "none", borderRadius: 6, padding: "8px 20px", cursor: "pointer", fontSize: 14, fontWeight: 600 }}>
          + 客户充值
        </button>
      </div>

      {/* Balance card */}
      <div style={{ background: "linear-gradient(135deg, #667eea, #764ba2)", color: "white", borderRadius: 10, padding: "20px 24px", marginBottom: 16 }}>
        <div style={{ fontSize: 13, opacity: .85 }}>当前余额</div>
        <div style={{ fontSize: 32, fontWeight: "bold", margin: "4px 0" }}>
          ¥ {currentBalance.toLocaleString("zh-CN", { minimumFractionDigits: 2 })}
        </div>
        <div style={{ fontSize: 12, opacity: .7 }}>共 {entries.length} 条交易记录</div>
      </div>

      {/* Recharge Modal */}
      {showRecharge && (
        <div style={{ position: "fixed", inset: 0, background: "rgba(0,0,0,0.4)", display: "flex", alignItems: "center", justifyContent: "center", zIndex: 100 }}>
          <div style={{ background: "white", borderRadius: 12, padding: 24, width: 400 }}>
            <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
              <h2 style={{ fontSize: 18, fontWeight: "bold" }}>客户充值</h2>
              <button onClick={() => setShowRecharge(false)} style={{ background: "none", border: "none", fontSize: 20, color: "#9ca3af", cursor: "pointer" }}>✕</button>
            </div>
            <div style={{ marginBottom: 12 }}>
              <label style={{ fontSize: 13, color: "#374151", display: "block", marginBottom: 4 }}>充值金额</label>
              <input type="number" value={amount} onChange={e => setAmount(e.target.value)} placeholder="输入金额"
                style={{ width: "100%", padding: "10px 12px", border: "1px solid #d1d5db", borderRadius: 6, fontSize: 16, boxSizing: "border-box" }} />
            </div>
            <div style={{ marginBottom: 16 }}>
              <label style={{ fontSize: 13, color: "#374151", display: "block", marginBottom: 4 }}>支付方式</label>
              <select value={method} onChange={e => setMethod(e.target.value)}
                style={{ width: "100%", padding: "10px 12px", border: "1px solid #d1d5db", borderRadius: 6, fontSize: 14 }}>
                <option>银行转账</option>
                <option>微信支付</option>
                <option>支付宝</option>
                <option>现金</option>
              </select>
            </div>
            <button onClick={() => {
                const amt = parseFloat(amount)
                if (amt > 0) rechargeMut.mutate({ client_id: 1, amount: amt, method })
              }} disabled={!amount || rechargeMut.isPending}
              style={{ width: "100%", padding: "10px", background: "#10b981", color: "white", border: "none", borderRadius: 6, fontSize: 14, fontWeight: 600, cursor: "pointer", opacity: !amount ? .5 : 1 }}>
              {rechargeMut.isPending ? "处理中..." : "确认充值"}
            </button>
          </div>
        </div>
      )}

      {/* Ledger table */}
      <div style={{ background: "white", borderRadius: 8, border: "1px solid #e5e7eb", overflow: "auto" }}>
        <table style={{ width: "100%", borderCollapse: "collapse" }}>
          <thead>
            <tr>
              <th style={th}>#</th>
              <th style={th}>类型</th>
              <th style={th}>金额</th>
              <th style={th}>余额</th>
              <th style={th}>描述</th>
              <th style={th}>时间</th>
            </tr>
          </thead>
          <tbody>
            {entries.map((e: any, i: number) => (
              <tr key={e.id || i} style={{ background: i % 2 === 0 ? "#fafafa" : "white" }}>
                <td style={td}>{i + 1}</td>
                <td style={td}>
                  <span style={{ fontSize: 12, padding: "2px 8px", borderRadius: 10, background: e.type === "recharge" ? "#dcfce7" : "#fef2f2", color: e.type === "recharge" ? "#16a34a" : "#ef4444", fontWeight: 600 }}>
                    {e.type === "recharge" ? "充值" : e.type === "order_deduct" ? "扣费" : e.type}
                  </span>
                </td>
                <td style={{ ...td, fontFamily: "monospace", color: e.amount > 0 ? "#16a34a" : "#ef4444", fontWeight: 600 }}>
                  {e.amount > 0 ? "+" : ""}{e.amount?.toFixed(2)}
                </td>
                <td style={{ ...td, fontFamily: "monospace", fontWeight: 600 }}>
                  ¥ {e.balance_after?.toLocaleString?.("zh-CN", { minimumFractionDigits: 2 }) || e.balance_after}
                </td>
                <td style={td}>{e.description || "-"}</td>
                <td style={{ ...td, color: "#9ca3af", fontSize: 12 }}>{e.created_at?.slice?.(0, 19) || e.created_at || "-"}</td>
              </tr>
            ))}
            {entries.length === 0 && (
              <tr><td colSpan={6} style={{ ...td, textAlign: "center", color: "#9ca3af" }}>暂无交易记录</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}

const th: React.CSSProperties = { padding: "10px 12px", textAlign: "left", borderBottom: "1px solid #e5e7eb", color: "#6b7280", fontWeight: 600, fontSize: 13 }
const td: React.CSSProperties = { padding: "10px 12px", borderBottom: "1px solid #f3f4f6", fontSize: 13 }
