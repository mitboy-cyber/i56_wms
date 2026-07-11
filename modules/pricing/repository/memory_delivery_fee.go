package repository

import "sync"

// ClientDeliveryFeeDisplay represents carrier delivery fee for client portal display.
// Fields match the client_delivery_fees.html template.
type ClientDeliveryFeeDisplay struct {
	Carrier       string `json:"carrier"`
	CustomsPoint  string `json:"customs_point"`  // 台北/台中/高雄
	Area          string `json:"area"`            // 預設/台湾北部/台湾中部/台湾南部/台湾东部
	DeliveryType  string `json:"delivery_type"`   // 宅配/專車/自取
	Condition     string `json:"condition"`       // 重量>39.8/單邊長>=600或重量>=500
	Price         string `json:"price"`
	FreeThreshold string `json:"free_threshold"`  // >=10kg or empty
}

// MemDeliveryFeeRepo is an in-memory seed repo for carrier delivery fees.
type MemDeliveryFeeRepo struct {
	mu   sync.RWMutex
	fees []ClientDeliveryFeeDisplay
}

func NewMemDeliveryFeeRepo() *MemDeliveryFeeRepo {
	return &MemDeliveryFeeRepo{
		fees: []ClientDeliveryFeeDisplay{
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "台北",
				Area:          "預設",
				DeliveryType:  "宅配",
				Condition:     "重量>39.8",
				Price:         "20",
				FreeThreshold: ">=10kg",
			},
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "台北",
				Area:          "預設",
				DeliveryType:  "專車",
				Condition:     "單邊長>=600或重量>=500",
				Price:         "3500",
				FreeThreshold: "",
			},
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "台中",
				Area:          "預設",
				DeliveryType:  "宅配",
				Condition:     "重量>39.8",
				Price:         "20",
				FreeThreshold: ">=10kg",
			},
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "高雄",
				Area:          "預設",
				DeliveryType:  "宅配",
				Condition:     "重量>39.8",
				Price:         "25",
				FreeThreshold: ">=10kg",
			},
			{
				Carrier:       "新竹物流",
				CustomsPoint:  "台北",
				Area:          "東部",
				DeliveryType:  "宅配",
				Condition:     "宜蘭/花蓮/台東",
				Price:         "50",
				FreeThreshold: "",
			},
		},
	}
}

// List returns all delivery fees.
func (r *MemDeliveryFeeRepo) List() []ClientDeliveryFeeDisplay {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]ClientDeliveryFeeDisplay, len(r.fees))
	copy(result, r.fees)
	return result
}
