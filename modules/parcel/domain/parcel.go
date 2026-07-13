package domain

import "time"

type ParcelStatus string

const (
	StatusPreDeclared   ParcelStatus = "pre_declared"
	StatusReceived      ParcelStatus = "received"
	StatusWeighed       ParcelStatus = "weighed"
	StatusStored        ParcelStatus = "stored"
	StatusPicked        ParcelStatus = "picked"
	StatusPacked        ParcelStatus = "packed"
	StatusOutbound      ParcelStatus = "outbound"
	StatusContainerArea ParcelStatus = "container_area"
	StatusLoaded        ParcelStatus = "loaded"
	StatusShipped       ParcelStatus = "shipped"
	StatusCustoms       ParcelStatus = "customs"
	StatusDelivering    ParcelStatus = "delivering"
	StatusDelivered     ParcelStatus = "delivered"
	StatusAbnormal      ParcelStatus = "abnormal"
	StatusReturned      ParcelStatus = "returned"
)

type Parcel struct {
	ID              int64        `json:"id"`
	TenantID        int64        `json:"tenant_id"`
	WarehouseID     int64        `json:"warehouse_id"`
	ClientID        int64        `json:"client_id"`
	TrackingNumber  string       `json:"tracking_number"`
	CourierCode     string       `json:"courier_code"`
	CargoType       string       `json:"cargo_type"`
	ProductName     string       `json:"product_name"`
	ParcelName      string       `json:"parcel_name"`
	ActualWeight    float64      `json:"actual_weight"`
	Length          float64      `json:"length"`
	Width           float64      `json:"width"`
	Height          float64      `json:"height"`
	Status          ParcelStatus `json:"status"`
	IsAbnormal      bool         `json:"is_abnormal"`
	LocationCode    string       `json:"location_code"`
	ImageURLs       []string     `json:"image_urls"`
	OrderID         *int64       `json:"order_id"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func ValidTransitions() map[ParcelStatus][]ParcelStatus {
	return map[ParcelStatus][]ParcelStatus{
		StatusPreDeclared:   {StatusReceived, StatusReturned},
		StatusReceived:      {StatusWeighed, StatusReturned},
		StatusWeighed:       {StatusStored, StatusReturned},
		StatusStored:        {StatusPicked, StatusContainerArea, StatusReturned, StatusAbnormal},
		StatusPicked:        {StatusPacked, StatusReturned},
		StatusPacked:        {StatusShipped, StatusLoaded, StatusContainerArea},
		StatusContainerArea: {StatusLoaded, StatusShipped},
		StatusLoaded:        {StatusShipped, StatusCustoms},
		StatusCustoms:       {StatusDelivering, StatusDelivered, StatusReturned, StatusAbnormal},
		StatusDelivering:    {StatusDelivered, StatusAbnormal},
		StatusShipped:       {StatusDelivering, StatusDelivered, StatusCustoms},
		StatusAbnormal:      {StatusStored, StatusReturned},
	}
}

func (p *Parcel) CanTransitionTo(target ParcelStatus) bool {
	allowed, ok := ValidTransitions()[p.Status]
	if !ok { return false }
	for _, s := range allowed {
		if s == target { return true }
	}
	return false
}

func (p *Parcel) VolumetricWeight() float64 {
	if p.Length <= 0 || p.Width <= 0 || p.Height <= 0 {
		return 0
	}
	return (p.Length * p.Width * p.Height) / 6000.0
}

func (p *Parcel) ChargeableWeight() float64 {
	vw := p.VolumetricWeight()
	if vw > p.ActualWeight { return vw }
	return p.ActualWeight
}
