package main

import (
	"html/template"
	"net/http"
	"github.com/i56/framework/core/router"
	weightDomain "github.com/i56/modules/weight/domain"
)

func registerWeightUIRoutes(r *router.Router, adminTmpl map[string]*template.Template, weightRepo *weightDomain.MemWeightRepo) {
	r.GET("/admin/weight-records", func(w http.ResponseWriter, req *http.Request) {
		records, total := weightRepo.List(1, 0, 50)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := weighAdminHTML(records, total)
		w.Write([]byte(html))
	})
}

func weighAdminHTML(records []weightDomain.WeightRecord, total int64) string {
	h := `<!DOCTYPE html><html lang="zh-TW"><head><meta charset="UTF-8"><title>称重记录管理</title>
<link rel="stylesheet" href="/static/css/i56-bdl.css"><script src="/static/js/i56-theme.js"></script>
</head><body style="padding:16px">
<h5 style="color:var(--i56-brand);font-size:15px;margin-bottom:12px">称重记录管理 <span class="i56-badge i56-badge-brand">` + i64toa(total) + `</span></h5>
<table style="width:100%;border-collapse:collapse;background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;overflow:hidden;font-size:12px"><thead><tr style="background:var(--i56-bg-base)">
<th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">ID</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">快递单号</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">快递公司</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">会员ID</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">平台</th>
<th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">重量(kg)</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">尺寸(cm)</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">体积(cm³)</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">包裹数</th>
<th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">类型</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">品名</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">备注</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">时间</th>
</tr></thead><tbody>`
	for _, r := range records {
		v := int64(r.Volume)
		h += `<tr><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + i64toa(r.ID) + `</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + r.TrackingNumber + `</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + r.CourierCompany +
			`</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + r.MemberID + `</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + r.Platform +
			`</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + f64toa(r.Weight) + `</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + i64toa(int64(r.Length)) + `×` + i64toa(int64(r.Width)) + `×` + i64toa(int64(r.Height)) + `</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + i64toa(v) + `</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + i64toa(int64(r.ParcelCount)) +
			`</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)"><span class="i56-badge i56-badge-brand">` + string(r.ParcelType) + `</span></td>` +
			`<td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + r.ProductName + `</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)">` + r.Remark + `</td><td style="padding:8px 12px;border-bottom:1px solid var(--i56-border)"><small>` + r.CreatedAt.Format("01-02 15:04") + `</small></td></tr>`
	}
	if len(records) == 0 {
		h += `<tr><td colspan="13" style="padding:32px;text-align:center;color:var(--i56-text-secondary)">暂无称重记录。通过API POST /api/v1/weight-records 创建。</td></tr>`
	}
	h += `</tbody></table></body></html>`
	return h
}

func i64toa(n int64) string {
	if n == 0 { return "0" }
	s := ""
	for n > 0 { s = string(rune('0'+n%10)) + s; n /= 10 }
	return s
}

func f64toa(f float64) string {
	w := int(f)
	fr := int((f - float64(w) + 0.005) * 100)
	s := i64toa(int64(w)) + "."
	if fr < 10 { s += "0" }
	s += i64toa(int64(fr))
	return s
}
