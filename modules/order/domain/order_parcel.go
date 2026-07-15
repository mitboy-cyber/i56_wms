package domain

import "time"

// OrderParcelStatus represents the status of a parcel within an order.
type OrderParcelStatus string

const (
	OPStatusPending   OrderParcelStatus = "pending"
	OPStatusPicked    OrderParcelStatus = "picked"
	OPStatusPacked    OrderParcelStatus = "packed"
	OPStatusLoaded    OrderParcelStatus = "loaded"
	OPStatusShipped   OrderParcelStatus = "shipped"
	OPStatusDelivered OrderParcelStatus = "delivered"
)

// OrderParcel represents the many-to-many relationship between orders and parcels.
// This is the join entity that tracks which parcels are in which orders and their status.
type OrderParcel struct {
	ID        int64            `json:"id"`
	OrderID   int64            `json:"order_id"`
	ParcelID  int64            `json:"parcel_id"`
	Status    OrderParcelStatus `json:"status"`
	SortOrder int              `json:"sort_order"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func OrderParcelValidTransitions() map[OrderParcelStatus][]OrderParcelStatus {
	return map[OrderParcelStatus][]OrderParcelStatus{
		OPStatusPending:   {OPStatusPicked},
		OPStatusPicked:    {OPStatusPacked},
		OPStatusPacked:    {OPStatusLoaded},
		OPStatusLoaded:    {OPStatusShipped},
		OPStatusShipped:   {OPStatusDelivered},
		OPStatusDelivered: {},
	}
}

func (op *OrderParcel) CanTransitionTo(target OrderParcelStatus) bool {
	allowed, ok := OrderParcelValidTransitions()[op.Status]
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
