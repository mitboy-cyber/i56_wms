package service

import (
	"math"
)

// PricingEngine calculates freight charges using the 4D pricing matrix
type PricingEngine struct{}

// PricingInput holds all parameters needed for pricing
type PricingInput struct {
	WeightKG        float64 // actual weight in KG
	LengthCM        float64
	WidthCM         float64
	HeightCM        float64
	TransportType   string  // air, sea, sea_express, air_special
	CargoType       string  // general, furniture, class1~class6, fragile
	TaxType         string  // full_inclusive, frequent_inclusive, non_inclusive
}

// PricingRule defines a single pricing row
type PricingRule struct {
	WeightUnitPrice      float64
	VolumeUnitPrice      float64
	MinCharge            float64
	FirstWeight          float64
	FirstWeightPrice     float64
	AdditionalWeightPrice float64
	FirstVolume          float64
	FirstVolumePrice     float64
	AdditionalVolumePrice float64
}

// PricingResult holds the calculated charges
type PricingResult struct {
	ChargeableWeight float64
	ChargeableVolume float64
	WeightCharge     float64
	VolumeCharge     float64
	TotalCharge      float64
	IsVolumeCharge   bool // true if volume charge > weight charge
}

// Calculate computes freight charges based on input and matching rule
func (pe *PricingEngine) Calculate(input PricingInput, rule PricingRule) PricingResult {
	r := PricingResult{}
	
	// Step 1: Calculate chargeable weight (actual vs dimensional)
	r.ChargeableWeight = pe.calcChargeableWeight(input.WeightKG, input.LengthCM, input.WidthCM, input.HeightCM, input.TransportType)
	r.ChargeableVolume = pe.calcVolume(input.LengthCM, input.WidthCM, input.HeightCM)
	
	// Step 2: Calculate weight-based charge
	r.WeightCharge = pe.calcTieredCharge(r.ChargeableWeight, rule.FirstWeight, rule.FirstWeightPrice, rule.AdditionalWeightPrice)
	
	// Step 3: Calculate volume-based charge (for sea freight)
	if rule.FirstVolumePrice > 0 || rule.AdditionalVolumePrice > 0 {
		r.VolumeCharge = pe.calcTieredCharge(r.ChargeableVolume, rule.FirstVolume, rule.FirstVolumePrice, rule.AdditionalVolumePrice)
	}
	
	// Step 4: Take the greater of weight vs volume charge
	if r.VolumeCharge > r.WeightCharge {
		r.TotalCharge = r.VolumeCharge
		r.IsVolumeCharge = true
	} else {
		r.TotalCharge = r.WeightCharge
	}
	
	// Step 5: Apply minimum charge
	if r.TotalCharge < rule.MinCharge {
		r.TotalCharge = rule.MinCharge
	}
	
	// Step 6: Apply cargo type surcharge for special goods
	r.TotalCharge *= pe.cargoSurchargeMultiplier(input.CargoType, input.TransportType)
	
	r.TotalCharge = math.Round(r.TotalCharge*100) / 100
	return r
}

// calcChargeableWeight returns max(actual_weight, dimensional_weight)
func (pe *PricingEngine) calcChargeableWeight(actualKG, l, w, h float64, transportType string) float64 {
	dimWeight := pe.calcDimWeight(l, w, h, transportType)
	if dimWeight > actualKG {
		return dimWeight
	}
	return actualKG
}

// calcDimWeight: L*W*H / divisor (air=6000, sea=5000, sea_express=6000)
func (pe *PricingEngine) calcDimWeight(l, w, h float64, transportType string) float64 {
	divisor := 6000.0
	switch transportType {
	case "sea":
		divisor = 5000.0
	case "sea_express":
		divisor = 6000.0
	case "air", "air_special":
		divisor = 6000.0
	}
	return l * w * h / divisor
}

// calcVolume: L*W*H in cubic cm / 27872 (1才 = 27872 cm³)
func (pe *PricingEngine) calcVolume(l, w, h float64) float64 {
	return l * w * h / 27872.0
}

// calcTieredCharge: first_weight → first_price, then add_price per extra unit
func (pe *PricingEngine) calcTieredCharge(value, first, firstPrice, addPrice float64) float64 {
	if value <= first {
		return firstPrice
	}
	return firstPrice + math.Ceil(value-first)*addPrice
}

// cargoSurchargeMultiplier returns multiplier for special cargo types
func (pe *PricingEngine) cargoSurchargeMultiplier(cargoType, transportType string) float64 {
	if transportType == "air" || transportType == "air_special" {
		switch cargoType {
		case "class5", "class6":
			return 1.5 // hazardous / large machinery
		case "fragile":
			return 1.2
		}
	}
	// No surcharge for general/sea freight
	return 1.0
}

// NewPricingEngine creates a new pricing engine
func NewPricingEngine() *PricingEngine {
	return &PricingEngine{}
}
