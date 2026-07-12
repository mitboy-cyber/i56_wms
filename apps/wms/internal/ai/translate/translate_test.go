package translate

import (
	"testing"
)

func TestTranslatorGlossary(t *testing.T) {
	tr := New(nil) // No AI gateway — glossary only

	tests := []struct {
		text string
		from string
		to   string
		want string
	}{
		{"手机壳", "zh", "en", "Phone Case"},
		{"运动鞋", "zh", "en", "Sports Shoes"},
		{"T恤", "zh", "en", "T-Shirt"},
		{"化妆品套装", "zh", "en", "Cosmetics Set"},
		{"充电宝", "zh", "en", "Power Bank"},
		{"蓝牙耳机", "zh", "en", "Bluetooth Earphones"},

		{"手机壳", "zh", "tw", "手機殼"},
		{"运动鞋", "zh", "tw", "運動鞋"},
		{"T恤", "zh", "tw", "T恤"},
		{"化妆品套装", "zh", "tw", "化妝品套裝"},
		{"充电宝", "zh", "tw", "行動電源"},
		{"蓝牙耳机", "zh", "tw", "藍牙耳機"},
	}

	for _, tt := range tests {
		t.Run(tt.text+"→"+tt.to, func(t *testing.T) {
			got := tr.Translate(tt.text, tt.from, tt.to)
			if got != tt.want {
				t.Errorf("Translate(%q, %s, %s): got %q, want %q", tt.text, tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestTranslateProductName(t *testing.T) {
	tr := New(nil)

	en, tw := tr.TranslateProductName("手机壳")
	if en != "Phone Case" {
		t.Errorf("en: got %q, want %q", en, "Phone Case")
	}
	if tw != "手機殼" {
		t.Errorf("tw: got %q, want %q", tw, "手機殼")
	}
}

func TestTranslatorUnknown(t *testing.T) {
	tr := New(nil)
	got := tr.Translate("神秘物品", "zh", "en")
	// Without AI gateway, should return original text
	if got != "神秘物品" {
		t.Errorf("unknown item should return original, got %q", got)
	}
}

func TestTranslatorSameLang(t *testing.T) {
	tr := New(nil)
	got := tr.Translate("手机壳", "zh", "zh")
	if got != "手机壳" {
		t.Errorf("same language should return original, got %q", got)
	}
}

func TestTranslatorEmpty(t *testing.T) {
	tr := New(nil)
	got := tr.Translate("", "zh", "en")
	if got != "" {
		t.Errorf("empty text should return empty, got %q", got)
	}
}

func TestAddGlossaryEntry(t *testing.T) {
	tr := New(nil)
	tr.AddGlossaryEntry("zh", "测试物品", "en", "Test Item")
	got := tr.Translate("测试物品", "zh", "en")
	if got != "Test Item" {
		t.Errorf("got %q, want %q", got, "Test Item")
	}
}
