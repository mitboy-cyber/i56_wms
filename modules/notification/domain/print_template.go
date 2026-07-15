package domain

import "time"

// PrintTemplateType represents the type of print template.
type PrintTemplateType string

const (
	PrintTypeWaybill  PrintTemplateType = "waybill"
	PrintTypeCustoms  PrintTemplateType = "customs"
	PrintTypeCarrier  PrintTemplateType = "carrier"
	PrintTypeLabel    PrintTemplateType = "label"
	PrintTypeInvoice  PrintTemplateType = "invoice"
)

// PrintTemplate represents a template for printing documents (waybills, customs forms, etc.).
type PrintTemplate struct {
	ID          int64             `json:"id"`
	TenantID    int64             `json:"tenant_id"`
	Name        string            `json:"name"`
	Code        string            `json:"code"`
	Type        PrintTemplateType `json:"type"`
	Content     string            `json:"content"`
	Width       float64           `json:"width_mm"`
	Height      float64           `json:"height_mm"`
	IsDefault   bool              `json:"is_default"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func DefaultPrintTemplates() []PrintTemplate {
	return []PrintTemplate{
		{Name: "标准面单-100x150", Code: "waybill_100x150", Type: PrintTypeWaybill, Content: `<div class="waybill"><h3>{{.OrderNo}}</h3><p>收件人: {{.RecipientName}}</p><p>地址: {{.Address}}</p><p>重量: {{.Weight}}kg</p><p>件数: {{.Pieces}}</p><p>承运商: {{.Carrier}}</p></div>`, Width: 100, Height: 150},
		{Name: "清关申报单", Code: "customs_decl", Type: PrintTypeCustoms, Content: `<div class="customs"><h3>海关申报单</h3><p>申报人: {{.DeclarantName}}</p><p>品名: {{.ProductName}}</p><p>数量: {{.Quantity}}</p><p>单价: ¥{{.UnitPrice}}</p><p>总价: ¥{{.TotalValue}}</p></div>`, Width: 210, Height: 297},
		{Name: "承运商面单-A4", Code: "carrier_a4", Type: PrintTypeCarrier, Content: `<div class="carrier"><h3>{{.CarrierName}}</h3><p>运单号: {{.CarrierTrackingNo}}</p><p>目的地: {{.Destination}}</p><p>件数: {{.Pieces}}件</p><p>总重: {{.TotalWeight}}kg</p></div>`, Width: 210, Height: 297},
		{Name: "货架标签-50x30", Code: "shelf_label", Type: PrintTypeLabel, Content: `<div class="label"><p>{{.LocationCode}}</p><p>运单: {{.TrackingNumber}}</p></div>`, Width: 50, Height: 30},
		{Name: "客户发票", Code: "client_invoice", Type: PrintTypeInvoice, Content: `<div class="invoice"><h3>发票</h3><p>客户: {{.ClientName}}</p><p>周期: {{.Period}}</p><p>金额: ¥{{.Amount}}</p></div>`, Width: 210, Height: 297},
	}
}
