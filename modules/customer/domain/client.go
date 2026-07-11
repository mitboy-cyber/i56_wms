// Package domain defines the Customer bounded context domain model.
package domain

import "time"

// ClientType represents the client category.
type ClientType string

const (
	ClientTypePlatform ClientType = "platform"  // 平台客户
	ClientTypeShopee   ClientType = "shopee"    // 虾皮商家
	ClientTypeMajor    ClientType = "major"     // 大客户
	ClientTypePeer     ClientType = "peer"      // 同行
	ClientTypeNormal   ClientType = "normal"    // 普通客户
)

// Client is the aggregate root for a tenant's business customer.
type Client struct {
	ID           int64      `json:"id"`
	TenantID     int64      `json:"tenant_id"`
	Name         string     `json:"name"`
	Code         string     `json:"code"`
	ClientType   ClientType `json:"client_type"`
	ContactName  string     `json:"contact_name"`
	ContactPhone string     `json:"contact_phone"`
	ContactEmail string     `json:"contact_email"`
	Balance      float64    `json:"balance"`
	IsActive     bool       `json:"is_active"`
	Remark       string     `json:"remark"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ValidClientTypes returns all valid client type values.
func ValidClientTypes() []string {
	return []string{"platform", "shopee", "major", "peer", "normal"}
}
