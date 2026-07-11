package domain

import "time"

// Permission represents an individual permission in the system.
type Permission struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Module      string    `json:"module"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Role represents a named role with a set of permissions.
type Role struct {
	ID            int64     `json:"id"`
	TenantID      int64     `json:"tenant_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	PermissionIDs []int64   `json:"permission_ids"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ClientPermission represents permissions granted to a specific client.
type ClientPermission struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenant_id"`
	ClientID        int64     `json:"client_id"`
	ClientName      string    `json:"client_name"`
	PermissionSlugs []string  `json:"permission_slugs"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// User represents a system user with a role assignment.
type User struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	RealName  string    `json:"real_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	RoleID    int64     `json:"role_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DefaultPermissions returns the set of seed permissions in 繁體中文.
func DefaultPermissions() []Permission {
	modules := []struct {
		Module string
		Name   string
		Slug   string
		Desc   string
	}{
		// 儀表板
		{Module: "儀表板", Name: "查看儀表板", Slug: "dashboard:view", Desc: "查看系統儀表板"},
		// 包裹管理
		{Module: "包裹管理", Name: "查看包裹", Slug: "parcel:list", Desc: "查看包裹列表"},
		{Module: "包裹管理", Name: "創建包裹", Slug: "parcel:create", Desc: "預報入庫包裹"},
		{Module: "包裹管理", Name: "編輯包裹", Slug: "parcel:update", Desc: "修改包裹資訊"},
		{Module: "包裹管理", Name: "刪除包裹", Slug: "parcel:delete", Desc: "刪除包裹記錄"},
		{Module: "包裹管理", Name: "入庫操作", Slug: "parcel:inbound", Desc: "執行入庫稱重上架"},
		{Module: "包裹管理", Name: "匯出包裹", Slug: "parcel:export", Desc: "匯出包裹報表"},
		// 訂單管理
		{Module: "訂單管理", Name: "查看訂單", Slug: "order:list", Desc: "查看集運訂單列表"},
		{Module: "訂單管理", Name: "創建訂單", Slug: "order:create", Desc: "創建集運訂單"},
		{Module: "訂單管理", Name: "編輯訂單", Slug: "order:update", Desc: "修改訂單資訊"},
		{Module: "訂單管理", Name: "刪除訂單", Slug: "order:delete", Desc: "刪除訂單"},
		{Module: "訂單管理", Name: "審核訂單", Slug: "order:approve", Desc: "審核通過訂單"},
		{Module: "訂單管理", Name: "取消訂單", Slug: "order:cancel", Desc: "取消訂單"},
		{Module: "訂單管理", Name: "匯出訂單", Slug: "order:export", Desc: "匯出訂單報表"},
		// 客戶管理
		{Module: "客戶管理", Name: "查看客戶", Slug: "client:list", Desc: "查看客戶列表"},
		{Module: "客戶管理", Name: "創建客戶", Slug: "client:create", Desc: "新增客戶"},
		{Module: "客戶管理", Name: "編輯客戶", Slug: "client:update", Desc: "修改客戶資訊"},
		{Module: "客戶管理", Name: "刪除客戶", Slug: "client:delete", Desc: "刪除客戶"},
		{Module: "客戶管理", Name: "查看客戶財務", Slug: "client:finance", Desc: "查看客戶餘額與帳單"},
		{Module: "客戶管理", Name: "管理客戶權限", Slug: "client:permissions", Desc: "設定客戶端功能權限"},
		// 倉庫管理
		{Module: "倉庫管理", Name: "查看倉庫", Slug: "warehouse:list", Desc: "查看倉庫列表"},
		{Module: "倉庫管理", Name: "創建倉庫", Slug: "warehouse:create", Desc: "新增倉庫"},
		{Module: "倉庫管理", Name: "編輯倉庫", Slug: "warehouse:update", Desc: "修改倉庫資訊"},
		{Module: "倉庫管理", Name: "刪除倉庫", Slug: "warehouse:delete", Desc: "刪除倉庫"},
		{Module: "倉庫管理", Name: "管理庫位", Slug: "warehouse:locations", Desc: "管理庫位與區域"},
		{Module: "倉庫管理", Name: "查看入庫看板", Slug: "warehouse:inbound-board", Desc: "查看入庫看板"},
		// 線路管理
		{Module: "線路管理", Name: "查看線路", Slug: "route:list", Desc: "查看運輸線路"},
		{Module: "線路管理", Name: "創建線路", Slug: "route:create", Desc: "新增運輸線路"},
		{Module: "線路管理", Name: "編輯線路", Slug: "route:update", Desc: "修改線路資訊"},
		{Module: "線路管理", Name: "刪除線路", Slug: "route:delete", Desc: "刪除線路"},
		{Module: "線路管理", Name: "管理報價", Slug: "route:pricing", Desc: "設定線路價格模板"},
		// 財務管理
		{Module: "財務管理", Name: "查看帳務", Slug: "finance:list", Desc: "查看財務帳務記錄"},
		{Module: "財務管理", Name: "客戶充值", Slug: "finance:recharge", Desc: "為客戶充值"},
		{Module: "財務管理", Name: "查看報表", Slug: "finance:reports", Desc: "查看財務報表"},
		{Module: "財務管理", Name: "匯出報表", Slug: "finance:export", Desc: "匯出財務報表"},
		{Module: "財務管理", Name: "對帳結算", Slug: "finance:settlement", Desc: "月結對帳操作"},
		// 系統管理
		{Module: "系統管理", Name: "查看用戶", Slug: "user:list", Desc: "查看系統用戶列表"},
		{Module: "系統管理", Name: "創建用戶", Slug: "user:create", Desc: "新增系統用戶"},
		{Module: "系統管理", Name: "編輯用戶", Slug: "user:update", Desc: "修改用戶資訊"},
		{Module: "系統管理", Name: "刪除用戶", Slug: "user:delete", Desc: "刪除用戶"},
		{Module: "系統管理", Name: "查看角色", Slug: "role:list", Desc: "查看角色列表"},
		{Module: "系統管理", Name: "創建角色", Slug: "role:create", Desc: "新增角色"},
		{Module: "系統管理", Name: "編輯角色", Slug: "role:update", Desc: "修改角色與權限"},
		{Module: "系統管理", Name: "刪除角色", Slug: "role:delete", Desc: "刪除角色"},
		{Module: "系統管理", Name: "系統設定", Slug: "system:settings", Desc: "管理系統參數配置"},
		{Module: "系統管理", Name: "查看操作日誌", Slug: "system:audit", Desc: "查看操作審計日誌"},
		// 附加服務
		{Module: "附加服務", Name: "查看服務", Slug: "service:list", Desc: "查看附加服務"},
		{Module: "附加服務", Name: "管理服務", Slug: "service:manage", Desc: "管理附加服務項目"},
		// 報表
		{Module: "報表", Name: "查看報表", Slug: "report:list", Desc: "查看各類營運報表"},
		{Module: "報表", Name: "匯出報表", Slug: "report:export", Desc: "匯出報表數據"},
	}

	perms := make([]Permission, 0, len(modules))
	for i, m := range modules {
		perms = append(perms, Permission{
			ID:          int64(i + 1),
			Name:        m.Name,
			Slug:        m.Slug,
			Module:      m.Module,
			Description: m.Desc,
			IsActive:    true,
		})
	}
	return perms
}

// DefaultRoles returns the 4 seed roles in 繁體中文.
func DefaultRoles() []Role {
	allPermIDs := func() []int64 {
		perms := DefaultPermissions()
		ids := make([]int64, len(perms))
		for i, p := range perms {
			ids[i] = p.ID
		}
		return ids
	}

	// 倉庫管理員權限：包裹、倉庫、入庫相關
	whPermIDs := func() []int64 {
		var ids []int64
		for _, p := range DefaultPermissions() {
			switch p.Slug {
			case "dashboard:view",
				"parcel:list", "parcel:create", "parcel:update", "parcel:inbound", "parcel:export",
				"order:list", "order:update",
				"client:list",
				"warehouse:list", "warehouse:locations", "warehouse:inbound-board",
				"route:list",
				"service:list",
				"report:list":
				ids = append(ids, p.ID)
			}
		}
		return ids
	}

	// 財務專員權限：財務、報表、客戶財務
	finPermIDs := func() []int64 {
		var ids []int64
		for _, p := range DefaultPermissions() {
			switch p.Slug {
			case "dashboard:view",
				"parcel:list",
				"order:list",
				"client:list", "client:finance",
				"finance:list", "finance:recharge", "finance:reports", "finance:export", "finance:settlement",
				"report:list", "report:export":
				ids = append(ids, p.ID)
			}
		}
		return ids
	}

	// 客服專員權限：客戶、訂單查看
	csPermIDs := func() []int64 {
		var ids []int64
		for _, p := range DefaultPermissions() {
			switch p.Slug {
			case "dashboard:view",
				"parcel:list",
				"order:list", "order:create", "order:update",
				"client:list",
				"warehouse:list",
				"route:list",
				"service:list":
				ids = append(ids, p.ID)
			}
		}
		return ids
	}

	return []Role{
		{
			ID:            1,
			TenantID:      1,
			Name:          "超级管理员",
			Slug:          "super_admin",
			Description:   "擁有系統所有權限，可管理全部模組與設定",
			PermissionIDs: allPermIDs(),
			IsActive:      true,
		},
		{
			ID:            2,
			TenantID:      1,
			Name:          "倉庫管理員",
			Slug:          "warehouse_manager",
			Description:   "負責倉庫日常營運，管理包裹入庫、庫位、訂單處理",
			PermissionIDs: whPermIDs(),
			IsActive:      true,
		},
		{
			ID:            3,
			TenantID:      1,
			Name:          "財務專員",
			Slug:          "finance_specialist",
			Description:   "負責財務管理、客戶充值、對帳結算與報表",
			PermissionIDs: finPermIDs(),
			IsActive:      true,
		},
		{
			ID:            4,
			TenantID:      1,
			Name:          "客服專員",
			Slug:          "customer_service",
			Description:   "處理客戶諮詢，查看包裹與訂單狀態",
			PermissionIDs: csPermIDs(),
			IsActive:      true,
		},
	}
}

// DefaultUsers returns 3 seed users.
func DefaultUsers() []User {
	return []User{
		{
			ID:       1,
			TenantID: 1,
			Username: "admin",
			Password: "admin",
			RealName: "系統管理員",
			Email:    "admin@i56.com",
			Phone:    "13800000010",
			RoleID:   1, // 超级管理员
			IsActive: true,
		},
		{
			ID:       2,
			TenantID: 1,
			Username: "dabao",
			Password: "dabao123",
			RealName: "大寶",
			Email:    "dabao@i56.com",
			Phone:    "13800000015",
			RoleID:   2, // 倉庫管理員
			IsActive: true,
		},
		{
			ID:       3,
			TenantID: 1,
			Username: "xiaolin",
			Password: "xiaolin123",
			RealName: "小林",
			Email:    "xiaolin@i56.com",
			Phone:    "13800000020",
			RoleID:   3, // 財務專員
			IsActive: true,
		},
	}
}

// DefaultClientPermissions returns 2 seed client permissions.
func DefaultClientPermissions() []ClientPermission {
	return []ClientPermission{
		{
			ID:              1,
			TenantID:        1,
			ClientID:        1,
			ClientName:      "EZ集運通",
			PermissionSlugs: []string{"parcel:create", "parcel:list", "order:create", "order:list", "client:list"},
			IsActive:        true,
		},
		{
			ID:              2,
			TenantID:        1,
			ClientID:        2,
			ClientName:      "嗨購商城",
			PermissionSlugs: []string{"parcel:create", "parcel:list", "order:create", "order:list"},
			IsActive:        true,
		},
	}
}
