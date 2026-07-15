package domain

import "time"

// ConsolidationOrderStatus represents the status of a consolidation (集运) order.
type ConsolidationOrderStatus string

const (
	ConsolStatusDraft         ConsolidationOrderStatus = "draft"
	ConsolStatusPendingMerge  ConsolidationOrderStatus = "pending_merge"
	ConsolStatusMerged        ConsolidationOrderStatus = "merged"
	ConsolStatusWeighing      ConsolidationOrderStatus = "weighing"
	ConsolStatusPacking       ConsolidationOrderStatus = "packing"
	ConsolStatusPacked        ConsolidationOrderStatus = "packed"
	ConsolStatusOutbound      ConsolidationOrderStatus = "outbound"
	ConsolStatusShipped       ConsolidationOrderStatus = "shipped"
	ConsolStatusCompleted     ConsolidationOrderStatus = "completed"
	ConsolStatusCancelled     ConsolidationOrderStatus = "cancelled"
)

// ConsolidationOrder represents a consolidation (集运) order that groups multiple parcels
// for a single client member to be shipped together.
type ConsolidationOrder struct {
	ID                int64                     `json:"id"`
	TenantID          int64                     `json:"tenant_id"`
	OrderNo           string                    `json:"order_no"`
	ClientID          int64                     `json:"client_id"`
	MemberID          int64                     `json:"member_id"`
	MemberAddressID   int64                     `json:"member_address_id"`
	WarehouseID       int64                     `json:"warehouse_id"`
	RouteID           int64                     `json:"route_id"`
	ParcelIDs         []int64                   `json:"parcel_ids"`
	ParcelCount       int                       `json:"parcel_count"`
	TotalWeight       float64                   `json:"total_weight"`
	TotalChargeable   float64                   `json:"total_chargeable_weight"`
	ShippingFee       float64                   `json:"shipping_fee"`
	ServiceFee        float64                   `json:"service_fee"`
	TotalPrice        float64                   `json:"total_price"`
	Status            ConsolidationOrderStatus  `json:"status"`
	DeclarantID       int64                     `json:"declarant_id"`
	CustomsNumber     string                    `json:"customs_number"`
	CarrierTrackingNo string                    `json:"carrier_tracking_no"`
	Remark            string                    `json:"remark"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
}

func ConsolidationValidTransitions() map[ConsolidationOrderStatus][]ConsolidationOrderStatus {
	return map[ConsolidationOrderStatus][]ConsolidationOrderStatus{
		ConsolStatusDraft:        {ConsolStatusPendingMerge, ConsolStatusCancelled},
		ConsolStatusPendingMerge: {ConsolStatusMerged},
		ConsolStatusMerged:       {ConsolStatusWeighing},
		ConsolStatusWeighing:     {ConsolStatusPacking},
		ConsolStatusPacking:      {ConsolStatusPacked},
		ConsolStatusPacked:       {ConsolStatusOutbound},
		ConsolStatusOutbound:     {ConsolStatusShipped},
		ConsolStatusShipped:      {ConsolStatusCompleted},
	}
}

func (o *ConsolidationOrder) CanTransitionTo(target ConsolidationOrderStatus) bool {
	allowed, ok := ConsolidationValidTransitions()[o.Status]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == target {
			return true
		}
	}
	return false
}
