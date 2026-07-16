// Package domain provides the core domain types and in-memory store for the WMS backend.
package domain

import (
	"sync"
	"sync/atomic"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// Generic in-memory store
// ═══════════════════════════════════════════════════════════════

// Store[T] is a generic concurrency-safe in-memory store.
type Store[T any] struct {
	mu     sync.RWMutex
	Items  []T
	NextID int64
}

// NewStore creates a new empty Store.
func NewStore[T any]() *Store[T] { return &Store[T]{Items: []T{}} }

// NextID returns the next auto-increment ID.
func (s *Store[T]) NextIdentifier() int64 { return atomic.AddInt64(&s.NextID, 1) }

// List returns a copy of all items.
func (s *Store[T]) List() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, len(s.Items))
	copy(result, s.Items)
	return result
}

// Add appends an item and returns it.
func (s *Store[T]) Add(item T) T {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Items = append(s.Items, item)
	return item
}

// Seed appends seed items and sets NextID.
func (s *Store[T]) Seed(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Items = append(s.Items, items...)
	s.NextID = int64(len(s.Items)) + 1
}

// Update replaces the item at the given index.
func (s *Store[T]) Update(idx int, item T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if idx < 0 || idx >= len(s.Items) { return false }
	s.Items[idx] = item
	return true
}

// Delete removes the item at the given index.
func (s *Store[T]) Delete(idx int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if idx < 0 || idx >= len(s.Items) { return false }
	s.Items = append(s.Items[:idx], s.Items[idx+1:]...)
	return true
}

// ═══════════════════════════════════════════════════════════════
// Domain models — TMS
// ═══════════════════════════════════════════════════════════════

type AreaGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Desc string `json:"description"`
}

type CargoType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type TransportMode struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type CustomsBroker struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	License string `json:"license"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
}

type CustomsPoint struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Port    string `json:"port"`
	Country string `json:"country"`
}

type ShippingProvider struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Contact string `json:"contact"`
}

type ContainerLoading struct {
	ID          int64     `json:"id"`
	ContainerNo string    `json:"container_no"`
	Vessel      string    `json:"vessel"`
	PortFrom    string    `json:"port_from"`
	PortTo      string    `json:"port_to"`
	ParcelCount int       `json:"parcel_count"`
	LoadedAt    time.Time `json:"loaded_at"`
}

type LogisticsTracking struct {
	ID          int64     `json:"id"`
	TrackingNo  string    `json:"tracking_no"`
	Route       string    `json:"route"`
	Status      string    `json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
	CourierName string    `json:"courier_name"`
	Detail      string    `json:"detail"`
}

type RouteTemplate struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	From      string `json:"from"`
	To        string `json:"to"`
	CarrierID int64  `json:"carrier_id"`
	EstDays   int    `json:"est_days"`
}

// ═══════════════════════════════════════════════════════════════
// Domain models — CRM
// ═══════════════════════════════════════════════════════════════

type ClientAccount struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	RealName string  `json:"real_name"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
}

type ClientRecharge struct {
	ID       int64     `json:"id"`
	ClientID int64     `json:"client_id"`
	Amount   float64   `json:"amount"`
	Method   string    `json:"method"`
	Remark   string    `json:"remark"`
	Time     time.Time `json:"time"`
}

type ClientPricing struct {
	ID       int64   `json:"id"`
	ClientID int64   `json:"client_id"`
	RouteID  int64   `json:"route_id"`
	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`
}

type ClientPermission struct {
	ID       int64  `json:"id"`
	ClientID int64  `json:"client_id"`
	Module   string `json:"module"`
	CanRead  bool   `json:"can_read"`
	CanWrite bool   `json:"can_write"`
}

type MonthlyStatement struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	Period     string    `json:"period"`
	Total      float64   `json:"total"`
	PaidAmount float64   `json:"paid_amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// ═══════════════════════════════════════════════════════════════
// Domain models — WMS / Services
// ═══════════════════════════════════════════════════════════════

type PDAWorkorderTemplate struct {
	ID              int64     `json:"id"`
	Warehouse       string    `json:"warehouse"`
	TemplateID      string    `json:"template_id"`
	Name            string    `json:"name"`
	WorkType        string    `json:"work_type"`
	WorkflowID      string    `json:"workflow_id"`
	DefaultPriority int       `json:"default_priority"`
	IsEnabled       bool      `json:"is_enabled"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WorkflowProcess struct {
	ID           int64     `json:"id"`
	Warehouse    string    `json:"warehouse"`
	ProcessID    string    `json:"process_id"`
	Name         string    `json:"name"`
	Steps        string    `json:"steps"`
	TriggerEvent string    `json:"trigger_event"`
	IsEnabled    bool      `json:"is_enabled"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type ServiceTemplate struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Type string  `json:"type"`
	Desc string  `json:"description"`
	Fee  float64 `json:"fee"`
}

type ServiceType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type ServiceWorkorder struct {
	ID         int64     `json:"id"`
	OrderNo    string    `json:"order_no"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	AssignedTo string    `json:"assigned_to"`
	CreatedAt  time.Time `json:"created_at"`
}

type PricingService struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Fee  float64 `json:"fee"`
	Unit string  `json:"unit"`
}

type Exception struct {
	ID          int64     `json:"id"`
	ParcelID    int64     `json:"parcel_id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type AIException struct {
	ID         int64     `json:"id"`
	ParcelID   int64     `json:"parcel_id"`
	Reason     string    `json:"reason"`
	Confidence float64   `json:"confidence"`
	Reviewed   bool      `json:"reviewed"`
	CreatedAt  time.Time `json:"created_at"`
}

type ExceptionReport struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Count     int       `json:"count"`
	Period    string    `json:"period"`
	CreatedAt time.Time `json:"created_at"`
}

type PDASession struct {
	ID             int64     `json:"id"`
	Warehouse      string    `json:"warehouse"`
	WorkerName     string    `json:"worker_name"`
	Device         string    `json:"device"`
	LoginAt        time.Time `json:"login_at"`
	LastHeartbeat  time.Time `json:"last_heartbeat"`
	OnlineDuration string    `json:"online_duration"`
	CurrentPage    string    `json:"current_page"`
	CurrentArea    string    `json:"current_area"`
	CurrentLocation string   `json:"current_location"`
	IsOnline       bool      `json:"is_online"`
	LogoutAt       *time.Time `json:"logout_at"`
}

type Printer struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	IP   string `json:"ip"`
}

type StorageConfig struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Bucket   string `json:"bucket"`
	Region   string `json:"region"`
}

type InboundBoardEntry struct {
	ID         int64     `json:"id"`
	ParcelNo   string    `json:"parcel_no"`
	Warehouse  string    `json:"warehouse"`
	Status     string    `json:"status"`
	ExpectedAt time.Time `json:"expected_at"`
}

type WarehouseBoardEntry struct {
	ID             int64 `json:"id"`
	PendingReceive int   `json:"pending_receive"`
	InStock        int   `json:"in_stock"`
	Picking        int   `json:"picking"`
	Outbound       int   `json:"outbound"`
}

type WarehouseConsoleEntry struct {
	ID          int64  `json:"id"`
	WarehouseID int64  `json:"warehouse_id"`
	Name        string `json:"name"`
	Machine     string `json:"machine"`
	Status      string `json:"status"`
}

// ═══════════════════════════════════════════════════════════════
// Domain models — System
// ═══════════════════════════════════════════════════════════════

type Notification struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Channel   string    `json:"channel"`
	Recipient string    `json:"recipient"`
	Sent      bool      `json:"sent"`
	CreatedAt time.Time `json:"created_at"`
}

type SystemParam struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Group string `json:"group"`
	Label string `json:"label"`
}

type BrandSetting struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Group string `json:"group"`
}

type APIConfig struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key"`
	Status   string `json:"status"`
}

type AIChatMessage struct {
	ID      int64     `json:"id"`
	Role    string    `json:"role"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

type SchedulerJob struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Cron    string `json:"cron"`
	Enabled bool   `json:"enabled"`
	LastRun string `json:"last_run"`
}

type AuditLog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

type Report struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ── BFT56-aligned domain models ──

// ClientMember 客户会员 (BFT56 client member)
type ClientMember struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	IDNumber  string    `json:"id_number"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"created_at"`
}

// BalanceLog 余额日志 (BFT56 balance log)
type BalanceLog struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Balance   float64   `json:"balance"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// RechargeRecord 充值记录 (BFT56 recharge log)
type RechargeRecord struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Container 集装柜 (BFT56 container)
type Container struct {
	ID          int64     `json:"id"`
	Warehouse   string    `json:"warehouse"`
	ContainerNo string    `json:"container_no"`
	RouteName   string    `json:"route_name"`
	Status      string    `json:"status"`
	MaxWeight   float64   `json:"max_weight"`
	CreatedAt   time.Time `json:"created_at"`
}

// ClientPanelPerm 客户端权限 (BFT56 client panel permission)
type ClientPanelPerm struct {
	ID          int64     `json:"id"`
	ClientID    int64     `json:"client_id"`
	ClientName  string    `json:"client_name"`
	Module      string    `json:"module"`
	MenuName    string    `json:"menu_name"`
	CanView     bool      `json:"can_view"`
	CanCreate   bool      `json:"can_create"`
	CanEdit     bool      `json:"can_edit"`
	CanDelete   bool      `json:"can_delete"`
	CanExport   bool      `json:"can_export"`
	Level       string    `json:"level"`
	Status      string    `json:"status"`
	GrantedAt   time.Time `json:"granted_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Remarks     string    `json:"remarks"`
}

type NotificationChannel struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
}

// ── 设备 (打印机/扫码枪/接地板) ──
type Device struct {
	ID          int64  `json:"id"`
	DeviceName  string `json:"device_name"`
	DeviceType  string `json:"device_type"`
	DeviceCode  string `json:"device_code"`
	IPAddress   string `json:"ip_address"`
	Status      string `json:"status"`
	WarehouseID int64  `json:"warehouse_id"`
}

// ── 仓位/货架 ──
type Shelf struct {
	ID          int64  `json:"id"`
	WarehouseID int64  `json:"warehouse_id"`
	Code        string `json:"code"`
	Zone        string `json:"zone"`
	Row         string `json:"row"`
	Level       int    `json:"level"`
	Status      string `json:"status"`
}

// ═══════════════════════════════════════════════════════════════
// Global data stores — singleton pattern
// ═══════════════════════════════════════════════════════════════

var (
	AreaGroupStore           = NewStore[AreaGroup]()
	CargoTypeStore           = NewStore[CargoType]()
	TransportModeStore       = NewStore[TransportMode]()
	CustomsBrokerStore       = NewStore[CustomsBroker]()
	CustomsPointStore        = NewStore[CustomsPoint]()
	ShippingProviderStore    = NewStore[ShippingProvider]()
	ContainerLoadingStore    = NewStore[ContainerLoading]()
	LogisticsTrackingStore   = NewStore[LogisticsTracking]()
	RouteTemplateStore       = NewStore[RouteTemplate]()
	ClientAccountStore       = NewStore[ClientAccount]()
	ClientRechargeStore      = NewStore[ClientRecharge]()
	ClientPricingStore       = NewStore[ClientPricing]()
	ClientPermissionStore    = NewStore[ClientPermission]()
	MonthlyStatementStore    = NewStore[MonthlyStatement]()
	NotificationStore        = NewStore[Notification]()
	PDAWorkorderTplStore     = NewStore[PDAWorkorderTemplate]()
	WorkflowProcessStore     = NewStore[WorkflowProcess]()
	RoleStore                = NewStore[Role]()
	ClientMemberStore        = NewStore[ClientMember]()
	BalanceLogStore          = NewStore[BalanceLog]()
	RechargeRecordStore      = NewStore[RechargeRecord]()
	ContainerStore           = NewStore[Container]()
	ClientPanelPermStore     = NewStore[ClientPanelPerm]()
	DeviceStore              = NewStore[Device]()
	ShelfStore               = NewStore[Shelf]()
	ServiceTemplateStore     = NewStore[ServiceTemplate]()
	ServiceTypeStore         = NewStore[ServiceType]()
	ServiceWorkorderStore    = NewStore[ServiceWorkorder]()
	PricingServiceStore      = NewStore[PricingService]()
	ExceptionStore           = NewStore[Exception]()
	PDASessionStore          = NewStore[PDASession]()
	PrinterStore             = NewStore[Printer]()
	StorageConfigStore       = NewStore[StorageConfig]()
	SystemParamStore         = NewStore[SystemParam]()
	BrandSettingStore        = NewStore[BrandSetting]()
	APIConfigStore           = NewStore[APIConfig]()
	AIChatStore              = NewStore[AIChatMessage]()
	SchedulerJobStore        = NewStore[SchedulerJob]()
	AuditLogStore            = NewStore[AuditLog]()
	ReportStore              = NewStore[Report]()
	NotificationChannelStore = NewStore[NotificationChannel]()
	InboundBoardStore        = NewStore[InboundBoardEntry]()
	WarehouseBoardStore      = NewStore[WarehouseBoardEntry]()
	WarehouseConsoleStore    = NewStore[WarehouseConsoleEntry]()
	AIExceptionStore         = NewStore[AIException]()
	ExceptionReportStore     = NewStore[ExceptionReport]()
)
