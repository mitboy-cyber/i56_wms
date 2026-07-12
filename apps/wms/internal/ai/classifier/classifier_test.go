package classifier

import (
	"testing"
)

func TestClassifierRuleEngine(t *testing.T) {
	c := New(nil) // No AI gateway — rules only

	tests := []struct {
		desc          string
		wantCategory  string
		wantHSCode    string
		wantRiskLevel string
		wantInspect   bool
	}{
		{"手机壳", "general", "3926.90.9090", "low", false},
		{"蓝牙耳机", "electronic", "8518.30.0000", "low", false},
		{"充电宝", "dangerous", "8507.60.0090", "high", true},
		{"运动鞋", "general", "6404.11.0000", "low", false},
		{"T恤", "general", "6109.10.0000", "low", false},
		{"化妆品套装", "sensitive", "3304.99.0090", "medium", false},
		{"香水", "dangerous", "3303.00.0000", "high", true},
		{"洗发水", "liquid", "3305.10.0000", "medium", false},
		{"玻璃杯", "fragile", "7013.37.0000", "high", true},
		{"数据线", "general", "8544.42.0000", "low", false},
		{"电脑", "electronic", "8471.30.0000", "medium", true},
		{"面膜", "sensitive", "3304.99.0090", "medium", false},
		{"保健品", "sensitive", "2106.90.9090", "medium", false},
		{"茶叶", "sensitive", "0902.20.0000", "low", false},
		{"玩具", "general", "9503.00.0000", "low", false},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := c.Classify(tt.desc)
			if result.Category != tt.wantCategory {
				t.Errorf("category: got %s, want %s", result.Category, tt.wantCategory)
			}
			if result.HSCode != tt.wantHSCode {
				t.Errorf("hs_code: got %s, want %s", result.HSCode, tt.wantHSCode)
			}
			if result.RiskLevel != tt.wantRiskLevel {
				t.Errorf("risk_level: got %s, want %s", result.RiskLevel, tt.wantRiskLevel)
			}
			if result.NeedInspect != tt.wantInspect {
				t.Errorf("need_inspect: got %v, want %v", result.NeedInspect, tt.wantInspect)
			}
			if tt.desc != "unknown item xyz" && result.Confidence < 0.9 {
				t.Errorf("confidence too low for known item: %.2f", result.Confidence)
			}
		})
	}
}

func TestClassifierBatch(t *testing.T) {
	c := New(nil)
	descs := []string{"手机壳", "运动鞋", "蓝牙耳机", "充电宝"}
	results := c.ClassifyBatch(descs)
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	for i, r := range results {
		if r.Category == "" {
			t.Errorf("result[%d] (%s): empty category", i, descs[i])
		}
	}
}

func TestClassifierCache(t *testing.T) {
	c := New(nil)
	desc := "手机壳"
	r1 := c.Classify(desc)
	r2 := c.Classify(desc)
	if r1.Category != r2.Category || r1.HSCode != r2.HSCode {
		t.Error("cached result should match original")
	}
}

func TestCategoryCN(t *testing.T) {
	tests := map[string]string{
		"general":    "普货",
		"electronic": "电子产品",
		"dangerous":  "危险品",
		"liquid":     "液体",
		"fragile":    "易碎品",
		"sensitive":  "特货",
		"powder":     "粉末",
	}
	for cat, want := range tests {
		got := CategoryCN(cat)
		if got != want {
			t.Errorf("CategoryCN(%s): got %s, want %s", cat, got, want)
		}
	}
}

func TestClassifierUnknown(t *testing.T) {
	c := New(nil)
	result := c.Classify("some unknown product xyz")
	// Without AI gateway, should return generic result
	if result.Category != "general" {
		t.Errorf("unknown item should default to general, got %s", result.Category)
	}
	if result.Confidence != 0.3 {
		t.Errorf("unknown item without AI should have confidence 0.3, got %.2f", result.Confidence)
	}
}
