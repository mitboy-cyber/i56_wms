package domain

import "time"

type OrderStatus string

const (
	StatusPendingPicking   OrderStatus = "pending_picking"
	StatusPicking          OrderStatus = "picking"
	StatusPendingPacking   OrderStatus = "pending_packing"
	StatusPendingLoading   OrderStatus = "pending_loading"
	StatusLoaded           OrderStatus = "loaded"
	StatusInTransit        OrderStatus = "in_transit"
	StatusCustomsClearance OrderStatus = "customs_clearance"
	StatusOutForDelivery   OrderStatus = "out_for_delivery"
	StatusCompleted        OrderStatus = "completed"
	StatusCancelled        OrderStatus = "cancelled"
	StatusShipped          OrderStatus = "shipped" // alias for in_transit in TMS context
)

type Order struct {
	ID                   int64       `json:"id"`
	TenantID             int64       `json:"tenant_id"`
	OrderNo              string      `json:"order_no"`
	WarehouseID          int64       `json:"warehouse_id"`
	ClientID             int64       `json:"client_id"`
	MemberID             int64       `json:"member_id"`
	RouteID              int64       `json:"route_id"`
	RecipientName        string      `json:"recipient_name"`
	TrackingNumbers      string      `json:"tracking_numbers"`
	ParcelCount          int         `json:"parcel_count"`
	Status               OrderStatus `json:"status"`
	ContainerLoadingID   *int64      `json:"container_loading_id"`
	TotalActualWeight    float64     `json:"total_actual_weight"`
	TotalChargeableWeight float64    `json:"total_chargeable_weight"`
	TotalPrice           float64     `json:"total_price"`
	CustomsNumber        string      `json:"customs_number"`
	CarrierTrackingNo    string      `json:"carrier_tracking_no"`
	Remark               string      `json:"remark"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

func ValidTransitions() map[OrderStatus][]OrderStatus {
	return map[OrderStatus][]OrderStatus{
		StatusPendingPicking:   {StatusPicking, StatusCancelled},
		StatusPicking:          {StatusPendingPacking},
		StatusPendingPacking:   {StatusPendingLoading},
		StatusPendingLoading:   {StatusLoaded},
		StatusLoaded:           {StatusInTransit},
		StatusInTransit:        {StatusCustomsClearance},
		StatusCustomsClearance: {StatusOutForDelivery},
		StatusOutForDelivery:   {StatusCompleted},
	}
}

func (o *Order) CanTransitionTo(target OrderStatus) bool {
	allowed, ok := ValidTransitions()[o.Status]
	if !ok { return false }
	for _, s := range allowed {
		if s == target { return true }
	}
	return false
}

func (o *Order) IsCancellable() bool {
	return o.Status == StatusPendingPicking
}
