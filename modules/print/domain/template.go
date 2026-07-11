package domain
import "time"

type PrintTemplate struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // waybill | customs | carrier
	Content   string    `json:"content"` // HTML template
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DefaultTemplates() []PrintTemplate {
	return []PrintTemplate{
		{Name:"标准面单",Type:"waybill",Content:`<div style="font-family:sans-serif;padding:10mm"><h3>{{.OrderNo}}</h3><p>收件人:{{.RecipientName}}</p><p>地址:{{.Address}}</p><p>重量:{{.Weight}}kg</p><p>承运商:{{.Carrier}}</p></div>`},
		{Name:"清关单",Type:"customs",Content:`<div style="font-family:sans-serif"><h3>海关申报单</h3><p>申报人:{{.Declarant}}</p><p>品名:{{.ProductName}}</p><p>数量:{{.Qty}}</p><p>价值:¥{{.Value}}</p></div>`},
		{Name:"承运商面单",Type:"carrier",Content:`<div style="font-family:sans-serif"><h3>{{.CarrierName}}</h3><p>单号:{{.CarrierTN}}</p><p>目的:{{.Destination}}</p><p>件数:{{.Pieces}}</p></div>`},
	}
}
