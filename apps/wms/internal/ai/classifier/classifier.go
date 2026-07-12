// Package classifier provides smart cargo classification using rule engine + AI fallback.
// Common items are classified by rules (instant), uncommon items fall back to AI.
package classifier

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/i56/framework/ai/gateway"
)

// CargoClassifier classifies product descriptions into cargo types, HS codes,
// risk levels, and related metadata.
type CargoClassifier struct {
	gateway gateway.Gateway
	cache   map[string]CargoResult
	mu      sync.RWMutex
}

// CargoResult holds the classification result for a product.
type CargoResult struct {
	Category    string  // general/sensitive/dangerous/fragile/liquid/electronic/powder
	HSCode      string  // HS code
	SubCategory string  // sub-category detail
	RiskLevel   string  // low/medium/high
	Confidence  float64 // 0.0-1.0
	NeedInspect bool    // whether inspection is recommended
}

// RuleEntry defines a keyword-based classification rule.
type RuleEntry struct {
	Keywords    []string
	Category    string
	HSCode      string
	SubCategory string
	RiskLevel   string
	NeedInspect bool
}

// Pre-defined rule engine for common items.
var defaultRules = []RuleEntry{
	// 普货 — general cargo
	{Keywords: []string{"手机壳", "phone case"}, Category: "general", HSCode: "3926.90.9090", SubCategory: "手机配件", RiskLevel: "low"},
	{Keywords: []string{"运动鞋", "sports shoes", "球鞋", "跑鞋"}, Category: "general", HSCode: "6404.11.0000", SubCategory: "鞋类", RiskLevel: "low"},
	{Keywords: []string{"t恤", "t-shirt", "t shirt", "T恤"}, Category: "general", HSCode: "6109.10.0000", SubCategory: "服装", RiskLevel: "low"},
	{Keywords: []string{"衣服", "服装", "clothing", "衬衫", "裤子", "裙子", "外套", "夹克"}, Category: "general", HSCode: "6100.00.0000", SubCategory: "服装", RiskLevel: "low"},
	{Keywords: []string{"数据线", "充电线", "usb线", "cable"}, Category: "general", HSCode: "8544.42.0000", SubCategory: "线材", RiskLevel: "low"},
	{Keywords: []string{"贴膜", "钢化膜", "screen protector"}, Category: "general", HSCode: "3926.90.9090", SubCategory: "手机配件", RiskLevel: "low"},
	{Keywords: []string{"书籍", "书本", "book"}, Category: "general", HSCode: "4901.99.0000", SubCategory: "印刷品", RiskLevel: "low"},
	{Keywords: []string{"玩具", "toy", "公仔", "玩偶", "模型"}, Category: "general", HSCode: "9503.00.0000", SubCategory: "玩具", RiskLevel: "low"},
	{Keywords: []string{"文具", "笔", "笔记本", "stationery"}, Category: "general", HSCode: "4820.10.0000", SubCategory: "文具", RiskLevel: "low"},
	{Keywords: []string{"包包", "手提包", "背包", "bag"}, Category: "general", HSCode: "4202.22.0000", SubCategory: "箱包", RiskLevel: "low"},

	// 电子产品 — electronic
	{Keywords: []string{"蓝牙耳机", "bluetooth earphone", "无线耳机"}, Category: "electronic", HSCode: "8518.30.0000", SubCategory: "音频设备", RiskLevel: "low"},
	{Keywords: []string{"耳机", "earphone", "headphone", "耳塞"}, Category: "electronic", HSCode: "8518.30.0000", SubCategory: "音频设备", RiskLevel: "low"},
	{Keywords: []string{"音箱", "speaker", "音响"}, Category: "electronic", HSCode: "8518.22.0000", SubCategory: "音频设备", RiskLevel: "low"},
	{Keywords: []string{"手机", "phone", "智能手机", "smartphone"}, Category: "electronic", HSCode: "8517.12.0000", SubCategory: "通讯设备", RiskLevel: "medium", NeedInspect: true},
	{Keywords: []string{"平板", "tablet", "ipad"}, Category: "electronic", HSCode: "8471.30.0000", SubCategory: "计算设备", RiskLevel: "medium", NeedInspect: true},
	{Keywords: []string{"电脑", "笔记本", "laptop", "computer"}, Category: "electronic", HSCode: "8471.30.0000", SubCategory: "计算设备", RiskLevel: "medium", NeedInspect: true},
	{Keywords: []string{"手表", "watch", "智能手表", "smartwatch"}, Category: "electronic", HSCode: "9102.12.0000", SubCategory: "穿戴设备", RiskLevel: "low"},
	{Keywords: []string{"充电器", "charger", "适配器", "adapter"}, Category: "electronic", HSCode: "8504.40.0000", SubCategory: "电源配件", RiskLevel: "low"},
	{Keywords: []string{"存储卡", "sd卡", "memory card", "u盘"}, Category: "electronic", HSCode: "8523.51.0000", SubCategory: "存储设备", RiskLevel: "low"},
	{Keywords: []string{"相机", "camera", "摄像机"}, Category: "electronic", HSCode: "8525.80.0000", SubCategory: "影像设备", RiskLevel: "medium", NeedInspect: true},

	// 危险品 — dangerous goods (battery/power bank)
	{Keywords: []string{"充电宝", "power bank", "移动电源", "充电电池"}, Category: "dangerous", HSCode: "8507.60.0090", SubCategory: "锂电池", RiskLevel: "high", NeedInspect: true},
	{Keywords: []string{"电池", "battery", "锂电池", "锂电"}, Category: "dangerous", HSCode: "8507.60.0090", SubCategory: "锂电池", RiskLevel: "high", NeedInspect: true},
	{Keywords: []string{"打火机", "lighter"}, Category: "dangerous", HSCode: "9613.80.0000", SubCategory: "易燃品", RiskLevel: "high", NeedInspect: true},
	{Keywords: []string{"香水", "perfume", "精油", "essential oil"}, Category: "dangerous", HSCode: "3303.00.0000", SubCategory: "易燃液体", RiskLevel: "high", NeedInspect: true},
	{Keywords: []string{"喷雾", "spray", "气雾剂", "aerosol"}, Category: "dangerous", HSCode: "3307.49.0000", SubCategory: "压力容器", RiskLevel: "high", NeedInspect: true},

	// 液体 — liquid
	{Keywords: []string{"洗发水", "shampoo", "沐浴露", "body wash"}, Category: "liquid", HSCode: "3305.10.0000", SubCategory: "洗护用品", RiskLevel: "medium"},
	{Keywords: []string{"饮料", "drink", "果汁", "juice"}, Category: "liquid", HSCode: "2202.10.0000", SubCategory: "饮品", RiskLevel: "medium"},
	{Keywords: []string{"食用油", "cooking oil", "橄榄油", "olive oil"}, Category: "liquid", HSCode: "1515.90.0000", SubCategory: "食用油", RiskLevel: "medium"},
	{Keywords: []string{"酱油", "soy sauce", "醋", "vinegar"}, Category: "liquid", HSCode: "2103.10.0000", SubCategory: "调味品", RiskLevel: "medium"},

	// 易碎品 — fragile
	{Keywords: []string{"玻璃杯", "glass cup", "水杯"}, Category: "fragile", HSCode: "7013.37.0000", SubCategory: "玻璃器皿", RiskLevel: "high", NeedInspect: true},
	{Keywords: []string{"花瓶", "vase", "瓷器", "porcelain", "陶瓷", "ceramic"}, Category: "fragile", HSCode: "6911.10.0000", SubCategory: "陶瓷制品", RiskLevel: "high", NeedInspect: true},
	{Keywords: []string{"镜子", "mirror"}, Category: "fragile", HSCode: "7009.92.0000", SubCategory: "玻璃制品", RiskLevel: "high", NeedInspect: true},

	// 特货 — sensitive/special
	{Keywords: []string{"化妆品", "cosmetics", "makeup", "彩妆", "护肤"}, Category: "sensitive", HSCode: "3304.99.0090", SubCategory: "化妆品", RiskLevel: "medium"},
	{Keywords: []string{"口红", "lipstick", "唇膏"}, Category: "sensitive", HSCode: "3304.10.0000", SubCategory: "化妆品", RiskLevel: "medium"},
	{Keywords: []string{"面膜", "mask", "facial mask"}, Category: "sensitive", HSCode: "3304.99.0090", SubCategory: "化妆品", RiskLevel: "medium"},
	{Keywords: []string{"化妆品套装", "cosmetics set"}, Category: "sensitive", HSCode: "3304.99.0090", SubCategory: "化妆品套装", RiskLevel: "medium"},
	{Keywords: []string{"保健品", "supplement", "维生素", "vitamin"}, Category: "sensitive", HSCode: "2106.90.9090", SubCategory: "保健品", RiskLevel: "medium"},
	{Keywords: []string{"药品", "medicine", "药物", "药丸"}, Category: "sensitive", HSCode: "3004.90.0000", SubCategory: "药品", RiskLevel: "high", NeedInspect: true},
	{Keywords: []string{"食品", "food", "零食", "snack", "饼干", "cookies"}, Category: "sensitive", HSCode: "1905.90.0000", SubCategory: "食品", RiskLevel: "medium"},
	{Keywords: []string{"茶叶", "tea", "茶"}, Category: "sensitive", HSCode: "0902.20.0000", SubCategory: "茶叶", RiskLevel: "low"},

	// 粉末 — powder
	{Keywords: []string{"奶粉", "milk powder", "蛋白粉", "protein powder"}, Category: "powder", HSCode: "0402.21.0000", SubCategory: "粉末食品", RiskLevel: "medium", NeedInspect: true},
	{Keywords: []string{"面粉", "flour", "调味粉", "spice powder"}, Category: "powder", HSCode: "1101.00.0000", SubCategory: "食品粉末", RiskLevel: "medium"},
}

// New creates a new CargoClassifier with the given AI gateway.
func New(gw gateway.Gateway) *CargoClassifier {
	return &CargoClassifier{
		gateway: gw,
		cache:   make(map[string]CargoResult),
	}
}

// Classify classifies a single product description.
// Rules are tried first; AI is used as a fallback for unknown items.
func (c *CargoClassifier) Classify(description string) CargoResult {
	desc := strings.TrimSpace(description)
	if desc == "" {
		return CargoResult{Category: "general", HSCode: "0000.00.0000", RiskLevel: "low", Confidence: 0.0}
	}

	// Check cache
	c.mu.RLock()
	if result, ok := c.cache[strings.ToLower(desc)]; ok {
		c.mu.RUnlock()
		return result
	}
	c.mu.RUnlock()

	// Try rules first
	result, matched := c.matchRules(desc)
	if matched {
		result.Confidence = 0.95
		c.mu.Lock()
		c.cache[strings.ToLower(desc)] = result
		c.mu.Unlock()
		return result
	}

	// Fallback: use AI for unknown items
	result = c.aiClassify(desc)
	c.mu.Lock()
	c.cache[strings.ToLower(desc)] = result
	c.mu.Unlock()
	return result
}

// ClassifyBatch classifies multiple descriptions in parallel.
func (c *CargoClassifier) ClassifyBatch(descriptions []string) []CargoResult {
	results := make([]CargoResult, len(descriptions))
	var wg sync.WaitGroup
	for i, desc := range descriptions {
		wg.Add(1)
		go func(idx int, d string) {
			defer wg.Done()
			results[idx] = c.Classify(d)
		}(i, desc)
	}
	wg.Wait()
	return results
}

// matchRules checks all rule entries against the description.
func (c *CargoClassifier) matchRules(desc string) (CargoResult, bool) {
	lower := strings.ToLower(desc)
	for _, rule := range defaultRules {
		for _, kw := range rule.Keywords {
			if strings.Contains(lower, strings.ToLower(kw)) {
				return CargoResult{
					Category:    rule.Category,
					HSCode:      rule.HSCode,
					SubCategory: rule.SubCategory,
					RiskLevel:   rule.RiskLevel,
					NeedInspect: rule.NeedInspect,
				}, true
			}
		}
	}
	return CargoResult{}, false
}

// aiClassify uses the AI gateway to classify an unknown product.
func (c *CargoClassifier) aiClassify(desc string) CargoResult {
	if c.gateway == nil {
		return CargoResult{Category: "general", HSCode: "0000.00.0000", RiskLevel: "low", Confidence: 0.3}
	}

	resp, err := c.gateway.Chat(context.Background(), &gateway.ChatRequest{
		Messages: []gateway.Message{
			{Role: gateway.RoleSystem, Content: classifySystemPrompt},
			{Role: gateway.RoleUser, Content: desc},
		},
		Temperature: 0.1,
		MaxTokens:   200,
	})
	if err != nil || resp == nil {
		return CargoResult{Category: "general", HSCode: "0000.00.0000", RiskLevel: "low", Confidence: 0.3}
	}

	return parseAIResult(resp.Content)
}

const classifySystemPrompt = `You are a cargo classification expert for international logistics.
Classify the product into: general/electronic/dangerous/liquid/fragile/sensitive/powder.
Respond ONLY in this exact JSON format:
{"category":"<type>","hs_code":"<hs code>","sub_category":"<detail>","risk_level":"low|medium|high","need_inspect":true|false}`

// parseAIResult parses the AI response into a CargoResult.
func parseAIResult(content string) CargoResult {
	result := CargoResult{
		Category:   "general",
		HSCode:     "0000.00.0000",
		RiskLevel:  "low",
		Confidence: 0.5,
	}

	// Basic parsing of the expected JSON response
	lower := strings.ToLower(content)
	if strings.Contains(lower, `"category"`) {
		for _, cat := range []string{"dangerous", "fragile", "liquid", "sensitive", "powder", "electronic", "general"} {
			if strings.Contains(lower, fmt.Sprintf(`"%s"`, cat)) || strings.Contains(lower, fmt.Sprintf(`"%s`, cat)) {
				result.Category = cat
				break
			}
		}
	}
	if strings.Contains(lower, `"risk_level"`) {
		for _, level := range []string{"high", "medium", "low"} {
			if strings.Contains(lower, fmt.Sprintf(`"%s"`, level)) {
				result.RiskLevel = level
				break
			}
		}
	}
	if strings.Contains(lower, `"need_inspect":true`) {
		result.NeedInspect = true
	}

	return result
}

// CategoryCN returns Chinese display name for a category.
func CategoryCN(category string) string {
	switch category {
	case "general":
		return "普货"
	case "electronic":
		return "电子产品"
	case "dangerous":
		return "危险品"
	case "liquid":
		return "液体"
	case "fragile":
		return "易碎品"
	case "sensitive":
		return "特货"
	case "powder":
		return "粉末"
	default:
		return category
	}
}

// CategoryFromCN maps a Chinese category name back to the enum.
func CategoryFromCN(cn string) string {
	switch cn {
	case "普货":
		return "general"
	case "电子产品":
		return "electronic"
	case "危险品":
		return "dangerous"
	case "液体":
		return "liquid"
	case "易碎品":
		return "fragile"
	case "特货":
		return "sensitive"
	case "粉末":
		return "powder"
	default:
		return cn
	}
}
