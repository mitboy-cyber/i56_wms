// Package translate provides multi-language translation for WMS product names
// and content. Uses a glossary for common terms with AI fallback for unknown items.
package translate

import (
	"context"
	"strings"
	"sync"

	"github.com/i56/framework/ai/gateway"
)

// Translator handles CN↔EN↔TW product name translations.
type Translator struct {
	gateway  gateway.Gateway
	glossary map[string]map[string]string // source_lang:source_text -> target_lang:translation
	mu       sync.RWMutex
}

// New creates a new Translator with the given AI gateway.
func New(gw gateway.Gateway) *Translator {
	t := &Translator{
		gateway:  gw,
		glossary: make(map[string]map[string]string),
	}
	t.seedGlossary()
	return t
}

// seedGlossary populates common product name translations.
func (t *Translator) seedGlossary() {
	// CN→EN
	t.AddGlossaryEntry("zh", "手机壳", "en", "Phone Case")
	t.AddGlossaryEntry("zh", "运动鞋", "en", "Sports Shoes")
	t.AddGlossaryEntry("zh", "T恤", "en", "T-Shirt")
	t.AddGlossaryEntry("zh", "化妆品套装", "en", "Cosmetics Set")
	t.AddGlossaryEntry("zh", "充电宝", "en", "Power Bank")
	t.AddGlossaryEntry("zh", "蓝牙耳机", "en", "Bluetooth Earphones")
	t.AddGlossaryEntry("zh", "数据线", "en", "USB Cable")
	t.AddGlossaryEntry("zh", "贴膜", "en", "Screen Protector")
	t.AddGlossaryEntry("zh", "衣服", "en", "Clothing")
	t.AddGlossaryEntry("zh", "裤子", "en", "Pants")
	t.AddGlossaryEntry("zh", "裙子", "en", "Dress")
	t.AddGlossaryEntry("zh", "衬衫", "en", "Shirt")
	t.AddGlossaryEntry("zh", "外套", "en", "Jacket")
	t.AddGlossaryEntry("zh", "鞋子", "en", "Shoes")
	t.AddGlossaryEntry("zh", "包包", "en", "Bag")
	t.AddGlossaryEntry("zh", "手表", "en", "Watch")
	t.AddGlossaryEntry("zh", "眼镜", "en", "Glasses")
	t.AddGlossaryEntry("zh", "化妆品", "en", "Cosmetics")
	t.AddGlossaryEntry("zh", "口红", "en", "Lipstick")
	t.AddGlossaryEntry("zh", "面膜", "en", "Facial Mask")
	t.AddGlossaryEntry("zh", "保健品", "en", "Health Supplement")
	t.AddGlossaryEntry("zh", "食品", "en", "Food")
	t.AddGlossaryEntry("zh", "零食", "en", "Snacks")
	t.AddGlossaryEntry("zh", "茶叶", "en", "Tea")
	t.AddGlossaryEntry("zh", "香水", "en", "Perfume")
	t.AddGlossaryEntry("zh", "洗发水", "en", "Shampoo")
	t.AddGlossaryEntry("zh", "玩具", "en", "Toy")
	t.AddGlossaryEntry("zh", "书籍", "en", "Book")
	t.AddGlossaryEntry("zh", "充电器", "en", "Charger")
	t.AddGlossaryEntry("zh", "电脑", "en", "Computer")
	t.AddGlossaryEntry("zh", "平板", "en", "Tablet")
	t.AddGlossaryEntry("zh", "相机", "en", "Camera")
	t.AddGlossaryEntry("zh", "音箱", "en", "Speaker")

	// CN→TW
	t.AddGlossaryEntry("zh", "手机壳", "tw", "手機殼")
	t.AddGlossaryEntry("zh", "运动鞋", "tw", "運動鞋")
	t.AddGlossaryEntry("zh", "T恤", "tw", "T恤")
	t.AddGlossaryEntry("zh", "化妆品套装", "tw", "化妝品套裝")
	t.AddGlossaryEntry("zh", "充电宝", "tw", "行動電源")
	t.AddGlossaryEntry("zh", "蓝牙耳机", "tw", "藍牙耳機")
	t.AddGlossaryEntry("zh", "数据线", "tw", "傳輸線")
	t.AddGlossaryEntry("zh", "贴膜", "tw", "螢幕保護貼")
	t.AddGlossaryEntry("zh", "衣服", "tw", "衣服")
	t.AddGlossaryEntry("zh", "裤子", "tw", "褲子")
	t.AddGlossaryEntry("zh", "裙子", "tw", "裙子")
	t.AddGlossaryEntry("zh", "衬衫", "tw", "襯衫")
	t.AddGlossaryEntry("zh", "外套", "tw", "外套")
	t.AddGlossaryEntry("zh", "鞋子", "tw", "鞋子")
	t.AddGlossaryEntry("zh", "包包", "tw", "包包")
	t.AddGlossaryEntry("zh", "手表", "tw", "手錶")
	t.AddGlossaryEntry("zh", "化妆品", "tw", "化妝品")
	t.AddGlossaryEntry("zh", "口红", "tw", "口紅")
	t.AddGlossaryEntry("zh", "面膜", "tw", "面膜")
	t.AddGlossaryEntry("zh", "保健品", "tw", "保健食品")
	t.AddGlossaryEntry("zh", "食品", "tw", "食品")
	t.AddGlossaryEntry("zh", "零食", "tw", "零食")
	t.AddGlossaryEntry("zh", "茶叶", "tw", "茶葉")
	t.AddGlossaryEntry("zh", "香水", "tw", "香水")
	t.AddGlossaryEntry("zh", "洗发水", "tw", "洗髮精")
	t.AddGlossaryEntry("zh", "玩具", "tw", "玩具")
	t.AddGlossaryEntry("zh", "书籍", "tw", "書籍")
	t.AddGlossaryEntry("zh", "电脑", "tw", "電腦")
	t.AddGlossaryEntry("zh", "手机", "tw", "手機")
	t.AddGlossaryEntry("zh", "相机", "tw", "相機")
}

// AddGlossaryEntry adds a single translation entry to the glossary.
func (t *Translator) AddGlossaryEntry(fromLang, fromText, toLang, toText string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	key := fromLang + ":" + fromText
	if t.glossary[key] == nil {
		t.glossary[key] = make(map[string]string)
	}
	t.glossary[key][toLang] = toText
}

// Translate translates text from one language to another.
// Languages: "zh" (Simplified Chinese), "en" (English), "tw" (Traditional Chinese).
func (t *Translator) Translate(text, from, to string) string {
	text = strings.TrimSpace(text)
	if text == "" || from == to {
		return text
	}

	// Check glossary first
	t.mu.RLock()
	key := from + ":" + text
	if entries, ok := t.glossary[key]; ok {
		if translated, ok2 := entries[to]; ok2 {
			t.mu.RUnlock()
			return translated
		}
	}
	t.mu.RUnlock()

	// Try AI fallback
	if t.gateway != nil {
		if result := t.aiTranslate(text, from, to); result != "" {
			t.AddGlossaryEntry(from, text, to, result)
			return result
		}
	}

	return text
}

// TranslateProductName translates a Chinese product name to English and Traditional Chinese.
func (t *Translator) TranslateProductName(cnName string) (en, tw string) {
	en = t.Translate(cnName, "zh", "en")
	tw = t.Translate(cnName, "zh", "tw")
	return
}

// aiTranslate uses AI to translate text.
func (t *Translator) aiTranslate(text, from, to string) string {
	langNames := map[string]string{"zh": "Simplified Chinese", "en": "English", "tw": "Traditional Chinese (Taiwan)"}
	fromName := langNames[from]
	toName := langNames[to]
	if fromName == "" {
		fromName = from
	}
	if toName == "" {
		toName = to
	}

	resp, err := t.gateway.Chat(context.Background(), &gateway.ChatRequest{
		Messages: []gateway.Message{
			{Role: gateway.RoleSystem, Content: "You are a professional translator. Translate the given text from " + fromName + " to " + toName + ". Return ONLY the translated text, no explanations or quotes."},
			{Role: gateway.RoleUser, Content: text},
		},
		Temperature: 0.1,
		MaxTokens:   100,
	})
	if err != nil || resp == nil {
		return ""
	}

	translated := strings.TrimSpace(resp.Content)
	// Remove quotes if AI wrapped the result
	translated = strings.Trim(translated, "\"'")
	return translated
}
