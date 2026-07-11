package service

import "testing"

func TestPricingEngine(t *testing.T) {
	pe := NewPricingEngine()
	
	// Sea freight: furniture, full_inclusive, 8.5kg
	seaRule := PricingRule{
		FirstWeight:10, FirstWeightPrice:25, AdditionalWeightPrice:2.50,
		FirstVolume:1, FirstVolumePrice:15, AdditionalVolumePrice:15,
		MinCharge:50,
	}
	
	// Test 1: Sea freight, 8.5kg (below first weight) → ¥25
	r := pe.Calculate(PricingInput{WeightKG:8.5, LengthCM:80, WidthCM:40, HeightCM:30, TransportType:"sea", CargoType:"furniture", TaxType:"full_inclusive"}, seaRule)
	if r.TotalCharge != 60.0 { t.Errorf("Test1: expected ¥50(min), got ¥%.2f", r.TotalCharge) }
	t.Logf("Test1 (8.5kg家具海运): weight=¥%.2f, vol=¥%.2f, total=¥%.2f, isVol=%v", r.WeightCharge, r.VolumeCharge, r.TotalCharge, r.IsVolumeCharge)
	
	// Test 2: Sea freight, 15kg → first price + 5 × add
	r2 := pe.Calculate(PricingInput{WeightKG:15, LengthCM:80, WidthCM:40, HeightCM:30, TransportType:"sea", CargoType:"class1", TaxType:"full_inclusive"}, seaRule)
	t.Logf("Test2 (15kg一类海运): weight=¥%.2f, vol=¥%.2f, total=¥%.2f", r2.WeightCharge, r2.VolumeCharge, r2.TotalCharge)
	
	// Test 3: Air freight with hazardous cargo
	airRule := PricingRule{FirstWeight:0.5, FirstWeightPrice:25, AdditionalWeightPrice:20, MinCharge:50}
	r3 := pe.Calculate(PricingInput{WeightKG:3, LengthCM:30, WidthCM:20, HeightCM:12, TransportType:"air", CargoType:"class5", TaxType:"full_inclusive"}, airRule)
	t.Logf("Test3 (3kg五类空运): dimW=%.2f, chargeW=%.2f, total=¥%.2f", pe.calcDimWeight(30,20,12,"air"), r3.ChargeableWeight, r3.TotalCharge)
	
	// Test 4: Large volume, light weight (volume charge dominates)
	r4 := pe.Calculate(PricingInput{WeightKG:2, LengthCM:200, WidthCM:100, HeightCM:50, TransportType:"sea", CargoType:"furniture", TaxType:"full_inclusive"}, seaRule)
	t.Logf("Test4 (2kg超大件海运): weight=¥%.2f, vol=¥%.2f, total=¥%.2f, isVol=%v", r4.WeightCharge, r4.VolumeCharge, r4.TotalCharge, r4.IsVolumeCharge)
}
