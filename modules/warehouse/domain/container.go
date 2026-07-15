package domain

import "time"

// ContainerType represents the type of shipping container.
type ContainerType string

const (
	Container20GP ContainerType = "20GP"
	Container40GP ContainerType = "40GP"
	Container40HQ ContainerType = "40HQ"
	Container45HQ ContainerType = "45HQ"
)

// ContainerStatus represents the status of a container in the warehouse.
type ContainerStatus string

const (
	ContainerStatusAvailable ContainerStatus = "available"
	ContainerStatusLoading   ContainerStatus = "loading"
	ContainerStatusLoaded    ContainerStatus = "loaded"
	ContainerStatusSealed    ContainerStatus = "sealed"
	ContainerStatusShipped   ContainerStatus = "shipped"
)

// Container represents a shipping container in the warehouse.
type Container struct {
	ID             int64           `json:"id"`
	WarehouseID    int64           `json:"warehouse_id"`
	ContainerNo    string          `json:"container_no"`
	ContainerType  ContainerType   `json:"container_type"`
	SealNo         string          `json:"seal_no"`
	RouteID        int64           `json:"route_id"`
	Status         ContainerStatus `json:"status"`
	MaxCapacity    float64         `json:"max_capacity"`
	CurrentWeight  float64         `json:"current_weight"`
	ParcelCount    int             `json:"parcel_count"`
	LoadedBy       string          `json:"loaded_by"`
	LoadedAt       *time.Time      `json:"loaded_at"`
	SealedAt       *time.Time      `json:"sealed_at"`
	ShippedAt      *time.Time      `json:"shipped_at"`
	Remark         string          `json:"remark"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func ContainerValidTransitions() map[ContainerStatus][]ContainerStatus {
	return map[ContainerStatus][]ContainerStatus{
		ContainerStatusAvailable: {ContainerStatusLoading},
		ContainerStatusLoading:   {ContainerStatusLoaded},
		ContainerStatusLoaded:    {ContainerStatusSealed},
		ContainerStatusSealed:    {ContainerStatusShipped},
		ContainerStatusShipped:   {},
	}
}

func (c *Container) CanTransitionTo(target ContainerStatus) bool {
	allowed, ok := ContainerValidTransitions()[c.Status]
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
