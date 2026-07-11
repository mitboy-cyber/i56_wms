package repository

import "sync"

// ClientSurchargeDisplay represents carrier surcharge data for client portal display.
// Fields match the client_carrier_surcharges.html template.
type ClientSurchargeDisplay struct {
	Carrier       string `json:"carrier"`
	CustomsPoint  string `json:"customs_point"`
	SurchargeName string `json:"surcharge_name"` // 超长费/超材费/偏远费/上楼费
	Condition     string `json:"condition"`       // 单边>150cm/体积重>实重2倍
	Rule          string `json:"rule"`            // 每件加收/按体积重计费差额
	Price         string `json:"price"`           // 100/50/20
}

// MemSurchargeRepo is an in-memory seed repo for carrier surcharge data.
type MemSurchargeRepo struct {
	mu         sync.RWMutex
	surcharges []ClientSurchargeDisplay
}

func NewMemSurchargeRepo() *MemSurchargeRepo {
	return &MemSurchargeRepo{
		surcharges: []ClientSurchargeDisplay{
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "台北",
				SurchargeName: "超長費",
				Condition:     "單邊>150cm",
				Rule:          "每件加收",
				Price:         "100",
			},
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "台北",
				SurchargeName: "超材費",
				Condition:     "體積重>實重2倍",
				Rule:          "按體積重計費差額",
				Price:         "按差",
			},
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "台中",
				SurchargeName: "偏遠費",
				Condition:     "偏遠/離島地區",
				Rule:          "每票加收",
				Price:         "50",
			},
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "台北",
				SurchargeName: "上樓費",
				Condition:     "無電梯4樓以上",
				Rule:          "每層加收",
				Price:         "20",
			},
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "高雄",
				SurchargeName: "偏遠費",
				Condition:     "偏遠/離島地區",
				Rule:          "每票加收",
				Price:         "60",
			},
		},
	}
}

// List returns all surcharges.
func (r *MemSurchargeRepo) List() []ClientSurchargeDisplay {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]ClientSurchargeDisplay, len(r.surcharges))
	copy(result, r.surcharges)
	return result
}
